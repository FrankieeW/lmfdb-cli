# LMFDB CLI

[English](./README.md) | [中文](./README.zh-CN.md)

A command-line tool for querying the [LMFDB](https://www.lmfdb.org/) (L-Functions and Modular Forms Database).

## Features

- Pure Go, single binary, no dependencies
- Bypasses reCAPTCHA via headless Chrome (`--browser`)
- Syntax-highlighted JSON output
- Table / JSON / CSV output formats
- Flexible filtering, sorting, and pagination
- Cross-platform (Linux, macOS, Windows)

## Installation

### Homebrew (macOS / Linux)
```bash
brew tap frankieew/tap
brew install lmfdb-cli
```

### From Release
Download pre-built binaries from [GitHub Releases](https://github.com/FrankieeW/lmfdb-cli/releases).

### From Source
```bash
git clone https://github.com/FrankieeW/lmfdb-cli.git
cd lmfdb-cli
go build -o lmfdb ./cmd/lmfdb
```

## Quick Start

```bash
lmfdb nf -d 2 -n 10               # quadratic number fields
lmfdb ec -r 2 --fmt json           # rank-2 elliptic curves, JSON output
lmfdb nf --id 2.0.3.1              # lookup by label
lmfdb nf --browser -n 50           # bypass reCAPTCHA
lmfdb list                         # available collections
```

## Commands

| Command | Description |
|---------|-------------|
| `nf` | Query number fields |
| `ec` | Query elliptic curves |
| `list` / `ls` | List available API collections |
| `version` / `v` | Show version |
| `install-browser` | Install Chrome for reCAPTCHA bypass |

Use `lmfdb <command> -h` for detailed help on each command.

## Output Formats

```bash
lmfdb nf -d 2 -n 5                 # table (default)
lmfdb nf -d 2 -n 5 --fmt json      # syntax-highlighted JSON
lmfdb nf -d 2 -n 5 --fmt csv       # CSV
lmfdb nf -d 2 -n 5 -o out.json     # save to file
```

## Documentation

- [Usage Guide](./docs/GUIDE.md) - Detailed usage with examples
- [API Reference](./docs/API.md) - LMFDB API query syntax
- [Documentation Index](./docs/INDEX.md)

## License

MIT

## Credits

- [LMFDB](https://www.lmfdb.org/) - The L-Functions and Modular Forms Database
