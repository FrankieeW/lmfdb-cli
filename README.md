# LMFDB CLI

[English](./README.md) | [中文](./README.zh-CN.md)

A command-line tool for querying the **LMFDB** (L-Functions and Modular Forms Database).

## Features

- Pure Go implementation
- Bypasses reCAPTCHA using chromedp (with `--browser` flag)
- Support all LMFDB API collections
- Table/JSON/CSV output formats
- Filter and sort support
- Cross-platform (Linux, macOS, Windows)

## Installation

### Homebrew (macOS/Linux)
```bash
brew tap frankieew/tap
brew install lmfdb-cli

# Install browser for reCAPTCHA bypass
lmfdb install-browser
```

### From Source
```bash
git clone https://github.com/FrankieeW/lmfdb-cli.git
cd lmfdb-cli
go build -o lmfdb ./cmd/lmfdb
./lmfdb install-browser  # Optional: for reCAPTCHA bypass
```

### From Release
Download pre-built binaries from [GitHub Releases](https://github.com/FrankieeW/lmfdb-cli/releases).

## Quick Start

```bash
# Query quadratic fields (default: degree=2)
lmfdb nf -d 2 -n 10

# Query elliptic curves
lmfdb ec -n 10

# Use browser to bypass reCAPTCHA
lmfdb nf -d 2 -n 10 --browser

# List available collections
lmfdb list
```

## Commands

### nf - Number Fields

```bash
lmfdb nf -d 2              # quadratic fields
lmfdb nf -d 3              # cubic fields
lmfdb nf -d 2 --disc -5    # filter by discriminant
lmfdb nf -d 2 -n 20        # limit results
lmfdb nf --id 2.0.3.1      # specific field by label
lmfdb nf -d 2 --fmt json   # JSON output
lmfdb nf --browser -n 50   # use browser to bypass reCAPTCHA
```

Options:
- `-d` — Number field degree (default: 2)
- `--disc` — Filter by discriminant
- `--class` — Filter by class number
- `--sig` — Filter by signature (e.g., `0,1`)
- `-n` — Number of results (default: 10)
- `--offset` — Result offset for pagination
- `--sort` — Sort by field (prefix `-` for descending)
- `-f` — Fields to return (comma-separated)
- `-o` — Output file
- `--fmt` — Output format: `table`, `json`, `csv`
- `--id` — Get specific field by label
- `-q` — Quiet mode
- `--browser` — Use browser (bypasses reCAPTCHA)

### ec - Elliptic Curves

```bash
lmfdb ec -n 10              # list curves
lmfdb ec -r 2               # filter by rank
lmfdb ec -t 5               # filter by torsion
lmfdb ec --conductor 11     # filter by conductor
```

Options:
- `-r` — Filter by Mordell-Weil rank
- `-t` — Filter by torsion order
- `--conductor` — Filter by conductor
- `-n`, `--offset`, `--sort`, `-f`, `-o`, `--fmt`, `-q`, `--browser` — same as `nf`

### list (ls) - Available Collections

```bash
lmfdb list
```

| Collection | Description |
|------------|-------------|
| `nf_fields` | Number fields |
| `ec_curvedata` | Elliptic curves |
| `ec_classdata` | Elliptic curve isogeny classes |
| `g2c_curves` | Genus 2 curves |
| `char_dirichlet` | Dirichlet characters |
| `maass_newforms` | Maass forms |
| `mf_newforms` | Modular forms |
| `lf_fields` | Local fields |
| `artin` | Artin representations |
| `belyi` | Belyi maps |

## License

MIT

## Credits

- [LMFDB](https://www.lmfdb.org/) - The L-Functions and Modular Forms Database
