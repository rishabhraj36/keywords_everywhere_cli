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
