# Keywords Everywhere CLI Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go CLI (`ke`) that wraps the Keywords Everywhere API with JSON output for skill/automation use.

**Architecture:** Cobra-based CLI with subcommands. API client in `internal/api/` handles auth, requests, and response parsing. Each command is a thin wrapper calling the API client.

**Tech Stack:** Go 1.25, Cobra CLI, standard library HTTP/JSON.

---

## Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `main.go`
- Create: `cmd/root.go`

**Step 1: Initialize Go module**

Run:
```bash
go mod init github.com/mrm/keywords-everywhere-cli
```

**Step 2: Create main.go**

Create `main.go`:
```go
package main

import "github.com/mrm/keywords-everywhere-cli/cmd"

func main() {
	cmd.Execute()
}
```

**Step 3: Create root command**

Create `cmd/root.go`:
```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	country  string
	currency string
	source   string
	limit    int
)

var rootCmd = &cobra.Command{
	Use:   "ke",
	Short: "Keywords Everywhere CLI",
	Long:  "A CLI for the Keywords Everywhere API. Set KEYWORDS_EVERYWHERE_API_KEY environment variable.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&country, "country", "c", "us", "Country code")
	rootCmd.PersistentFlags().StringVar(&currency, "currency", "usd", "Currency code")
	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "gkp", "Data source: gkp|cli")
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "Max results (0 = no limit)")
}
```

**Step 4: Add Cobra dependency**

Run:
```bash
go get github.com/spf13/cobra@latest
go mod tidy
```

**Step 5: Verify it builds**

Run:
```bash
go build -o ke .
./ke --help
```

Expected: Shows help with "Keywords Everywhere CLI" and global flags.

**Step 6: Commit**

```bash
git add go.mod go.sum main.go cmd/root.go
git commit -m "feat: scaffold CLI with Cobra and global flags"
```

---

## Task 2: API Client Foundation

**Files:**
- Create: `internal/api/client.go`
- Create: `internal/api/client_test.go`

**Step 1: Write failing test for API client creation**

Create `internal/api/client_test.go`:
```go
package api

import (
	"os"
	"testing"
)

func TestNewClient_WithAPIKey(t *testing.T) {
	os.Setenv("KEYWORDS_EVERYWHERE_API_KEY", "test-key")
	defer os.Unsetenv("KEYWORDS_EVERYWHERE_API_KEY")

	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
}

func TestNewClient_MissingAPIKey(t *testing.T) {
	os.Unsetenv("KEYWORDS_EVERYWHERE_API_KEY")

	_, err := NewClient()
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}
```

**Step 2: Run test to verify it fails**

Run:
```bash
go test ./internal/api/... -v
```

Expected: FAIL - package doesn't exist yet.

**Step 3: Implement API client**

Create `internal/api/client.go`:
```go
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const baseURL = "https://api.keywordseverywhere.com/v1"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	apiKey := os.Getenv("KEYWORDS_EVERYWHERE_API_KEY")
	if apiKey == "" {
		return nil, errors.New("KEYWORDS_EVERYWHERE_API_KEY environment variable not set")
	}
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

func (c *Client) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *Client) post(endpoint string, data url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", baseURL+endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Response wraps API responses with credits used
type Response struct {
	Credits int             `json:"credits"`
	Data    json.RawMessage `json:"data"`
}
```

**Step 4: Run tests**

Run:
```bash
go test ./internal/api/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/client.go internal/api/client_test.go
git commit -m "feat: add API client with auth handling"
```

---

## Task 3: Credits Command

**Files:**
- Create: `cmd/credits.go`
- Create: `cmd/credits_test.go`

**Step 1: Write failing test**

Create `cmd/credits_test.go`:
```go
package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestCreditsCommand_Output(t *testing.T) {
	// Test that the command exists and has correct structure
	if creditsCmd.Use != "credits" {
		t.Errorf("expected Use 'credits', got %s", creditsCmd.Use)
	}
}
```

