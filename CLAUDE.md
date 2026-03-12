# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

LMFDB CLI is a command-line tool for querying the [LMFDB](https://www.lmfdb.org/) (L-Functions and Modular Forms Database). The project has dual implementations:
- **Go**: Primary binary (cmd/lmfdb/main.go) - cross-platform standalone
- **Python**: Full-featured version (lmfdb_cli/) - uses Playwright to bypass reCAPTCHA

## Common Commands

### Python Development
```bash
# Install dependencies
pip install -e ".[dev]"

# Install Playwright browsers (required for reCAPTCHA bypass)
playwright install chromium

# Run tests
pytest

# Code linting
ruff check .
ruff format .
```

### Go Development
```bash
# Build binary
go build -o lmfdb ./cmd/lmfdb

# Run
./lmfdb nf -d 2 -n 10
```

### CLI Commands
```bash
# Number fields
lmfdb nf -d 2                    # quadratic fields
lmfdb nf -d 3                    # cubic fields
lmfdb nf --id 2.0.3.1            # get specific field

# Elliptic curves
lmfdb ec -n 10
lmfdb ec -r 2                    # filter by rank

# List collections
lmfdb list-collections
```

## Architecture

- **cmd/lmfdb/main.go**: Go implementation with flag-based CLI
- **lmfdb_cli/main.py**: Python CLI using Typer
- **lmfdb_cli/client.py**: HTTP client for LMFDB API
- **share/lmfdb_cli/**: Package data (shell completions)
- **.github/workflows/build.yml**: Builds Go binaries for linux/macOS/Windows (amd64/arm64)

The Python version includes Playwright-based reCAPTCHA bypass; the Go version makes direct API calls and may be blocked by reCAPTCHA for certain queries.

## Output Formats

Both implementations support `--fmt` (Go) or output options for:
- `table` - Default terminal table
- `json` - JSON output
- `csv` - CSV output
