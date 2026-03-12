# LMFDB API 参考文档

[English](./API.md) | 中文

本文档介绍 LMFDB API 的查询语法。

## 基础 URL

```
https://www.lmfdb.org/api/
```

## 查询格式

### 类型前缀

查询值必须带类型前缀：

| 前缀 | 类型 | 示例 |
|------|------|------|
| `i` | 整数 | `degree=i2` |
| `f` | 浮点数 | `rank=f2.5` |
| `s` | 字符串 | `label=s11a1` |
| `li` | 整数列表 | `torsion=li2;4` |
| `ls` | 字符串列表 | `coeffs=ls1;-1;1` |
| `py` | Python 字面量 | `matrix=py[[1,2],[3,4]]` |
| `cs` | 包含字符串 | `label=cs11a` |
| `ci` | 包含整数 | `torsion_structure=ci2` |

### 元参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `_format` | 输出格式 (json/yaml/html) | `_format=json` |
| `_fields` | 返回字段 | `_fields=label,degree,disc` |
| `_sort` | 排序字段（`-` 前缀表示降序） | `_sort=-degree,label` |
| `_limit` | 最大结果数 | `_limit=100` |
| `_offset` | 结果偏移量 | `_offset=50` |
| `_delim` | 列表分隔符（默认: `,`） | `_delim=;` |

## 查询示例

### 数域 (Number Fields)

```bash
# 二次域
https://www.lmfdb.org/api/nf_fields/?degree=i2&_format=json

# 按判别式筛选
https://www.lmfdb.org/api/nf_fields/?disc=i-5&_format=json

# 类数为 1
https://www.lmfdb.org/api/nf_fields/?class_number=i1&_format=json

# 多条件组合
https://www.lmfdb.org/api/nf_fields/?degree=i2&class_number=i1&_format=json

# 选择特定字段
https://www.lmfdb.org/api/nf_fields/?degree=i2&_fields=label,degree,disc,class_number&_format=json
```

### 椭圆曲线 (Elliptic Curves)

```bash
# 秩 = 2
https://www.lmfdb.org/api/ec_curvedata/?rank=i2&_format=json

# Torsion = 5
https://www.lmfdb.org/api/ec_curvedata/?torsion=i5&_format=json

# Conductor = 11
https://www.lmfdb.org/api/ec_curvedata/?conductor=i11&_format=json

# 按秩降序排序
https://www.lmfdb.org/api/ec_curvedata/?_sort=-rank&_format=json
```

### 使用 curl

```bash
# 查询数域
curl "https://www.lmfdb.org/api/nf_fields/?degree=i2&_format=json"

# 指定返回字段
curl "https://www.lmfdb.org/api/nf_fields/?degree=i2&_fields=label,disc&_limit=10&_format=json"
```

## 常用字段名

### nf_fields (数域)

| 字段 | 说明 |
|------|------|
| `label` | LMFDB 标签 |
| `degree` | 次数 |
| `disc` | 判别式 |
| `class_number` | 类数 |
| `class_group` | 类群结构 |
| `signature` | 签名 [r1, r2] |
| `coefficients` | 定义多项式系数 |
| `cm` | 是否有复乘法 |

### ec_curvedata (椭圆曲线)

| 字段 | 说明 |
|------|------|
| `label` | LMFDB 标签 (如 "11a1") |
| `conductor` | 导子 |
| `rank` | Mordell-Weil 秩 |
| `torsion` | 挠系数 |
| `torsion_structure` | 挠结构 (如 [2,2]) |
| `equation` | 极小方程 |
| `j_invariant` | j-不变量 |
| `cid` | 等价类 ID |

## 限制

- 每次查询最大结果数: **10,000**
- 默认限制: **100**
- 使用 `_offset` 进行分页

## 速率限制

请合理使用 API。LMFDB 为研究目的提供此服务。

## 错误信息

- `Unable to locate the page you requested` - 无效的集合或字段名
- `recaptcha` 响应 - 请求被拦截，请重试或使用浏览器
- HTTP 404 - 无效的 URL
- HTTP 429 - 请求过于频繁
