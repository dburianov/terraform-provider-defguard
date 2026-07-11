package client

import (
	"testing"
)

func TestClientAuth(t *testing.T) {
	baseURL := "http://localhost:8000/api/v1"
	apiToken := "test-token"
	sessionCookie := "test-session"

	client := NewClient(baseURL, apiToken, sessionCookie)
	if client == nil {
		t.Fatal("Failed to create client")
	}

	client = NewClient(baseURL, apiToken, "")
	if client == nil {
		t.Fatal("Failed to create client with API token")
	}

	client = NewClient(baseURL, "", sessionCookie)
	if client == nil {
		t.Fatal("Failed to create client with session cookie")
	}
}