**Step 2: Run test to verify it fails**

Run:
```bash
go test ./cmd/... -v
```

Expected: FAIL - creditsCmd undefined.

**Step 3: Implement credits command**

Create `cmd/credits.go`:
```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var creditsCmd = &cobra.Command{
	Use:   "credits",
	Short: "Check account credit balance",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetCredits()
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(creditsCmd)
}
```

**Step 4: Add GetCredits to API client**

Add to `internal/api/client.go`:
```go
type CreditsResponse struct {
	Credits int `json:"credits"`
}

func (c *Client) GetCredits() (*CreditsResponse, error) {
	body, err := c.get("/account/credits")
	if err != nil {
		return nil, err
	}

	var resp CreditsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
```

**Step 5: Run tests**

Run:
```bash
go test ./... -v
go build -o ke .
```

Expected: PASS, builds successfully.

**Step 6: Commit**

```bash
git add cmd/credits.go cmd/credits_test.go internal/api/client.go
git commit -m "feat: add credits command"
```

---

## Task 4: Countries and Currencies Commands

**Files:**
- Create: `cmd/countries.go`
- Create: `cmd/currencies.go`

**Step 1: Implement countries command**

Create `cmd/countries.go`:
```go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var countriesCmd = &cobra.Command{
	Use:   "countries",
	Short: "List supported countries",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetCountries()
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(countriesCmd)
}
```

**Step 2: Implement currencies command**

Create `cmd/currencies.go`:
```go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var currenciesCmd = &cobra.Command{
	Use:   "currencies",
	Short: "List supported currencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetCurrencies()
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currenciesCmd)
}
```

**Step 3: Add API methods**

Add to `internal/api/client.go`:
```go
func (c *Client) GetCountries() (json.RawMessage, error) {
	return c.get("/countries")
}

func (c *Client) GetCurrencies() (json.RawMessage, error) {
	return c.get("/currencies")
}
```

**Step 4: Build and verify**

Run:
```bash
go build -o ke .
./ke countries --help
./ke currencies --help
```

Expected: Help shown for both commands.

**Step 5: Commit**

```bash
git add cmd/countries.go cmd/currencies.go internal/api/client.go
git commit -m "feat: add countries and currencies commands"
```

---

## Task 5: Keywords Command

**Files:**
- Create: `cmd/keywords.go`
- Modify: `internal/api/client.go`

**Step 1: Implement keywords command with stdin support**

Create `cmd/keywords.go`:
```go
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var keywordsCmd = &cobra.Command{
	Use:   "keywords [keyword...]",
	Short: "Get volume, CPC, and competition data for keywords",
	Long:  "Get keyword data. Pass keywords as args or pipe via stdin (one per line).",
	RunE: func(cmd *cobra.Command, args []string) error {
		keywords := args

		// Read from stdin if no args and stdin has data
		if len(keywords) == 0 {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					line := scanner.Text()
					if line != "" {
						keywords = append(keywords, line)
					}
				}
			}
		}

		if len(keywords) == 0 {
			return fmt.Errorf("no keywords provided")
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetKeywordData(keywords, country, currency, source)
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(keywordsCmd)
}
```

**Step 2: Add GetKeywordData to API client with batching**

Add to `internal/api/client.go`:
```go
type KeywordDataResponse struct {
	Credits int               `json:"credits"`
	Data    []json.RawMessage `json:"data"`
}

func (c *Client) GetKeywordData(keywords []string, country, currency, source string) (*KeywordDataResponse, error) {
	const batchSize = 100
	var allData []json.RawMessage
	totalCredits := 0

	for i := 0; i < len(keywords); i += batchSize {
		end := i + batchSize
		if end > len(keywords) {
			end = len(keywords)
		}
		batch := keywords[i:end]

		data := url.Values{}
		data.Set("country", country)
		data.Set("currency", currency)
		data.Set("dataSource", source)
		for _, kw := range batch {
			data.Add("kw[]", kw)
		}

		body, err := c.post("/get_keyword_data", data)
		if err != nil {
			return nil, err
		}

		var resp struct {
			Credits int               `json:"credits"`
			Data    []json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}

		totalCredits += resp.Credits
		allData = append(allData, resp.Data...)
	}

	return &KeywordDataResponse{
		Credits: totalCredits,
		Data:    allData,
	}, nil
}
```

