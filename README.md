# LMFDB CLI

[English](./README.md) | [中文](./README.zh-CN.md)

A command-line tool for querying the **LMFDB** (L-Functions and Modular Forms Database).

## Features

- ✅ Bypasses reCAPTCHA using Playwright
- ✅ Support all LMFDB API collections
- ✅ Beautiful table output with Rich
- ✅ Customizable fields
- ✅ JSON/YAML output format
- ✅ Filter and sort support

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/lmfdb-cli.git
cd lmfdb-cli

# Install dependencies
pip install -e .

# Install Playwright browsers
playwright install chromium
```

## Quick Start

```bash
# Query quadratic fields (default: degree=2)
lmfdb nf -d 2 -n 10

# Query elliptic curves
lmfdb ec -n 10

# List available collections
lmfdb list-collections

# Generic query
lmfdb query nf_fields -n 5
```

## Commands

### nf - Number Fields

Query number fields from LMFDB.

```bash
# Query quadratic fields
lmfdb nf -d 2

# Query cubic fields
lmfdb nf -d 3

# Filter by discriminant
lmfdb nf -d 2 --discriminant -5

# Limit results
lmfdb nf -d 2 -n 20

# Output to JSON file
lmfdb nf -d 2 -n 10 -o results.json
```

Options:
- `-d, --degree`: Number field degree (default: 2)
- `--disc`: Filter by discriminant
- `-h, --class-number`: Filter by class number
- `-n, --limit`: Number of results (default: 10)
- `-f, --fields`: Comma-separated fields to return
- `-o, --output`: Output file path (JSON)
- `--headless/--no-headless`: Browser mode

### ec - Elliptic Curves

Query elliptic curves from LMFDB.

```bash
# Query all elliptic curves
lmfdb ec -n 10

# Filter by rank
lmfdb ec -r 2

# Filter by torsion
lmfdb ec -t 5

# Filter by conductor
lmfdb ec --conductor 11

# Custom fields
lmfdb ec -n 10 -f label,conductor,rank,torsion
```

Options:
- `-r, --rank`: Filter by Mordell-Weil rank
- `-t, --torsion`: Filter by torsion order
- `--conductor`: Filter by conductor
- `-n, --limit`: Number of results (default: 10)
- `-f, --fields`: Comma-separated fields to return
- `-o, --output`: Output file path (JSON)

### query - Generic Query

Query any LMFDB API collection.

```bash
# List collections
lmfdb list-collections

# Query specific collection
lmfdb query nf_fields -n 5
lmfdb query ec_curvedata -n 10

# With custom parameters (JSON)
lmfdb query nf_fields -k '{"degree": "i3"}' -n 10
```

## Available Collections

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

## Configuration

### Environment Variables

```bash
# Set timeout (milliseconds)
export LMFDB_TIMEOUT=60000

# Set headless mode
export LMFDB_HEADLESS=true
```

## Development

```bash
# Install development dependencies
pip install -e ".[dev]"

# Run tests
pytest

# Code formatting
ruff check .
ruff format .
```

## License

MIT

## Credits

- [LMFDB](https://www.lmfdb.org/) - The L-Functions and Modular Forms Database
- [Playwright](https://playwright.dev/) - Browser automation
- [Rich](https://rich.readthedocs.io/) - Terminal formatting
