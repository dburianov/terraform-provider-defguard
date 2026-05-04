package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:8000", "test-token")

	if client.baseURL != "http://localhost:8000" {
		t.Errorf("expected baseURL to be 'http://localhost:8000', got '%s'", client.baseURL)
	}
	if client.apiToken != "test-token" {
		t.Errorf("expected apiToken to be 'test-token', got '%s'", client.apiToken)
	}
	if client.httpClient == nil {
		t.Error("expected httpClient to be not nil")
	}
}

func TestGet(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   1,
			"name": "test",
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "")

	resp, err := client.Get(context.Background(), "/test")
	if err != nil {
		t.Errorf("expected no error, got '%v'", err)
	}

	var result map[string]interface{}
	if err := resp.Unmarshal(&result); err != nil {
		t.Errorf("expected no error unmarshaling, got '%v'", err)
	}

	if result["id"] != 1.0 {
		t.Errorf("expected id to be 1, got '%v'", result["id"])
	}
	if result["name"] != "test" {
		t.Errorf("expected name to be 'test', got '%v'", result["name"])
	}
}

func TestPost(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   1,
			"name": "created",
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "")

	resp, err := client.Post(context.Background(), "/test", map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Errorf("expected no error, got '%v'", err)
	}

	var result map[string]interface{}
	if err := resp.Unmarshal(&result); err != nil {
		t.Errorf("expected no error unmarshaling, got '%v'", err)
	}

	if result["id"] != 1.0 {
		t.Errorf("expected id to be 1, got '%v'", result["id"])
	}
	if result["name"] != "created" {
		t.Errorf("expected name to be 'created', got '%v'", result["name"])
	}
}

func TestPut(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   1,
			"name": "updated",
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "")

	resp, err := client.Put(context.Background(), "/test/1", map[string]interface{}{
		"name": "updated",
	})
	if err != nil {
		t.Errorf("expected no error, got '%v'", err)
	}

	var result map[string]interface{}
	if err := resp.Unmarshal(&result); err != nil {
		t.Errorf("expected no error unmarshaling, got '%v'", err)
	}

	if result["id"] != 1.0 {
		t.Errorf("expected id to be 1, got '%v'", result["id"])
	}
	if result["name"] != "updated" {
		t.Errorf("expected name to be 'updated', got '%v'", result["name"])
	}
}

func TestDelete(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "")

	resp, err := client.Delete(context.Background(), "/test/1", nil)
	if err != nil {
		t.Errorf("expected no error, got '%v'", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got '%d'", resp.StatusCode)
	}
}

func TestDoRequestWithError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "bad request",
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "")

	resp, err := client.Get(context.Background(), "/test")
	if err != nil {
		t.Errorf("expected no error from doRequest, got '%v'", err)
	}

	// The response should have an error set
	if resp.Err == nil {
		t.Error("expected response error to be set for non-2xx status")
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code 400, got '%d'", resp.StatusCode)
	}
}

func TestDoRequestWithInvalidURL(t *testing.T) {
	client := NewClient("http://invalid-host-that-does-not-exist.example.com", "")

	_, err := client.Get(context.Background(), "/test")
	if err == nil {
		t.Error("expected error for invalid host")
	}
}
