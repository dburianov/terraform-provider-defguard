package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

// Client represents the Defguard API client
type Client struct {
	baseURL            string
	apiToken           string
	session            string // Session cookie name (default: "defguard_session")
	sessionCookieValue string // Session cookie value (for direct header setting)
	httpClient         *http.Client
}

// NewClient creates a new Defguard API client
func NewClient(endpoint string, apiToken string) *Client {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	return &Client{
		baseURL:  endpoint,
		apiToken: apiToken,
		session:  "defguard_session",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     cookieJar,
		},
	}
}

// SetSessionCookie sets the session cookie name
func (c *Client) SetSessionCookie(name string) {
	c.session = name
}

// SetSessionValue sets an existing session cookie value directly
func (c *Client) SetSessionValue(value string) error {
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse baseURL: %w", err)
	}

	c.sessionCookieValue = value

	cookie := &http.Cookie{
		Name:   c.session,
		Value:  value,
		Path:   "/",
		Domain: parsedURL.Hostname(),
	}

	c.httpClient.Jar.SetCookies(parsedURL, []*http.Cookie{cookie})
	return nil
}

// LoginResult represents the response from the auth endpoint
type LoginResult struct {
	Token   string `json:"token"`
	Message string `json:"msg"`
}

// LoginWithCredentials authenticates with the Defguard API using username and password
func (c *Client) LoginWithCredentials(ctx context.Context, username, password string) (*LoginResult, error) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	respObj, err := c.Post(ctx, "/api/v1/auth", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	var result LoginResult
	if err := respObj.Unmarshal(&result); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %w", err)
	}

	return &result, nil
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
	requestURL := fmt.Sprintf("%s%s", c.baseURL, req.Path)

	fmt.Fprintf(os.Stderr, "DEBUG doRequest: baseURL=%s, path=%s\n", c.baseURL, req.Path)
	fmt.Fprintf(os.Stderr, "DEBUG requestURL=%s\n", requestURL)
	fmt.Fprintf(os.Stderr, "DEBUG doRequest: %s %s\n", req.Method, requestURL)

	var bodyReader io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
		fmt.Fprintf(os.Stderr, "DEBUG Request body: %s\n", string(jsonBody))
	} else {
		fmt.Fprintf(os.Stderr, "DEBUG Request body is nil\n")
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Add session cookie directly to request headers if set (takes precedence over cookiejar)
	if c.apiToken != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
		fmt.Fprintf(os.Stderr, "DEBUG Using Authorization header: Bearer %s...\n", c.apiToken[:20])
	} else if c.sessionCookieValue != "" {
		// Use the cookiejar to add cookies instead of manually setting header
		// This ensures proper cookie handling by Go's HTTP client
		parsedURL, err := url.Parse(c.baseURL)
		if err == nil {
			cookies := []*http.Cookie{{Name: c.session, Value: c.sessionCookieValue, Path: "/", Domain: parsedURL.Hostname()}}
			c.httpClient.Jar.SetCookies(parsedURL, cookies)
			fmt.Fprintf(os.Stderr, "DEBUG Added cookie to jar: %s=...\n", c.session)
		}

		// Get cookies that will be sent
		jarCookies := c.httpClient.Jar.Cookies(httpReq.URL)
		if len(jarCookies) > 0 {
			var cookieStrings []string
			for _, cookie := range jarCookies {
				cookieStrings = append(cookieStrings, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
			}
			fmt.Fprintf(os.Stderr, "DEBUG Cookies to send: %s\n", strings.Join(cookieStrings, "; "))
		}

		// Set Cookie header directly for debugging (so we can see it in headers output)
		if c.sessionCookieValue != "" {
			cookieHeader := fmt.Sprintf("%s=%s", c.session, c.sessionCookieValue)
			httpReq.Header.Set("Cookie", cookieHeader)
			fmt.Fprintf(os.Stderr, "DEBUG Setting Cookie header: %s\n", cookieHeader)
		}
	}

	// Log all headers being sent
	fmt.Fprintf(os.Stderr, "DEBUG Headers being sent:\n")
	for k, v := range httpReq.Header {
		if len(v) > 0 {
			fmt.Fprintf(os.Stderr, "  %s: %v\n", k, v[0])
		}
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
