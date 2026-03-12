# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

LMFDB CLI is a command-line tool for querying the [LMFDB](https://www.lmfdb.org/) (L-Functions and Modular Forms Database). Pure Go implementation with chromedp for reCAPTCHA bypass.

## Common Commands

```bash
# Build binary
go build -o lmfdb ./cmd/lmfdb

# Run
./lmfdb nf -d 2 -n 10

# Number fields
lmfdb nf -d 2                    # quadratic fields
lmfdb nf -d 3                    # cubic fields
lmfdb nf --id 2.0.3.1            # get specific field
lmfdb nf --browser -n 50         # bypass reCAPTCHA

# Elliptic curves
lmfdb ec -n 10
lmfdb ec -r 2                    # filter by rank

# List collections
lmfdb list
```

## Architecture

- **cmd/lmfdb/main.go**: Single-file CLI with flag-based commands
- **.github/workflows/build.yml**: Builds Go binaries for linux/macOS/Windows (amd64/arm64)

Uses chromedp for headless browser reCAPTCHA bypass via `--browser` flag.

## Output Formats

Via `--fmt` flag:
- `table` - Default terminal table
- `json` - JSON output
- `csv` - CSV output
