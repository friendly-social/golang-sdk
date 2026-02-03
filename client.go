package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
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
func NewClient(url string) *Client {
	return &Client{
		url: url,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewLocalhostClient creates Client with localhost URL and provided port.
func NewLocalhostClient(port int) *Client {
	return NewClient(fmt.Sprintf("http://localhost:%d", port))
}

// NewProductionClient creates Client with actual backend.
func NewProductionClient() *Client {
	return NewClient("https://api.getfriend.ly/")
}

func (c *Client) do(ctx context.Context, auth *Authorization, method, path string, body any, result any) error {
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

func (c *Client) execute(req *http.Request, result any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	ct := resp.Header.Get("Content-Type")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if result == nil {
			return nil
		}

		mediaType, _, _ := mime.ParseMediaType(ct)
		switch mediaType {
		case "application/json":
			if err = json.Unmarshal(body, result); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
		case "text/plain":
			strPtr, ok := result.(*string)
			if !ok {
				return fmt.Errorf("expected *string result for text/plain response")
			}
			*strPtr = string(body)
		default:
			return fmt.Errorf("unexpected content type: %s", ct)
		}

		return nil
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("%s %s: %w\n%s", req.Method, req.URL.Path, ErrUnauthorized, body)
	case http.StatusForbidden:
		return fmt.Errorf("%s %s: %w\n%s", req.Method, req.URL.Path, ErrForbidden, body)
	case http.StatusNotFound:
		return fmt.Errorf("%s %s: %w\n%s", req.Method, req.URL.Path, ErrNotFound, body)
	default:
		return fmt.Errorf("%s %s: unexpected request with status code %d\n%s", req.Method, req.URL.Path, resp.StatusCode, body)
	}
}