**Step 3: Build and verify**

Run:
```bash
go build -o ke .
./ke keywords --help
```

Expected: Shows help with usage `keywords [keyword...]`.

**Step 4: Commit**

```bash
git add cmd/keywords.go internal/api/client.go
git commit -m "feat: add keywords command with batching and stdin support"
```

---

## Task 6: Related and PASF Commands

**Files:**
- Create: `cmd/related.go`
- Create: `cmd/pasf.go`

**Step 1: Implement related command**

Create `cmd/related.go`:
```go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var relatedCmd = &cobra.Command{
	Use:   "related <keyword>",
	Short: "Get related keywords",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetRelatedKeywords(args[0], country, currency, source)
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(relatedCmd)
}
```

**Step 2: Implement PASF command**

Create `cmd/pasf.go`:
```go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var pasfCmd = &cobra.Command{
	Use:   "pasf <keyword>",
	Short: "Get 'People Also Search For' keywords",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetPASFKeywords(args[0], country, currency, source)
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pasfCmd)
}
```

**Step 3: Add API methods**

Add to `internal/api/client.go`:
```go
func (c *Client) GetRelatedKeywords(keyword, country, currency, source string) (json.RawMessage, error) {
	data := url.Values{}
	data.Set("country", country)
	data.Set("currency", currency)
	data.Set("dataSource", source)
	data.Set("kw", keyword)

	return c.post("/get_related_keywords", data)
}

func (c *Client) GetPASFKeywords(keyword, country, currency, source string) (json.RawMessage, error) {
	data := url.Values{}
	data.Set("country", country)
	data.Set("currency", currency)
	data.Set("dataSource", source)
	data.Set("kw", keyword)

	return c.post("/get_pasf_keywords", data)
}
```

**Step 4: Build and verify**

Run:
```bash
go build -o ke .
./ke related --help
./ke pasf --help
```

Expected: Both commands show help.

**Step 5: Commit**

```bash
git add cmd/related.go cmd/pasf.go internal/api/client.go
git commit -m "feat: add related and pasf keyword commands"
```

---

## Task 7: Traffic Commands

**Files:**
- Create: `cmd/traffic.go`

**Step 1: Implement traffic commands**

Create `cmd/traffic.go`:
```go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var domainKeywordsCmd = &cobra.Command{
	Use:   "domain-keywords <domain>",
	Short: "Get keywords a domain ranks for",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetDomainKeywords(args[0], country, currency, limit)
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

var urlKeywordsCmd = &cobra.Command{
	Use:   "url-keywords <url>",
	Short: "Get keywords a URL ranks for",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetURLKeywords(args[0], country, currency, limit)
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

var domainTrafficCmd = &cobra.Command{
	Use:   "domain-traffic <domain>",
	Short: "Get traffic metrics for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetDomainTraffic(args[0])
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

var urlTrafficCmd = &cobra.Command{
	Use:   "url-traffic <url>",
	Short: "Get traffic metrics for a URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetURLTraffic(args[0])
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(domainKeywordsCmd)
	rootCmd.AddCommand(urlKeywordsCmd)
	rootCmd.AddCommand(domainTrafficCmd)
	rootCmd.AddCommand(urlTrafficCmd)
}
```

**Step 2: Add API methods**

