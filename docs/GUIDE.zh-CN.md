# LMFDB CLI 使用指南

[English](./GUIDE.md) | 中文

## 安装

### Homebrew（macOS / Linux）

```bash
brew tap frankieew/tap
brew install lmfdb-cli
```

升级到最新版本：

```bash
brew update && brew upgrade lmfdb-cli
```

### 预编译版本

从 [GitHub Releases](https://github.com/FrankieeW/lmfdb-cli/releases) 下载，将二进制文件放入 `$PATH`。

### 从源码构建

需要 Go 1.24+。

```bash
git clone https://github.com/FrankieeW/lmfdb-cli.git
cd lmfdb-cli
go build -o lmfdb ./cmd/lmfdb
sudo mv lmfdb /usr/local/bin/
```

### 浏览器设置（可选）

部分 LMFDB 查询受 reCAPTCHA 保护。`--browser` 参数使用无头 Chrome 绕过验证：

```bash
lmfdb install-browser
```

这会下载 chromedp 管理的 Chromium，无需安装系统 Chrome。

## 命令

### 帮助

```bash
lmfdb -h                # 全局帮助
lmfdb nf -h             # 数域帮助
lmfdb ec -h             # 椭圆曲线帮助
lmfdb version            # 显示版本
```

### 数域（`nf`）

查询 LMFDB 的 `nf_fields` 集合。

#### 基本查询

```bash
# 二次域（degree 2，默认）
lmfdb nf

# 三次域
lmfdb nf -d 3

# 五次域，20 条结果
lmfdb nf -d 5 -n 20
```

#### 过滤

```bash
# 按判别式
lmfdb nf -d 2 --disc -5

# 按类数
lmfdb nf -d 2 --class 1

# 按 signature (r1, r2)
lmfdb nf -d 4 --sig 2,1

# 组合过滤
lmfdb nf -d 3 --class 1 -n 50
```

#### 按标签查询

```bash
# 通过 LMFDB 标签获取特定数域
lmfdb nf --id 2.0.3.1
lmfdb nf --id 3.1.23.1
```

#### 排序

```bash
# 按类数升序
lmfdb nf -d 2 --sort class_number

# 按判别式降序
lmfdb nf -d 2 --sort -disc_abs
```

#### 分页

```bash
# 第一页
lmfdb nf -d 2 -n 20

# 第二页
lmfdb nf -d 2 -n 20 --offset 20

# 第三页
lmfdb nf -d 2 -n 20 --offset 40
```

#### 选择字段

```bash
# 只返回指定字段
lmfdb nf -d 2 -f label,degree,disc_abs,class_number
```

### 椭圆曲线（`ec`）

查询 LMFDB 的 `ec_curvedata` 集合。

#### 基本查询

```bash
# 列出椭圆曲线
lmfdb ec -n 10

# 按秩筛选
lmfdb ec -r 0
lmfdb ec -r 2

# 按 torsion 阶筛选
lmfdb ec -t 5

# 按 conductor 筛选
lmfdb ec --conductor 11
```

#### 组合过滤

```bash
# 秩为 0，torsion 为 2
lmfdb ec -r 0 -t 2 -n 20

# 指定 conductor，排序
lmfdb ec --conductor 389 --sort -rank
```

### 列出集合（`list` / `ls`）

```bash
lmfdb list
```

显示所有可用的 LMFDB API 集合：

| 集合 | 说明 |
|------|------|
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

## 输出格式

### 表格（默认）

```bash
lmfdb nf -d 2 -n 5
```

在终端打印格式化表格。自动选择列（最多 6 列），长值会被截断。

### JSON

```bash
lmfdb nf -d 2 -n 5 --fmt json
```

输出带语法高亮的 JSON：
- **键名** 青色
- **字符串** 绿色
- **数字** 黄色
- **布尔值 / null** 紫色
- **括号** 灰色

### CSV

```bash
lmfdb nf -d 2 -n 5 --fmt csv
```

输出带表头的 CSV，方便管道传递给其他工具：

```bash
lmfdb nf -d 2 -n 100 --fmt csv -q | head -1    # 仅表头
lmfdb nf -d 2 -n 100 --fmt csv -q | wc -l       # 计数
```

### 保存到文件

```bash
# 保存为 JSON
lmfdb nf -d 2 -n 100 -o results.json

# 保存为 CSV
lmfdb nf -d 2 -n 100 -o results.csv --fmt csv
```

### 静默模式

使用 `-q` 抑制状态信息（适合管道操作）：

```bash
lmfdb nf -d 2 -n 10 --fmt json -q | jq '.[0].label'
```

## 绕过 reCAPTCHA

LMFDB 使用 Google reCAPTCHA 保护 API。如果请求被拦截，你会看到：

```
Error: Blocked by reCAPTCHA
Tip: Use --browser to bypass reCAPTCHA
```

添加 `--browser` 使用无头 Chrome：

```bash
lmfdb nf -d 2 -n 50 --browser
lmfdb ec -r 2 --browser --fmt json
```

首次请求可能需要几秒钟等待 Chrome 启动，之后数据提取速度很快。

**前提：** 先运行 `lmfdb install-browser` 下载 Chromium。

## 参数速查表

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--fmt` | 输出格式：`table`、`json`、`csv` | `table` |
| `-n` | 结果数量 | `10` |
| `--offset` | 结果偏移量（分页） | `0` |
| `--sort` | 排序字段（`-` 前缀降序） | |
| `-f` | 返回字段（逗号分隔） | 全部 |
| `-o` | 输出文件路径 | |
| `-q` | 静默模式 | `false` |
| `--browser` | 使用无头 Chrome（绕过 reCAPTCHA） | `false` |

### 数域专用

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-d` | 数域次数 | `2` |
| `--disc` | 按判别式筛选 | |
| `--class` | 按类数筛选 | |
| `--sig` | 按 signature 筛选（如 `0,1`） | |
| `--id` | 按 LMFDB 标签查询 | |

### 椭圆曲线专用

| 参数 | 说明 |
|------|------|
| `-r` | 按 Mordell-Weil 秩筛选 |
| `-t` | 按 torsion 阶筛选 |
| `--conductor` | 按 conductor 筛选 |

## 使用示例

### 研究工作流

```bash
# 导出所有类数为 1 的二次域
lmfdb nf -d 2 --class 1 -n 1000 -o class1_quadratic.json -q

# 查找高秩椭圆曲线
lmfdb ec -r 3 --sort -conductor --fmt json

# 三次域判别式导出为 CSV
lmfdb nf -d 3 -f label,disc_abs,class_number --fmt csv -q > cubic.csv

# 分页遍历结果
for i in 0 100 200 300; do
  lmfdb nf -d 2 -n 100 --offset $i -o "batch_$i.json" -q
done
```

### 配合 jq 使用

```bash
# 提取标签
lmfdb nf -d 2 -n 10 --fmt json -q | jq -r '.[].label'

# 按类数分组计数
lmfdb nf -d 2 -n 100 --fmt json -q | jq 'group_by(.class_number) | map({class_number: .[0].class_number, count: length})'

# 统计 CM 域
lmfdb nf -d 2 -n 100 --fmt json -q | jq '[.[] | select(.cm == true)] | length'
```

## 常见问题

### "Error: Blocked by reCAPTCHA"

使用 `--browser` 参数。如未安装浏览器，先运行 `lmfdb install-browser`。

### "Browser error: ..."

- 确保已安装 Chromium：`lmfdb install-browser`
- Linux 服务器可能需要：`apt install -y libnss3 libatk-bridge2.0-0 libcups2`

### 首次查询较慢

使用 `--browser` 的首次查询需要几秒启动 Chrome，这是正常现象。

### 没有结果

- 检查过滤条件 — 某些组合可能返回空结果
- 先不带过滤条件测试集合是否正常
- 使用 `-f` 选择更少字段，减小响应体积
