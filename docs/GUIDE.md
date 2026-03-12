# LMFDB CLI Usage Guide

[English](./GUIDE.md) | [中文](./GUIDE.zh-CN.md)

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap frankieew/tap
brew install lmfdb-cli
```

Upgrade to the latest version:

```bash
brew update && brew upgrade lmfdb-cli
```

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/FrankieeW/lmfdb-cli/releases) and place the binary in your `$PATH`.

### From Source

Requires Go 1.24+.

```bash
git clone https://github.com/FrankieeW/lmfdb-cli.git
cd lmfdb-cli
go build -o lmfdb ./cmd/lmfdb
sudo mv lmfdb /usr/local/bin/
```

### Browser Setup (Optional)

Some LMFDB queries are protected by reCAPTCHA. The `--browser` flag uses headless Chrome to bypass this. To set it up:

```bash
lmfdb install-browser
```

This downloads a Chromium binary managed by chromedp. No system Chrome installation is required.

## Commands

### Help

```bash
lmfdb -h                # global help
lmfdb nf -h             # number fields help
lmfdb ec -h             # elliptic curves help
lmfdb version            # show version
```

### Number Fields (`nf`)

Query the `nf_fields` collection from LMFDB.

#### Basic Queries

```bash
# Quadratic fields (degree 2, default)
lmfdb nf

# Cubic fields
lmfdb nf -d 3

# Quintic fields, 20 results
lmfdb nf -d 5 -n 20
```

#### Filtering

```bash
# By discriminant
lmfdb nf -d 2 --disc -5

# By class number
lmfdb nf -d 2 --class 1

# By signature (r1, r2)
lmfdb nf -d 4 --sig 2,1

# Combine filters
lmfdb nf -d 3 --class 1 -n 50
```

#### Lookup by Label

```bash
# Get a specific number field by its LMFDB label
lmfdb nf --id 2.0.3.1
lmfdb nf --id 3.1.23.1
```

#### Sorting

```bash
# Sort by class number (ascending)
lmfdb nf -d 2 --sort class_number

# Sort by discriminant (descending)
lmfdb nf -d 2 --sort -disc_abs
```

#### Pagination

```bash
# First page
lmfdb nf -d 2 -n 20

# Second page
lmfdb nf -d 2 -n 20 --offset 20

# Third page
lmfdb nf -d 2 -n 20 --offset 40
```

#### Select Fields

```bash
# Only return specific fields
lmfdb nf -d 2 -f label,degree,disc_abs,class_number
```

### Elliptic Curves (`ec`)

Query the `ec_curvedata` collection from LMFDB.

#### Basic Queries

```bash
# List elliptic curves
lmfdb ec -n 10

# Filter by rank
lmfdb ec -r 0
lmfdb ec -r 2

# Filter by torsion order
lmfdb ec -t 5

# Filter by conductor
lmfdb ec --conductor 11
```

#### Combine Filters

```bash
# Rank 0, torsion 2
lmfdb ec -r 0 -t 2 -n 20

# Specific conductor, sorted
lmfdb ec --conductor 389 --sort -rank
```

### List Collections (`list` / `ls`)

```bash
lmfdb list
```

Shows all available LMFDB API collections:

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

## Output Formats

### Table (Default)

```bash
lmfdb nf -d 2 -n 5
```

Prints a formatted table to the terminal. Columns are automatically selected from the data (up to 6 columns). Long values are truncated.

### JSON

```bash
lmfdb nf -d 2 -n 5 --fmt json
```

Outputs syntax-highlighted JSON to the terminal:
- **Keys** in cyan
- **Strings** in green
- **Numbers** in yellow
- **Booleans / null** in magenta
- **Brackets** in gray

### CSV

```bash
lmfdb nf -d 2 -n 5 --fmt csv
```

Outputs CSV with a header row. Useful for piping to other tools:

```bash
lmfdb nf -d 2 -n 100 --fmt csv -q | head -1    # header only
lmfdb nf -d 2 -n 100 --fmt csv -q | wc -l       # count rows
```

### Save to File

```bash
# Save as JSON
lmfdb nf -d 2 -n 100 -o results.json

# Save as CSV
lmfdb nf -d 2 -n 100 -o results.csv --fmt csv
```

### Quiet Mode

Use `-q` to suppress status messages (useful for piping):

```bash
lmfdb nf -d 2 -n 10 --fmt json -q | jq '.[0].label'
```

## Bypassing reCAPTCHA

LMFDB protects its API with Google reCAPTCHA. If your request is blocked, you will see:

```
Error: Blocked by reCAPTCHA
Tip: Use --browser to bypass reCAPTCHA
```

Add `--browser` to use headless Chrome:

```bash
lmfdb nf -d 2 -n 50 --browser
lmfdb ec -r 2 --browser --fmt json
```

The first request may take a few seconds as Chrome starts up. Subsequent data extraction is fast.

**Requirements:** Run `lmfdb install-browser` once to download Chromium.

## Common Options Reference

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--fmt` | | Output format: `table`, `json`, `csv` | `table` |
| `-n` | | Number of results | `10` |
| `--offset` | | Result offset (pagination) | `0` |
| `--sort` | | Sort field (prefix `-` for descending) | |
| `-f` | | Fields to return (comma-separated) | all |
| `-o` | | Output file path | |
| `-q` | | Quiet mode (suppress status messages) | `false` |
| `--browser` | | Use headless Chrome (bypass reCAPTCHA) | `false` |

### Number Fields Only

| Flag | Description | Default |
|------|-------------|---------|
| `-d` | Number field degree | `2` |
| `--disc` | Filter by discriminant | |
| `--class` | Filter by class number | |
| `--sig` | Filter by signature (e.g., `0,1`) | |
| `--id` | Get specific field by LMFDB label | |

### Elliptic Curves Only

| Flag | Description |
|------|-------------|
| `-r` | Filter by Mordell-Weil rank |
| `-t` | Filter by torsion order |
| `--conductor` | Filter by conductor |

## Examples

### Research Workflows

```bash
# Export all quadratic fields with class number 1
lmfdb nf -d 2 --class 1 -n 1000 -o class1_quadratic.json -q

# Find high-rank elliptic curves
lmfdb ec -r 3 --sort -conductor --fmt json

# Get discriminants of cubic fields as CSV
lmfdb nf -d 3 -f label,disc_abs,class_number --fmt csv -q > cubic.csv

# Paginate through results
for i in 0 100 200 300; do
  lmfdb nf -d 2 -n 100 --offset $i -o "batch_$i.json" -q
done
```

### Piping with jq

```bash
# Extract labels
lmfdb nf -d 2 -n 10 --fmt json -q | jq -r '.[].label'

# Count fields by class number
lmfdb nf -d 2 -n 100 --fmt json -q | jq 'group_by(.class_number) | map({class_number: .[0].class_number, count: length})'

# Find CM fields
lmfdb nf -d 2 -n 100 --fmt json -q | jq '[.[] | select(.cm == true)] | length'
```

## Troubleshooting

### "Error: Blocked by reCAPTCHA"

Use `--browser` flag. If not installed, run `lmfdb install-browser` first.

### "Browser error: ..."

- Ensure Chromium is installed: `lmfdb install-browser`
- On Linux servers, you may need: `apt install -y libnss3 libatk-bridge2.0-0 libcups2`

### Slow first query

The first query with `--browser` takes a few seconds for Chrome startup. This is normal.

### No results found

- Check your filters — some combinations return empty results
- Try without filters first to verify the collection works
- Use `-f` to select fewer fields if the response is too large
