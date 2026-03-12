"""
LMFDB API Client - 使用 Playwright 绕过 reCAPTCHA
"""
from __future__ import annotations

import json
import time
from pathlib import Path
from typing import Any

import requests
from playwright.sync_api import sync_playwright

# LMFDB API Base URL
BASE_URL = "https://www.lmfdb.org"


class LMFDBClient:
    """LMFDB API 客户端 - 使用 Playwright 绕过 reCAPTCHA"""
    
    def __init__(self, headless: bool = True, timeout: int = 30000):
        self.headless = headless
        self.timeout = timeout
        self._playwright = None
        self._browser = None
        self._context = None
        self._page = None
        self._session_cookies = None
    
    def _init_browser(self):
        """初始化浏览器"""
        if self._playwright is None:
            self._playwright = sync_playwright().start()
            self._browser = self._playwright.chromium.launch(headless=self.headless)
            self._context = self._browser.new_context()
            self._page = self._context.new_page()
    
    def _wait_for_recaptcha(self) -> bool:
        """等待 reCAPTCHA 通过"""
        # 检查是否有 reCAPTCHA
        try:
            # 等待页面加载完成
            self._page.wait_for_load_state("networkidle", timeout=10000)
            
            # 检查是否有 recaptcha iframe
            recaptcha_frames = self._page.frames
            
            for frame in recaptcha_frames:
                try:
                    if "google.com/recaptcha" in frame.url:
                        # 等待 reCAPTCHA 执行
                        self._page.wait_for_timeout(3000)
                        break
                except:
                    continue
            
            # 检查是否还有 reCAPTCHA 保护
            content = self._page.content()
            if "recaptcha" in content.lower() and "checking your browser" in content.lower():
                # 等待更长时间
                self._page.wait_for_timeout(5000)
            
            return True
        except Exception as e:
            print(f"Warning: reCAPTCHA check: {e}")
            return False
    
    def _build_url(self, collection: str, params: dict[str, Any]) -> str:
        """构建 API URL"""
        url = f"{BASE_URL}/api/{collection}/"
        
        # 构建查询参数
        query_parts = []
        for key, value in params.items():
            if value is not None:
                query_parts.append(f"{key}={value}")
        
        if query_parts:
            url += "?" + "&".join(query_parts)
        
        return url
    
    def query(
        self,
        collection: str,
        params: dict[str, Any] | None = None,
        fields: list[str] | None = None,
        format: str = "json",
        limit: int = 100,
        offset: int = 0,
    ) -> dict[str, Any]:
        """
        查询 LMFDB API
        
        Args:
            collection: API 集合名称 (如 nf_fields, ec_curvedata 等)
            params: 查询参数
            fields: 返回字段列表
            format: 输出格式 (json/yaml/html)
            limit: 返回结果数量限制
            offset: 偏移量
            
        Returns:
            API 响应数据
        """
        if params is None:
            params = {}
        
        # 构建参数
        query_params = {}
        
        # 添加用户参数
        for key, value in params.items():
            if value is not None:
                query_params[key] = value
        
        # 添加格式参数
        if format:
            query_params["_format"] = format
        
        # 添加字段参数
        if fields:
            query_params["_fields"] = ",".join(fields)
        
        # 添加数量限制
        if limit:
            query_params["_limit"] = limit
        
        if offset:
            query_params["_offset"] = offset
        
        url = self._build_url(collection, query_params)
        
        # 尝试直接请求 (可能已经被缓存绕过)
        try:
            response = requests.get(url, timeout=30)
            
            # 检查是否被 reCAPTCHA 拦截
            if "recaptcha" in response.text.lower() or "checking your browser" in response.text.lower():
                # 需要使用 Playwright
                return self._query_with_playwright(collection, query_params)
            
            if format == "json":
                return response.json()
            else:
                return {"raw": response.text}
                
        except requests.RequestException as e:
            # 如果直接请求失败，使用 Playwright
            print(f"Direct request failed: {e}, trying with Playwright...")
            return self._query_with_playwright(collection, query_params)
    
    def _query_with_playwright(
        self,
        collection: str,
        params: dict[str, Any]
    ) -> dict[str, Any]:
        """使用 Playwright 查询 (绕过 reCAPTCHA)"""
        self._init_browser()
        
        url = self._build_url(collection, params)
        
        try:
            # 访问页面
            self._page.goto(url, wait_until="networkidle", timeout=self.timeout)
            
            # 等待 reCAPTCHA 通过
            self._wait_for_recaptcha()
            
            # 再次等待网络空闲
            self._page.wait_for_load_state("networkidle", timeout=10000)
            
            # 获取内容
            content = self._page.content()
            
            # 尝试解析 JSON
            try:
                # 找到 JSON 数据
                json_start = content.find('{')
                json_end = content.rfind('}') + 1
                
                if json_start >= 0 and json_end > json_start:
                    json_str = content[json_start:json_end]
                    return json.loads(json_str)
            except json.JSONDecodeError:
                pass
            
            # 返回原始内容
            return {"raw": content, "url": url}
            
        finally:
            pass  # 保持浏览器打开以便复用
    
    def get_collections(self) -> list[str]:
        """获取所有可用的集合"""
        data = self._query_with_playwright("", {"_format": "json"})
        
        collections = []
        if isinstance(data, dict):
            # 解析集合列表
            for key in data.keys():
                if not key.startswith('_'):
                    collections.append(key)
        
        return collections
    
    def close(self):
        """关闭浏览器"""
        if self._browser:
            self._browser.close()
        if self._playwright:
            self._playwright.stop()
    
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()


# 全局客户端实例
_client: LMFDBClient | None = None


def get_client(headless: bool = True) -> LMFDBClient:
    """获取全局客户端实例"""
    global _client
    if _client is None:
        _client = LMFDBClient(headless=headless)
    return _client


def close_client():
    """关闭全局客户端"""
    global _client
    if _client:
        _client.close()
        _client = None
