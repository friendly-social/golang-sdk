package sdk

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestClients(t *testing.T) {
	t.Run("Default Client", func(t *testing.T) {
		client := NewClient("https://example.com")
		require.Equal(t, client.url, "https://example.com")
		require.Equal(t, client.http, &http.Client{Timeout: 30 * time.Second})
	})

	t.Run("Localhost Client", func(t *testing.T) {
		client := NewLocalhostClient(8080)
		require.Equal(t, client.url, "http://localhost:8080")
		require.Equal(t, client.http, &http.Client{Timeout: 30 * time.Second})
	})

	t.Run("Production Client", func(t *testing.T) {
		client := NewProductionClient()
		require.Equal(t, client.url, "https://api.getfriend.ly/")
		require.Equal(t, client.http, &http.Client{Timeout: 30 * time.Second})
	})
}

func TestDoAndExecute(t *testing.T) {
	type input struct {
		host   string
		path   string
		method string
		body   any
		auth   *Authorization
	}

	type response struct {
		Trollge string
	}

	type testCase struct {
		name             string
		input            input
		mockStatus       int
		mockResponse     string
		mockError        error
		expectedHeaders  map[string]string
		expectedBody     string
		expectedResponse response
		expectError      bool
		expectedError    error
	}

	cases := []testCase{
		{
			name: "Success",
			input: input{
				host:   "https://getfriend.ly",
				method: "GET",
				path:   "/ping",
				body:   struct{ Field string }{"something interesting"},
				auth: &Authorization{
					Id:         1,
					Token:      Token("token"),
					AccessHash: UserAccessHash("hash"),
				},
			},
			mockStatus:   200,
			mockResponse: `{"Trollge":"trollge"}`,
			expectedBody: `{"Field":"something interesting"}`,
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-User-Id":    "1",
				"X-Token":      "token",
			},
			expectedResponse: response{
				Trollge: "trollge",
			},
		},
		{
			name: "Failed Marshal Body",
			input: input{
				body: func() {},
			},
			expectError: true,
		},
		{
			name: "Invalid Path",
			input: input{
				host: "::invalid",
			},
			expectError: true,
		},
		{
			name: "Invalid Request",
			input: input{
				method: "BAD METHOD",
			},
			expectError: true,
		},
		{
			name: "Network Error",
			input: input{
				host:   "https://getfriend.ly",
				method: "GET",
				path:   "/ping",
			},
			mockError:   fmt.Errorf("some shit"),
			expectError: true,
		},
		{
			name: "Unauthorized",
			input: input{
				host:   "https://getfriend.ly",
				method: "GET",
				path:   "/ping",
			},
			mockStatus:    401,
			expectedError: ErrUnauthorized,
		},
		{
			name: "Forbidden",
			input: input{
				host:   "https://getfriend.ly",
				method: "GET",
				path:   "/ping",
			},
			mockStatus:    403,
			expectedError: ErrForbidden,
		},
		{
			name: "Not Found",
			input: input{
				host:   "https://getfriend.ly",
				method: "GET",
				path:   "/ping",
			},
			mockStatus:    404,
			expectedError: ErrNotFound,
		},
		{
			name: "Something went wrong",
			input: input{
				host:   "https://getfriend.ly",
				method: "GET",
				path:   "/ping",
			},
			mockStatus:  418,
			expectError: true,
		},
		{
			name: "Invalid response",
			input: input{
				host:   "https://getfriend.ly",
				method: "GET",
				path:   "/ping",
			},
			mockStatus:   200,
			mockResponse: "invalid",
			expectError:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			r := gock.New(tc.input.host)
			switch tc.input.method {
			case "GET":
				r = r.Get(tc.input.path)
			case "POST":
				r = r.Post(tc.input.path)
			}

			for k, v := range tc.expectedHeaders {
				r = r.MatchHeader(k, v)
			}

			r = r.JSON(tc.expectedBody)
			if tc.mockError != nil {
				r.ReplyError(tc.mockError)
			} else {
				r.Reply(tc.mockStatus).
					JSON(tc.mockResponse)
			}

			var resp response
			client := NewClient(tc.input.host)
			err := client.do(context.Background(), tc.input.method, tc.input.path, tc.input.auth, tc.input.body, &resp)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResponse, resp)
			}
		})
	}
}

func TestDoWithCancellation(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/ping").
		Reply(200).
		Delay(100 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	client := NewClient("https://getfriend.ly")
	err := client.do(ctx, "GET", "/ping", nil, nil, nil)

	require.ErrorIs(t, err, context.DeadlineExceeded)
}
