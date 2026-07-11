package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents an API client for defguard
type Client struct {
	baseURL       string
	apiToken      string
	sessionCookie string
}

// Response represents the API response
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// NewClient creates a new API client
func NewClient(baseURL, apiToken, sessionCookie string) *Client {
	return &Client{
		baseURL:       baseURL,
		apiToken:      apiToken,
		sessionCookie: sessionCookie,
	}
}

// doRequest performs an HTTP request and returns the response
func (c *Client) doRequest(ctx context.Context, method, path string, body []byte) (*Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set authentication header
	if c.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiToken)
	}

	// Set session cookie if provided
	if c.sessionCookie != "" {
		req.AddCookie(&http.Cookie{
			Name:  "defguard_session",
			Value: c.sessionCookie,
		})
	}

	// Add context to request
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}

	return response, nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string) (*Response, error) {
	return c.doRequest(ctx, "GET", path, nil)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	return c.doRequest(ctx, "POST", path, reqBody)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	return c.doRequest(ctx, "PUT", path, reqBody)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string, body interface{}) (*Response, error) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	return c.doRequest(ctx, "DELETE", path, reqBody)
}
