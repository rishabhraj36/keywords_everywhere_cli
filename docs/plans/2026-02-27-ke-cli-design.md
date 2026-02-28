# Keywords Everywhere CLI Design

## Overview

A Go CLI (`ke`) wrapping the Keywords Everywhere API with JSON output for skill/automation use.

## Authentication

Environment variable: `KEYWORDS_EVERYWHERE_API_KEY`

## Commands

### Miscellaneous

| Command | Description |
|---------|-------------|
| `ke credits` | Check credit balance |
| `ke countries` | List supported countries |
| `ke currencies` | List supported currencies |

### Keyword Data

| Command | Description | Credits |
|---------|-------------|---------|
| `ke keywords <kw>...` | Volume/CPC/competition | 1/keyword |
| `ke related <keyword>` | Related keywords | 2/keyword |
| `ke pasf <keyword>` | People Also Search For | 2/keyword |

### Traffic Metrics

| Command | Description |
|---------|-------------|
| `ke domain-keywords <domain>` | Keywords a domain ranks for |
| `ke url-keywords <url>` | Keywords a URL ranks for |
| `ke domain-traffic <domain>` | Traffic metrics for domain |
| `ke url-traffic <url>` | Traffic metrics for URL |

### Backlinks

| Command | Description |
|---------|-------------|
| `ke domain-backlinks <domain>` | Backlinks for domain |
| `ke page-backlinks <url>` | Backlinks for page |

## Global Flags

```
--country, -c    Country code (default: us)
--currency       Currency code (default: usd)
--source, -s     Data source: gkp|cli (default: gkp)
--limit, -l      Max results for list endpoints
```

## Input

- **Keywords**: args or stdin (one per line)
- **Domains/URLs**: single arg
- Auto-batches keywords (100 per API call)

## Output

JSON to stdout. Errors to stderr with non-zero exit.

```json
{
  "credits": 1,
  "data": [...]
}
```

## Project Structure

```
keywords-everywhere-cli/
├── main.go
├── go.mod
├── cmd/
│   ├── root.go
│   ├── credits.go
│   ├── keywords.go
│   ├── related.go
│   ├── pasf.go
│   ├── traffic.go
│   └── backlinks.go
└── internal/
    └── api/
        └── client.go
```

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- Standard library for HTTP/JSON

## API Reference

Base URL: `https://api.keywordseverywhere.com/v1/`

Auth: `Authorization: Bearer <API_KEY>`

### Endpoints

| Endpoint | Method | Notes |
|----------|--------|-------|
| `/account/credits` | GET | Credit balance |
| `/countries` | GET | Supported countries |
| `/currencies` | GET | Supported currencies |
| `/get_keyword_data` | POST | Max 100 keywords |
| `/get_related_keywords` | POST | |
| `/get_pasf_keywords` | POST | Gold/Premium only |
| `/get_domain_keywords` | POST | |
| `/get_url_keywords` | POST | |
| `/get_domain_traffic` | POST | |
| `/get_url_traffic` | POST | |
| `/get_domain_backlinks` | POST | |
| `/get_page_backlinks` | POST | |
