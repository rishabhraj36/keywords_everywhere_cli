package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNewClient_WithAPIKey(t *testing.T) {
	t.Setenv("KEYWORDS_EVERYWHERE_API_KEY", "test-key")

	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("expected baseURL %s, got %s", defaultBaseURL, client.baseURL)
	}
}

func TestNewClient_MissingAPIKey(t *testing.T) {
	os.Unsetenv("KEYWORDS_EVERYWHERE_API_KEY")

	_, err := NewClient()
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
	if !strings.Contains(err.Error(), "KEYWORDS_EVERYWHERE_API_KEY") {
		t.Errorf("error should mention env var, got: %v", err)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	client := NewClientWithOptions("my-key", "https://custom.api.com", nil)

	if client.apiKey != "my-key" {
		t.Errorf("expected apiKey 'my-key', got %s", client.apiKey)
	}
	if client.baseURL != "https://custom.api.com" {
		t.Errorf("expected custom baseURL, got %s", client.baseURL)
	}
	if client.httpClient == nil {
		t.Error("expected default httpClient, got nil")
	}
}

func TestGetCredits(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/account/credits" {
			t.Errorf("expected /account/credits, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("expected Bearer auth, got %s", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"credits": 5000})
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetCredits()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Credits != 5000 {
		t.Errorf("expected 5000 credits, got %d", resp.Credits)
	}
}

func TestGetCredits_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid api key"}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("bad-key", server.URL, nil)
	_, err := client.GetCredits()
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should contain status code, got: %v", err)
	}
}

func TestGetCountries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/countries" {
			t.Errorf("expected /countries, got %s", r.URL.Path)
		}
		w.Write([]byte(`["us","uk","ca","au"]`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetCountries()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var countries []string
	if err := json.Unmarshal(resp, &countries); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(countries) != 4 {
		t.Errorf("expected 4 countries, got %d", len(countries))
	}
}

func TestGetCurrencies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/currencies" {
			t.Errorf("expected /currencies, got %s", r.URL.Path)
		}
		w.Write([]byte(`["usd","gbp","eur","aud"]`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetCurrencies()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var currencies []string
	if err := json.Unmarshal(resp, &currencies); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(currencies) != 4 {
		t.Errorf("expected 4 currencies, got %d", len(currencies))
	}
}

func TestGetKeywordData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/get_keyword_data" {
			t.Errorf("expected /get_keyword_data, got %s", r.URL.Path)
		}

		r.ParseForm()
		if r.Form.Get("country") != "us" {
			t.Errorf("expected country=us, got %s", r.Form.Get("country"))
		}
		if r.Form.Get("currency") != "usd" {
			t.Errorf("expected currency=usd, got %s", r.Form.Get("currency"))
		}
		if r.Form.Get("dataSource") != "gkp" {
			t.Errorf("expected dataSource=gkp, got %s", r.Form.Get("dataSource"))
		}

		keywords := r.Form["kw[]"]
		if len(keywords) != 2 {
			t.Errorf("expected 2 keywords, got %d", len(keywords))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"credits": 2,
			"data": []map[string]interface{}{
				{"keyword": "seo", "vol": 10000, "cpc": map[string]interface{}{"value": 2.5}},
				{"keyword": "marketing", "vol": 5000, "cpc": map[string]interface{}{"value": 1.8}},
			},
		})
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetKeywordData([]string{"seo", "marketing"}, "us", "usd", "gkp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Credits != 2 {
		t.Errorf("expected 2 credits, got %d", resp.Credits)
	}
	if len(resp.Data) != 2 {
		t.Errorf("expected 2 data items, got %d", len(resp.Data))
	}
}

func TestGetKeywordData_Batching(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		r.ParseForm()
		keywords := r.Form["kw[]"]

		if callCount == 1 && len(keywords) != 100 {
			t.Errorf("first batch should have 100 keywords, got %d", len(keywords))
		}
		if callCount == 2 && len(keywords) != 50 {
			t.Errorf("second batch should have 50 keywords, got %d", len(keywords))
		}

		data := make([]map[string]interface{}, len(keywords))
		for i, kw := range keywords {
			data[i] = map[string]interface{}{"keyword": kw, "vol": 100}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"credits": len(keywords),
			"data":    data,
		})
	}))
	defer server.Close()

	keywords := make([]string, 150)
	for i := range keywords {
		keywords[i] = "keyword" + string(rune('a'+i%26))
	}

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetKeywordData(keywords, "us", "usd", "gkp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 API calls for 150 keywords, got %d", callCount)
	}
	if resp.Credits != 150 {
		t.Errorf("expected 150 total credits, got %d", resp.Credits)
	}
	if len(resp.Data) != 150 {
		t.Errorf("expected 150 data items, got %d", len(resp.Data))
	}
}

