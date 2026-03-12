# LMFDB CLI

[English](./README.md) | 中文

一个用于查询 **LMFDB**（L-函数与模形式数据库）的命令行工具。

## 功能特性

- ✅ 使用 Playwright 绕过 reCAPTCHA 验证
- ✅ 支持所有 LMFDB API 集合
- ✅ 使用 Rich 美化表格输出
- ✅ 支持自定义返回字段
- ✅ 支持 JSON/YAML 输出格式
- ✅ 支持过滤和排序

## 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/lmfdb-cli.git
cd lmfdb-cli

# 安装依赖
pip install -e .

# 安装 Playwright 浏览器
playwright install chromium
```

## 快速开始

```bash
# 查询二次域（默认 degree=2）
lmfdb nf -d 2 -n 10

# 查询椭圆曲线
lmfdb ec -n 10

# 列出可用集合
lmfdb list-collections

# 通用查询
lmfdb query nf_fields -n 5
```

## 命令详解

### nf - 数域 (Number Fields)

查询 LMFDB 中的数域。

```bash
# 查询二次域
lmfdb nf -d 2

# 查询三次域
lmfdb nf -d 3

# 按判别式筛选
lmfdb nf -d 2 --discriminant -5

# 限制结果数量
lmfdb nf -d 2 -n 20

# 输出到 JSON 文件
lmfdb nf -d 2 -n 10 -o results.json
```

参数选项：
- `-d, --degree`: 数域次数（默认: 2）
- `--disc`: 按判别式筛选
- `-h, --class-number`: 按类数筛选
- `-n, --limit`: 结果数量限制（默认: 10）
- `-f, --fields`: 逗号分隔的返回字段
- `-o, --output`: 输出文件路径（JSON 格式）
- `--headless/--no-headless`: 浏览器模式

### ec - 椭圆曲线 (Elliptic Curves)

查询 LMFDB 中的椭圆曲线。

```bash
# 查询所有椭圆曲线
lmfdb ec -n 10

# 按秩筛选
lmfdb ec -r 2

# 按 torsion 筛选
lmfdb ec -t 5

# 按 conductor 筛选
lmfdb ec --conductor 11

# 自定义字段
lmfdb ec -n 10 -f label,conductor,rank,torsion
```

参数选项：
- `-r, --rank`: 按 Mordell-Weil 秩筛选
- `-t, --torsion`: 按 torsion 阶筛选
- `--conductor`: 按 conductor 筛选
- `-n, --limit`: 结果数量限制（默认: 10）
- `-f, --fields`: 逗号分隔的返回字段
- `-o, --output`: 输出文件路径（JSON 格式）

### query - 通用查询

查询任意 LMFDB API 集合。

```bash
# 列出所有集合
lmfdb list-collections

# 查询指定集合
lmfdb query nf_fields -n 5
lmfdb query ec_curvedata -n 10

# 带自定义参数（JSON 格式）
lmfdb query nf_fields -k '{"degree": "i3"}' -n 10
```

## 可用集合

| 集合名称 | 说明 |
|----------|------|
| `nf_fields` | 数域 |
| `ec_curvedata` | 椭圆曲线 |
| `ec_classdata` | 椭圆曲线等价类 |
| `g2c_curves` | 二 genus 曲线 |
| `char_dirichlet` | Dirichlet 特征标 |
| `maass_newforms` | Maass 形式 |
| `mf_newforms` | 模形式 |
| `lf_fields` | 局部域 |
| `artin` | Artin 表示 |

## 环境配置

```bash
# 设置超时时间（毫秒）
export LMFDB_TIMEOUT=60000

# 设置无头模式
export LMFDB_HEADLESS=true
```

## 开发

```bash
# 安装开发依赖
pip install -e ".[dev]"

# 运行测试
pytest

# 代码格式化
ruff check .
ruff format .
```

## 使用示例

### 查询二次域并导出

```bash
lmfdb nf -d 2 -n 100 -o quadratic_fields.json
```

### 查询特定秩的椭圆曲线

```bash
lmfdb ec -r 0 -n 50  # rank = 0 ( Mordell 猜想相关)
```

### 查询类数为 1 的数域

```bash
lmfdb nf -d 2 --class-number 1
```

## 常见问题

### Q: 为什么需要 Playwright？

A: LMFDB 使用 Google reCAPTCHA 保护其 API，直接请求会被拦截。Playwright 通过真实浏览器访问，可以自动绕过验证。

### Q: 第一次运行很慢？

A: 首次运行时，Playwright 需要启动浏览器并等待 reCAPTCHA 验证。后续请求会更快。

### Q: 如何调试？

A: 使用 `--no-headless` 选项可以看到浏览器操作过程：

```bash
lmfdb nf -d 2 --no-headless
```

## 相关链接

- [LMFDB 官网](https://www.lmfdb.org/)
- [LMFDB API 文档](https://www.lmfdb.org/api/)
- [Playwright 文档](https://playwright.dev/)
- [Rich 文档](https://rich.readthedocs.io/)

## 许可证

MIT

## 致谢

- [LMFDB](https://www.lmfdb.org/) - L-函数与模形式数据库
- [Playwright](https://playwright.dev/) - 浏览器自动化
- [Rich](https://rich.readthedocs.io/) - 终端美化输出
