# LMFDB CLI

[English](./README.md) | 中文

用于查询 [LMFDB](https://www.lmfdb.org/)（L-函数与模形式数据库）的命令行工具。

## 特性

- 纯 Go 实现，单一二进制文件，无依赖
- 通过无头 Chrome 绕过 reCAPTCHA（`--browser`）
- JSON 输出语法高亮
- 支持 Table / JSON / CSV 输出格式
- 灵活的过滤、排序和分页
- 跨平台（Linux、macOS、Windows）

## 安装

### Homebrew（macOS / Linux）
```bash
brew tap frankieew/tap
brew install lmfdb-cli
```

### 下载预编译版本
从 [GitHub Releases](https://github.com/FrankieeW/lmfdb-cli/releases) 下载。

### 从源码构建
```bash
git clone https://github.com/FrankieeW/lmfdb-cli.git
cd lmfdb-cli
go build -o lmfdb ./cmd/lmfdb
```

## 快速开始

```bash
lmfdb nf -d 2 -n 10               # 二次数域
lmfdb ec -r 2 --fmt json           # 秩为 2 的椭圆曲线，JSON 输出
lmfdb nf --id 2.0.3.1              # 按标签查询
lmfdb nf --browser -n 50           # 绕过 reCAPTCHA
lmfdb list                         # 可用集合
```

## 命令

| 命令 | 说明 |
|------|------|
| `nf` | 查询数域 |
| `ec` | 查询椭圆曲线 |
| `list` / `ls` | 列出可用的 API 集合 |
| `version` / `v` | 显示版本 |
| `install-browser` | 安装 Chrome 浏览器（用于绕过 reCAPTCHA） |

使用 `lmfdb <command> -h` 查看各命令的详细帮助。

## 输出格式

```bash
lmfdb nf -d 2 -n 5                 # 表格（默认）
lmfdb nf -d 2 -n 5 --fmt json      # 带语法高亮的 JSON
lmfdb nf -d 2 -n 5 --fmt csv       # CSV
lmfdb nf -d 2 -n 5 -o out.json     # 保存到文件
```

## 文档

- [使用指南](./docs/GUIDE.zh-CN.md) - 详细使用说明与示例
- [API 参考](./docs/API.zh-CN.md) - LMFDB API 查询语法
- [文档索引](./docs/INDEX.md)

## 许可证

MIT

## 致谢

- [LMFDB](https://www.lmfdb.org/) - L-函数与模形式数据库
