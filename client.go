package sdk

import (
	"bytes"
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
	ErrInvalidPath             = fmt.Errorf("invalid path")
	ErrFailedToMarshalBody     = fmt.Errorf("failed to marshal request body")
	ErrFailedToCreateRequest   = fmt.Errorf("failed to create request")
	ErrFailedToExecuteRequest  = fmt.Errorf("failed to execute request")
	ErrRequestUnauthorized     = fmt.Errorf("unauthorized")
	ErrRequestResourceNotFound = fmt.Errorf("not found")
	ErrFailedToDecodeResponse  = fmt.Errorf("failed to decode response")
	ErrInternalServerError     = fmt.Errorf("internal server error")
	ErrRequestFailed           = fmt.Errorf("request failed")
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

// NewMeetacyClient creates Client with Meetacy URL.
func NewMeetacyClient() *Client {
	return NewClient("https://meetacy.app/friendly")
}

// do creates and executes HTTP request to given path using provided data and fills unmarshalled response to result argument or returns an error if something went wrong.
func (c *Client) do(method, path string, auth *Authorization, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrFailedToMarshalBody, err)
		}

		bodyReader = bytes.NewReader(jsonData)
	}

	completePath, err := url.JoinPath(c.url, path)
	if err != nil {
		return fmt.Errorf("%w: %s + %s", ErrInvalidPath, c.url, path)
	}

	req, err := http.NewRequest(method, completePath, bodyReader)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateRequest, err)
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
		return fmt.Errorf("%w: %w", ErrFailedToExecuteRequest, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if result != nil {
			if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
				return fmt.Errorf("%w: %w", ErrFailedToDecodeResponse, err)
			}
		}
		return nil
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("%s %s: %w", req.Method, req.URL.Path, ErrRequestUnauthorized)
	case http.StatusNotFound:
		return fmt.Errorf("%s %s: %w", req.Method, req.URL.Path, ErrRequestResourceNotFound)
	case http.StatusInternalServerError:
		return fmt.Errorf("%s %s: %w", req.Method, req.URL.Path, ErrInternalServerError)
	default:
		return fmt.Errorf("%s %s with status %d: %w", req.Method, req.URL.Path, resp.StatusCode, ErrRequestFailed)
	}
}
