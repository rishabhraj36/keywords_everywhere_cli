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
		httpClient: &http.Client{Timeout: 30 * time.Second},
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
