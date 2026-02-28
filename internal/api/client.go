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
	"time"
)

const defaultBaseURL = "https://api.keywordseverywhere.com/v1"

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	apiKey := os.Getenv("KEYWORDS_EVERYWHERE_API_KEY")
	if apiKey == "" {
		return nil, errors.New("KEYWORDS_EVERYWHERE_API_KEY environment variable not set")
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// NewClientWithOptions creates a client with custom options (for testing)
func NewClientWithOptions(apiKey, baseURL string, httpClient *http.Client) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *Client) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
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
	req, err := http.NewRequest("POST", c.baseURL+endpoint, strings.NewReader(data.Encode()))
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

// CreditsResponse represents the account credits response
type CreditsResponse struct {
	Credits int `json:"credits"`
}

// GetCredits retrieves the current account credit balance
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

// GetCountries retrieves the list of supported countries
func (c *Client) GetCountries() ([]byte, error) {
	return c.get("/countries")
}

// GetCurrencies retrieves the list of supported currencies
func (c *Client) GetCurrencies() ([]byte, error) {
	return c.get("/currencies")
}

// KeywordDataResponse represents the response from keyword data API
type KeywordDataResponse struct {
	Credits int               `json:"credits"`
	Data    []json.RawMessage `json:"data"`
}

// GetKeywordData retrieves volume, CPC, and competition data for keywords
// Keywords are batched in groups of 100 to respect API limits
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

// GetRelatedKeywords retrieves related keywords for a given keyword
func (c *Client) GetRelatedKeywords(keyword, country, currency, source string) ([]byte, error) {
	data := url.Values{}
	data.Set("country", country)
	data.Set("currency", currency)
	data.Set("dataSource", source)
	data.Set("kw", keyword)

	return c.post("/get_related_keywords", data)
}

// GetPASFKeywords retrieves "People Also Search For" keywords for a given keyword
func (c *Client) GetPASFKeywords(keyword, country, currency, source string) ([]byte, error) {
	data := url.Values{}
	data.Set("country", country)
	data.Set("currency", currency)
	data.Set("dataSource", source)
	data.Set("kw", keyword)

	return c.post("/get_pasf_keywords", data)
}

// GetDomainKeywords retrieves keywords that a domain ranks for
func (c *Client) GetDomainKeywords(domain, country, currency string, limit int) ([]byte, error) {
	data := url.Values{}
	data.Set("domain", domain)
	data.Set("country", country)
	data.Set("currency", currency)
	if limit > 0 {
		data.Set("limit", fmt.Sprintf("%d", limit))
	}

	return c.post("/get_domain_keywords", data)
}

// GetURLKeywords retrieves keywords that a URL ranks for
func (c *Client) GetURLKeywords(targetURL, country, currency string, limit int) ([]byte, error) {
	data := url.Values{}
	data.Set("url", targetURL)
	data.Set("country", country)
	data.Set("currency", currency)
	if limit > 0 {
		data.Set("limit", fmt.Sprintf("%d", limit))
	}

	return c.post("/get_url_keywords", data)
}

// GetDomainTraffic retrieves traffic metrics for a domain
func (c *Client) GetDomainTraffic(domain string) ([]byte, error) {
	data := url.Values{}
	data.Set("domain", domain)

	return c.post("/get_domain_traffic", data)
}

// GetURLTraffic retrieves traffic metrics for a URL
func (c *Client) GetURLTraffic(targetURL string) ([]byte, error) {
	data := url.Values{}
	data.Set("url", targetURL)

	return c.post("/get_url_traffic", data)
}

// GetDomainBacklinks retrieves backlinks for a domain
func (c *Client) GetDomainBacklinks(domain string, limit int) ([]byte, error) {
	data := url.Values{}
	data.Set("domain", domain)
	if limit > 0 {
		data.Set("limit", fmt.Sprintf("%d", limit))
	}

	return c.post("/get_domain_backlinks", data)
}

// GetPageBacklinks retrieves backlinks for a specific page URL
func (c *Client) GetPageBacklinks(targetURL string, limit int) ([]byte, error) {
	data := url.Values{}
	data.Set("url", targetURL)
	if limit > 0 {
		data.Set("limit", fmt.Sprintf("%d", limit))
	}

	return c.post("/get_page_backlinks", data)
}
