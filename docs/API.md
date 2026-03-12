# LMFDB API Reference

[English](./API.md) | [中文](./API.zh-CN.md)

This document describes the LMFDB API query syntax.

## Base URL

```
https://www.lmfdb.org/api/
```

## Query Format

### Type Prefixes

Values must be prefixed with their type:

| Prefix | Type | Example |
|--------|------|---------|
| `i` | Integer | `degree=i2` |
| `f` | Float | `rank=f2.5` |
| `s` | String | `label=s11a1` |
| `li` | List of integers | `torsion=li2;4` |
| `ls` | List of strings | `coeffs=ls1;-1;1` |
| `py` | Python literal | `matrix=py[[1,2],[3,4]]` |
| `cs` | Contains string | `label=cs11a` |
| `ci` | Contains integer | `torsion_structure=ci2` |

### Meta Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `_format` | Output format (json/yaml/html) | `_format=json` |
| `_fields` | Fields to return | `_fields=label,degree,disc` |
| `_sort` | Sort by field (prefix `-` for descending) | `_sort=-degree,label` |
| `_limit` | Maximum results | `_limit=100` |
| `_offset` | Result offset | `_offset=50` |
| `_delim` | List delimiter (default: `,`) | `_delim=;` |

## Examples

### Number Fields

```bash
# Degree 2 fields
https://www.lmfdb.org/api/nf_fields/?degree=i2&_format=json

# Specific discriminant
https://www.lmfdb.org/api/nf_fields/?disc=i-5&_format=json

# Class number = 1
https://www.lmfdb.org/api/nf_fields/?class_number=i1&_format=json

# Multiple conditions
https://www.lmfdb.org/api/nf_fields/?degree=i2&class_number=i1&_format=json

# Select specific fields
https://www.lmfdb.org/api/nf_fields/?degree=i2&_fields=label,degree,disc,class_number&_format=json
```

### Elliptic Curves

```bash
# Rank = 2
https://www.lmfdb.org/api/ec_curvedata/?rank=i2&_format=json

# Torsion = 5
https://www.lmfdb.org/api/ec_curvedata/?torsion=i5&_format=json

# Conductor = 11
https://www.lmfdb.org/api/ec_curvedata/?conductor=i11&_format=json

# Sorted by rank descending
https://www.lmfdb.org/api/ec_curvedata/?_sort=-rank&_format=json
```

### Using curl

```bash
# Query number fields
curl "https://www.lmfdb.org/api/nf_fields/?degree=i2&_format=json"

# With specific fields
curl "https://www.lmfdb.org/api/nf_fields/?degree=i2&_fields=label,disc&_limit=10&_format=json"
```

## Common Field Names

### nf_fields (Number Fields)

| Field | Description |
|-------|-------------|
| `label` | LMFDB label |
| `degree` | Field degree |
| `disc` | Discriminant |
| `class_number` | Class number |
| `class_group` | Class group structure |
| `signature` | Signature [r1, r2] |
| `coefficients` | Defining polynomial coefficients |
| `cm` | Has complex multiplication |

### ec_curvedata (Elliptic Curves)

| Field | Description |
|-------|-------------|
| `label` | LMFDB label (e.g., "11a1") |
| `conductor` | Conductor |
| `rank` | Mordell-Weil rank |
| `torsion` | Torsion order |
| `torsion_structure` | Torsion structure (e.g., [2,2]) |
| `equation` | Minimal equation |
| `j_invariant` | j-invariant |
| `cid` | Isogeny class ID |

## Limits

- Maximum results per query: **10,000**
- Default limit: **100**
- Use `_offset` for pagination

## Rate Limiting

Please use the API responsibly. LMFDB provides this service for research purposes.

## Error Messages

- `Unable to locate the page you requested` - Invalid collection or field name
- `recaptcha` response - Request blocked, try again or use browser
- HTTP 404 - Invalid URL
- HTTP 429 - Too many requests
