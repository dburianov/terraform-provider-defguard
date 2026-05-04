package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the Defguard API client
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

// NewClient creates a new Defguard API client
func NewClient(endpoint string, apiToken string) *Client {
	return &Client{
		baseURL:    endpoint,
		apiToken:   apiToken,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Request is a generic request structure
type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

// Response is a generic response structure
type Response struct {
	StatusCode int
	Body       []byte
	Err        error
}

// doRequest performs an HTTP request to the Defguard API
func (c *Client) doRequest(ctx context.Context, req Request) (*Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, req.Path)

	var bodyReader io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if c.apiToken != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &Response{
			StatusCode: resp.StatusCode,
			Body:       body,
			Err:        fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body)),
		}, nil
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
		Err:        nil,
	}, nil
}

// UnmarshalResponse unmarshals a JSON response into a struct
func UnmarshalResponse[T any](resp *Response) (*T, error) {
	if resp.Err != nil {
		return nil, resp.Err
	}

	var result T
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Unmarshal unmarshals a JSON response into a struct
func (r *Response) Unmarshal(v interface{}) error {
	if r.Err != nil {
		return r.Err
	}
	return json.Unmarshal(r.Body, v)
}

// UnmarshalListResponse unmarshals a JSON array response into a slice
func UnmarshalListResponse[T any](resp *Response) ([]T, error) {
	if resp.Err != nil {
		return nil, resp.Err
	}

	var result []T
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string) (*Response, error) {
	req := Request{
		Method: "GET",
		Path:   path,
	}
	return c.doRequest(ctx, req)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	req := Request{
		Method: "POST",
		Path:   path,
		Body:   body,
	}
	return c.doRequest(ctx, req)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	req := Request{
		Method: "PUT",
		Path:   path,
		Body:   body,
	}
	return c.doRequest(ctx, req)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string, body interface{}) (*Response, error) {
	req := Request{
		Method: "DELETE",
		Path:   path,
		Body:   body,
	}
	return c.doRequest(ctx, req)
}
