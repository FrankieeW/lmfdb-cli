# LMFDB CLI

[English](./README.md) | 中文

一个用于查询 **LMFDB**（L-函数与模形式数据库）的命令行工具。

## 功能特性

- 纯 Go 实现，跨平台
- 使用 chromedp 绕过 reCAPTCHA 验证（`--browser` 参数）
- 支持所有 LMFDB API 集合
- 支持 Table/JSON/CSV 输出格式
- 支持过滤和排序

## 安装

### Homebrew（macOS/Linux）
```bash
brew tap frankieew/tap
brew install lmfdb-cli

# 安装浏览器（用于绕过 reCAPTCHA）
lmfdb install-browser
```

### 从源码构建
```bash
git clone https://github.com/FrankieeW/lmfdb-cli.git
cd lmfdb-cli
go build -o lmfdb ./cmd/lmfdb
./lmfdb install-browser  # 可选：用于绕过 reCAPTCHA
```

### 下载预编译版本
从 [GitHub Releases](https://github.com/FrankieeW/lmfdb-cli/releases) 下载。

## 快速开始

```bash
# 查询二次域（默认 degree=2）
lmfdb nf -d 2 -n 10

# 查询椭圆曲线
lmfdb ec -n 10

# 使用浏览器绕过 reCAPTCHA
lmfdb nf -d 2 -n 10 --browser

# 列出可用集合
lmfdb list
```

## 命令详解

### nf - 数域

```bash
lmfdb nf -d 2              # 二次域
lmfdb nf -d 3              # 三次域
lmfdb nf -d 2 --disc -5    # 按判别式筛选
lmfdb nf -d 2 -n 20        # 限制结果数量
lmfdb nf --id 2.0.3.1      # 按标签查询
lmfdb nf -d 2 --fmt json   # JSON 输出
lmfdb nf --browser -n 50   # 使用浏览器绕过 reCAPTCHA
```

参数选项：
- `-d` — 数域次数（默认: 2）
- `--disc` — 按判别式筛选
- `--class` — 按类数筛选
- `--sig` — 按 signature 筛选（如 `0,1`）
- `-n` — 结果数量（默认: 10）
- `--offset` — 分页偏移量
- `--sort` — 排序字段（`-` 前缀为降序）
- `-f` — 返回字段（逗号分隔）
- `-o` — 输出文件
- `--fmt` — 输出格式：`table`、`json`、`csv`
- `--id` — 按标签获取
- `-q` — 静默模式
- `--browser` — 使用浏览器（绕过 reCAPTCHA）

### ec - 椭圆曲线

```bash
lmfdb ec -n 10              # 列出曲线
lmfdb ec -r 2               # 按秩筛选
lmfdb ec -t 5               # 按 torsion 筛选
lmfdb ec --conductor 11     # 按 conductor 筛选
```

参数选项：
- `-r` — 按 Mordell-Weil 秩筛选
- `-t` — 按 torsion 阶筛选
- `--conductor` — 按 conductor 筛选
- `-n`、`--offset`、`--sort`、`-f`、`-o`、`--fmt`、`-q`、`--browser` — 同 `nf`

### list (ls) - 可用集合

```bash
lmfdb list
```

| 集合名称 | 说明 |
|----------|------|
| `nf_fields` | 数域 |
| `ec_curvedata` | 椭圆曲线 |
| `ec_classdata` | 椭圆曲线等价类 |
| `g2c_curves` | 亏格 2 曲线 |
| `char_dirichlet` | Dirichlet 特征标 |
| `maass_newforms` | Maass 形式 |
| `mf_newforms` | 模形式 |
| `lf_fields` | 局部域 |
| `artin` | Artin 表示 |
| `belyi` | Belyi 映射 |

## 许可证

MIT

## 致谢

- [LMFDB](https://www.lmfdb.org/) - L-函数与模形式数据库
