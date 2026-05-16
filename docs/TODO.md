# TODO — lmfdb-cli 改进计划

来源：对照 LMFDB 官方仓库（`LMFDB/lmfdb`、`LMFDB/lmfdb-inventory`）和现有客户端
（`roed314/lmfdb-lite`、`sage2lmfdb`、`CBirkbeck/LeanBridge`）后整理出的改进空间。
顺序按 ROI（价值 / 工作量）排列，不代表必须按序执行。

---

## P0 — 高价值、低成本

### 1. 修正文档漂移（Go → Rust）
- `CLAUDE.md` 里仍写 "Pure Go implementation with chromedp"，实际已是 Rust + reqwest
- `README.md`、`README.zh-CN.md`、`docs/GUIDE*.md` 需要核对命令示例
- 历史包袱：`go.mod`、`go.sum`、`cmd/`、`bin/`、`lmfdb_cli.egg-info/` 应删除
- 验收：仓库内不再有 Go / Python 相关文件；`CLAUDE.md` 描述与 `Cargo.toml` 一致

### 2. 大查询自动分页（bug fix）
- LMFDB API 单次硬上限 100 行（`api.py` 中固定），`_limit` 参数被钳制到 0–10000
- 当前 `-n 500` 会静默只返回 100 条
- 方案：`fetch` 内部循环 `_offset += 100` 直到拿满 `limit` 或上游返回空
- 验收：`lmfdb nf -d 2 -n 250` 真的返回 250 条；新增单测覆盖 offset 累加

### 3. 通用 `query` 子命令（覆盖率从 2 → 全部 28 个 collection）
- 形态：`lmfdb query <table> --where field=op:val --fields ... --limit ...`
- `<table>` 例：`g2c_curves`、`gps_groups`、`lf_fields`、`mf_newforms`、
  `hmf_forms`、`lat_lattices`、`modcurve_*`、`artin_*`、`hgm_*`、`smf_samples` …
- 让用户用一个统一入口操作 inventory 中列出的所有表
- 现有 `nf` / `ec` 保留为便捷封装
- 验收：可成功查询至少 5 个原本没专用子命令的 table

### 4. 完整类型前缀支持
- 现状：只用了 `i`、`li`
- 需补：`s`（string）、`f`（float）、`ls`、`lf`、`py`（Python 字面量，矩阵嵌套列表）、
  `cs` / `ci` / `cf` / `cpy`（`$contains` 包含查询）
- CLI 表达：`--where conductor=i:11`、`--where bad_primes=ci:5`
- 验收：每种前缀有 URL builder 单测

### 5. Shell 补全
- 引入 `clap_complete`，加 `lmfdb completions <shell>` 子命令
- 输出 bash / zsh / fish / powershell 补全脚本
- README 写明安装方式
- 验收：`lmfdb completions fish | source` 后命令、子命令、长选项都能补全

---

## P1 — 中等价值

### 6. 解决 `--browser` 死分支
两个选项二选一：
- (a) 用 `chromiumoxide` 在 `--features browser` 下真正实现 reCAPTCHA 绕过
- (b) 直接删 `--browser` 旗标和 `install-browser` 子命令，避免运行时 `bail!` 误导用户
- 验收：调用 `--browser` 不会再 panic / bail，行为与 `--help` 一致

### 7. 磁盘缓存层
- LMFDB 数据基本静态，重复查询浪费
- 路径：`~/.cache/lmfdb-cli/<url-hash>.json`，TTL 默认 7 天
- 旗标：`--no-cache` 跳过；`--refresh` 强制刷新
- 验收：第二次同样查询 < 50ms；`--no-cache` 走网络

### 8. 重试 + 限速
- 当前一次 HTTP 失败立即 bail；上游偶发 5xx 体验差
- 指数退避（base 500ms，max 3 次），429/5xx 重试
- 全局 `--retries N`、`--timeout SECS`
- 验收：模拟 503 时观察到重试日志，最终成功或在 N 次后报错

### 9. `describe` 子命令（接 lmfdb-inventory）
- `lmfdb describe <table>` 列出该表所有字段、类型、含义
- 数据源：构建期拉取 `lmfdb-inventory` 仓库 YAML/JSON，嵌入二进制
- 这是 `lmfdb-lite` 和官方 web UI 都没做的差异化
- 验收：`lmfdb describe nf_fields` 显示字段表

---

## P2 — 长尾

### 10. `--open` 旗标（人类页面）
- 输出记录时同时支持 `--open` 在浏览器中打开 `lmfdb.org/NumberField/2.0.3.1`
- 衔接 CLI 探索 → 网页深读的常见工作流
- 验收：macOS / Linux / Windows 三平台各自调用 `open` / `xdg-open` / `start`

### 11. L-function 端点
- `lfunc_*` 是 LMFDB 招牌之一，目前未覆盖
- 至少支持 `lmfdb query lfunc_instances` 与按 label 查询
- 验收：`lmfdb query lfunc_lfunctions --where degree=i:2 -n 5` 工作正常

### 12. Sage / Lean 互操作输出
- `--fmt sage`：输出可在 Sage notebook 直接 `eval` 的 Python dict（含 `NumberField(...)` 之类）
- `--fmt lean`：输出 Lean 4 可读的结构体，参考 `CBirkbeck/LeanBridge`
- 验收：示例输出能被 Sage / Lean 直接解析（在 README 贴 round-trip 例子）

---

## 跨项任务（贯穿全部）

- 单元测试覆盖：每个新增 URL builder / 分页 / 重试逻辑都要补 `#[cfg(test)]`
- CI：保留现有 `cargo clippy --all-targets -- -D warnings` 门槛
- README & docs/GUIDE 双语同步更新（每个新特性 EN + ZH）

---

## 建议起步顺序

1. 先做 **#1（删 Go 残骸 + 修 CLAUDE.md）**：10 分钟，避免后续工作被误导
2. 再做 **#2（自动分页）**：bug fix，独立 PR
3. 然后 **#3 + #4 一起做**：通用 `query` 命令需要完整的类型前缀支持
4. **#5（补全）** 可在任何时点穿插
