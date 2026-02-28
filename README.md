# Keywords Everywhere CLI

A Go CLI wrapper for the [Keywords Everywhere API](https://api.keywordseverywhere.com/docs/). Outputs JSON for easy parsing in scripts and automation.

## Installation

```bash
go install github.com/mrm/keywords-everywhere-cli@latest
```

Or build from source:

```bash
git clone https://github.com/mrm/keywords-everywhere-cli
cd keywords-everywhere-cli
go install .
```

## Setup

Set your API key as an environment variable:

```bash
export KEYWORDS_EVERYWHERE_API_KEY="your-api-key"
```

Get your API key at [keywordseverywhere.com](https://keywordseverywhere.com).

## Commands

### Keyword Data

```bash
# Get volume, CPC, competition for keywords
ke keywords "seo tools" "keyword research"

# Pipe keywords from stdin
echo -e "seo\nmarketing\ncontent" | ke keywords

# Get related keywords
ke related "seo"

# Get "People Also Search For" keywords
ke pasf "seo"
```

### Traffic Metrics

```bash
# Keywords a domain ranks for
ke domain-keywords example.com

# Keywords a URL ranks for
ke url-keywords https://example.com/page

# Traffic metrics for domain
ke domain-traffic example.com

# Traffic metrics for URL
ke url-traffic https://example.com/page
```

### Backlinks

```bash
# Backlinks for a domain
ke domain-backlinks example.com

# Backlinks for a specific page
ke page-backlinks https://example.com/page
```

### Account

```bash
# Check credit balance
ke credits

# List supported countries
ke countries

# List supported currencies
ke currencies
```

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--country` | `-c` | `us` | Country code |
| `--currency` | | `usd` | Currency code |
| `--source` | `-s` | `gkp` | Data source: `gkp` (Google Keyword Planner) or `cli` (Clickstream) |
| `--limit` | `-l` | `0` | Max results (0 = no limit) |

Example with flags:

```bash
ke keywords -c uk --currency gbp "british tea"
ke domain-keywords -l 100 example.com
```

## Output

All commands output JSON to stdout:

```bash
ke keywords "seo" | jq '.data[0]'
```

Errors go to stderr with non-zero exit code.

## Credits

API calls consume credits from your Keywords Everywhere account:
- `keywords`: 1 credit per keyword
- `related`, `pasf`: 2 credits per keyword
- Other endpoints: varies (check API docs)

The `keywords` command automatically batches requests (100 keywords per API call).

## License

MIT
