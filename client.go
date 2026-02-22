package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is an entity for interacting with Friendly API.
type Client struct {
	url  string
	http *http.Client
}

// APIError represents some error returned by Friendly API.
type APIError struct {
	Code int
	Body []byte
}

func (e APIError) Error() string {
	return fmt.Sprintf("HTTP error with code %d, returned body:\n%s", e.Code, string(e.Body))
}

// NewClient creates basic Client with provided URL (on which backend is located).
func NewClient() *Client {
	return &Client{
		url: "https://api.getfriend.ly",
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WithHTTPClient sets custom http.Client for Friendly Client.
func (c *Client) WithHTTPClient(http *http.Client) *Client {
	c.http = http
	return c
}

// WithTimeout sets custom request timeout for Client.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.http.Timeout = timeout
	return c
}

// WithBaseURL sets custom URL for Client.
func (c *Client) WithBaseURL(url string) *Client {
	c.url = url
	return c
}

func (c *Client) do(ctx context.Context, auth *Authorization, method, path string, body any, result any) error {
	req, err := c.newRequest(ctx, auth, method, path, body)
	if err != nil {
		return err
	}

	return c.execute(req, result)
}

func (c *Client) newRequest(ctx context.Context, auth *Authorization, method, path string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}

		bodyReader = bytes.NewReader(jsonData)
	}

	completePath, err := url.JoinPath(c.url, path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s + %s", c.url, path)
	}

	req, err := http.NewRequestWithContext(ctx, method, completePath, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.Header.Set("X-User-Id", fmt.Sprintf("%d", auth.Id))
		req.Header.Set("X-Token", string(auth.Token.value))
	}

	return req, nil
}

func (c *Client) execute(req *http.Request, result any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if result == nil {
			return nil
		}

		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP error body")
	}

	return fmt.Errorf("unexpected response: %w", APIError{Code: resp.StatusCode, Body: body})
}