func TestGetRelatedKeywords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/get_related_keywords" {
			t.Errorf("expected /get_related_keywords, got %s", r.URL.Path)
		}
		r.ParseForm()
		if r.Form.Get("kw") != "seo" {
			t.Errorf("expected kw=seo, got %s", r.Form.Get("kw"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"credits": 2,
			"data":    []string{"seo tools", "seo marketing", "seo agency"},
		})
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetRelatedKeywords("seo", "us", "usd", "gkp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestGetPASFKeywords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/get_pasf_keywords" {
			t.Errorf("expected /get_pasf_keywords, got %s", r.URL.Path)
		}
		r.ParseForm()
		if r.Form.Get("kw") != "coffee" {
			t.Errorf("expected kw=coffee, got %s", r.Form.Get("kw"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"credits": 2,
			"data":    []string{"coffee near me", "coffee shops", "coffee beans"},
		})
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetPASFKeywords("coffee", "us", "usd", "gkp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestGetDomainKeywords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/get_domain_keywords" {
			t.Errorf("expected /get_domain_keywords, got %s", r.URL.Path)
		}
		r.ParseForm()
		if r.Form.Get("domain") != "example.com" {
			t.Errorf("expected domain=example.com, got %s", r.Form.Get("domain"))
		}
		if r.Form.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", r.Form.Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data": [{"keyword": "example", "vol": 1000}]}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetDomainKeywords("example.com", "us", "usd", 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestGetDomainKeywords_NoLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("limit") != "" {
			t.Errorf("expected no limit param, got %s", r.Form.Get("limit"))
		}
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	_, err := client.GetDomainKeywords("example.com", "us", "usd", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetURLKeywords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/get_url_keywords" {
			t.Errorf("expected /get_url_keywords, got %s", r.URL.Path)
		}
		r.ParseForm()
		if r.Form.Get("url") != "https://example.com/page" {
			t.Errorf("expected url param, got %s", r.Form.Get("url"))
		}

		w.Write([]byte(`{"data": []}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	_, err := client.GetURLKeywords("https://example.com/page", "us", "usd", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDomainTraffic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/get_domain_traffic" {
			t.Errorf("expected /get_domain_traffic, got %s", r.URL.Path)
		}
		r.ParseForm()
		if r.Form.Get("domain") != "example.com" {
			t.Errorf("expected domain=example.com, got %s", r.Form.Get("domain"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"traffic": 50000, "keywords": 1200}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetDomainTraffic("example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestGetURLTraffic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/get_url_traffic" {
			t.Errorf("expected /get_url_traffic, got %s", r.URL.Path)
		}
		r.ParseForm()
		if r.Form.Get("url") != "https://example.com/page" {
			t.Errorf("expected url param, got %s", r.Form.Get("url"))
		}

		w.Write([]byte(`{"traffic": 5000}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	_, err := client.GetURLTraffic("https://example.com/page")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDomainBacklinks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/get_domain_backlinks" {
			t.Errorf("expected /get_domain_backlinks, got %s", r.URL.Path)
		}
		r.ParseForm()
		if r.Form.Get("domain") != "example.com" {
			t.Errorf("expected domain=example.com, got %s", r.Form.Get("domain"))
		}
		if r.Form.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", r.Form.Get("limit"))
		}

		w.Write([]byte(`{"backlinks": [{"url": "https://other.com", "anchor": "example"}]}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	resp, err := client.GetDomainBacklinks("example.com", 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestGetPageBacklinks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/get_page_backlinks" {
			t.Errorf("expected /get_page_backlinks, got %s", r.URL.Path)
		}
		r.ParseForm()
		if r.Form.Get("url") != "https://example.com/page" {
			t.Errorf("expected url param, got %s", r.Form.Get("url"))
		}

		w.Write([]byte(`{"backlinks": []}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-key", server.URL, nil)
	_, err := client.GetPageBacklinks("https://example.com/page", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthorizationHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer my-secret-key" {
			t.Errorf("expected 'Bearer my-secret-key', got '%s'", auth)
		}
		w.Write([]byte(`{"credits": 100}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("my-secret-key", server.URL, nil)
	_, err := client.GetCredits()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