Add to `internal/api/client.go`:
```go
func (c *Client) GetDomainKeywords(domain, country, currency string, limit int) (json.RawMessage, error) {
	data := url.Values{}
	data.Set("domain", domain)
	data.Set("country", country)
	data.Set("currency", currency)
	if limit > 0 {
		data.Set("limit", fmt.Sprintf("%d", limit))
	}

	return c.post("/get_domain_keywords", data)
}

func (c *Client) GetURLKeywords(targetURL, country, currency string, limit int) (json.RawMessage, error) {
	data := url.Values{}
	data.Set("url", targetURL)
	data.Set("country", country)
	data.Set("currency", currency)
	if limit > 0 {
		data.Set("limit", fmt.Sprintf("%d", limit))
	}

	return c.post("/get_url_keywords", data)
}

func (c *Client) GetDomainTraffic(domain string) (json.RawMessage, error) {
	data := url.Values{}
	data.Set("domain", domain)

	return c.post("/get_domain_traffic", data)
}

func (c *Client) GetURLTraffic(targetURL string) (json.RawMessage, error) {
	data := url.Values{}
	data.Set("url", targetURL)

	return c.post("/get_url_traffic", data)
}
```

**Step 3: Build and verify**

Run:
```bash
go build -o ke .
./ke --help
```

Expected: Shows all traffic commands in help.

**Step 4: Commit**

```bash
git add cmd/traffic.go internal/api/client.go
git commit -m "feat: add traffic and domain/URL keyword commands"
```

---

## Task 8: Backlinks Commands

**Files:**
- Create: `cmd/backlinks.go`

**Step 1: Implement backlinks commands**

Create `cmd/backlinks.go`:
```go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var domainBacklinksCmd = &cobra.Command{
	Use:   "domain-backlinks <domain>",
	Short: "Get backlinks for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetDomainBacklinks(args[0], limit)
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

var pageBacklinksCmd = &cobra.Command{
	Use:   "page-backlinks <url>",
	Short: "Get backlinks for a specific page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetPageBacklinks(args[0], limit)
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(domainBacklinksCmd)
	rootCmd.AddCommand(pageBacklinksCmd)
}
```

**Step 2: Add API methods**

Add to `internal/api/client.go`:
```go
func (c *Client) GetDomainBacklinks(domain string, limit int) (json.RawMessage, error) {
	data := url.Values{}
	data.Set("domain", domain)
	if limit > 0 {
		data.Set("limit", fmt.Sprintf("%d", limit))
	}

	return c.post("/get_domain_backlinks", data)
}

func (c *Client) GetPageBacklinks(targetURL string, limit int) (json.RawMessage, error) {
	data := url.Values{}
	data.Set("url", targetURL)
	if limit > 0 {
		data.Set("limit", fmt.Sprintf("%d", limit))
	}

	return c.post("/get_page_backlinks", data)
}
```

**Step 3: Build and verify**

Run:
```bash
go build -o ke .
./ke --help
```

Expected: Shows domain-backlinks and page-backlinks in help.

**Step 4: Commit**

```bash
git add cmd/backlinks.go internal/api/client.go
git commit -m "feat: add backlinks commands"
```

---

## Task 9: Install Binary

**Step 1: Install to GOPATH/bin**

Run:
```bash
go install .
```

**Step 2: Verify installation**

Run:
```bash
ke --help
```

Expected: Shows full CLI help with all commands.

**Step 3: Final commit**

```bash
git add -A
git commit -m "chore: finalize CLI implementation"
```

---

## Summary

| Task | Description | Files |
|------|-------------|-------|
| 1 | Project scaffolding | main.go, cmd/root.go, go.mod |
| 2 | API client foundation | internal/api/client.go |
| 3 | Credits command | cmd/credits.go |
| 4 | Countries/currencies | cmd/countries.go, cmd/currencies.go |
| 5 | Keywords command | cmd/keywords.go |
| 6 | Related/PASF commands | cmd/related.go, cmd/pasf.go |
| 7 | Traffic commands | cmd/traffic.go |
| 8 | Backlinks commands | cmd/backlinks.go |
| 9 | Install binary | N/A |
