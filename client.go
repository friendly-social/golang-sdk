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

var (
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrForbidden    = fmt.Errorf("forbidden")
	ErrNotFound     = fmt.Errorf("not found")
)

// NewClient creates basic Client with provided URL (on which backend is located).
func NewClient(endpoint string) *Client {
	return &Client{
		url: endpoint,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewLocalhostClient creates Client with localhost URL and provided port.
func NewLocalhostClient(port int) *Client {
	return NewClient(fmt.Sprintf("http://localhost:%d", port))
}

// NewProductionClient creates Client with Meetacy URL.
func NewProductionClient() *Client {
	return NewClient("https://api.getfriend.ly/")
}

// do creates and executes HTTP request to given path using provided data and fills unmarshalled response to result argument or returns an error if something went wrong.
func (c *Client) do(ctx context.Context, method, path string, auth *Authorization, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal body: %w", err)
		}

		bodyReader = bytes.NewReader(jsonData)
	}

	completePath, err := url.JoinPath(c.url, path)
	if err != nil {
		return fmt.Errorf("invalid path: %s + %s", c.url, path)
	}

	req, err := http.NewRequestWithContext(ctx, method, completePath, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.Header.Set("X-User-Id", fmt.Sprintf("%d", auth.Id))
		req.Header.Set("X-Token", string(auth.Token))
	}

	return c.execute(req, result)
}

// execute send already created request and fills unmarshalled response to result argument or error if something went wrong.
func (c *Client) execute(req *http.Request, result any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if result != nil {
			if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
		}
		return nil
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("%s %s: %w", req.Method, req.URL.Path, ErrUnauthorized)
	case http.StatusForbidden:
		return fmt.Errorf("%s %s: %w", req.Method, req.URL.Path, ErrForbidden)
	case http.StatusNotFound:
		return fmt.Errorf("%s %s: %w", req.Method, req.URL.Path, ErrNotFound)
	default:
		return fmt.Errorf("%s %s: unexpected request with status code %d", req.Method, req.URL.Path, resp.StatusCode)
	}
}
