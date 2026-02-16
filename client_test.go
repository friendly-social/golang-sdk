package sdk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestClients(t *testing.T) {
	t.Run("Default Client", func(t *testing.T) {
		client := NewClient()
		require.Equal(t, client.url, "https://api.getfriend.ly")
		require.Equal(t, client.http.Timeout, 30*time.Second)
	})

	t.Run("Customized Client", func(t *testing.T) {
		client := NewClient().
			WithHTTPClient(&http.Client{}).
			WithTimeout(5 * time.Second).
			WithBaseURL("https://example.com")

		require.Equal(t, client.url, "https://example.com")
		require.Equal(t, client.http, &http.Client{Timeout: 5 * time.Second})
	})
}

func TestDo_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/ping").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"Field": "something interesting"}`).
		Reply(200).
		JSON(`{"Trollge": "trollge"}`)

	client := NewClient()
	auth := &Authorization{
		Id:         MockUserId(1),
		Token:      MockToken("token"),
		AccessHash: MockUserAccessHash("hash"),
	}

	var resp struct{ Trollge string }
	err := client.do(context.Background(), auth, "GET", "/ping", struct{ Field string }{"something interesting"}, &resp)

	require.NoError(t, err)
	require.Equal(t, struct{ Trollge string }{"trollge"}, resp)
}

func TestDo_Cancel(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/ping").
		Reply(200).
		Delay(100 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	client := NewClient()
	err := client.do(ctx, nil, "GET", "/ping", nil, nil)

	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestDo_FailedToMarshalBody(t *testing.T) {
	client := NewClient()
	err := client.do(context.Background(), nil, "GET", "/ping", func() {}, nil)
	require.Error(t, err)
}

func TestDo_InvalidURL(t *testing.T) {
	client := NewClient().
		WithBaseURL("::invalid")

	err := client.do(context.Background(), nil, "GET", "/ping", nil, nil)
	require.Error(t, err)
}

func TestDo_InvalidRequest(t *testing.T) {
	client := NewClient()
	err := client.do(context.Background(), nil, "SOMETHING BAD!!!", "/ping", nil, nil)
	require.Error(t, err)
}

func TestDo_APIError(t *testing.T) {
	cases := []struct {
		name string
		code int
		body []byte
	}{
		{"Unauthorized", 401, []byte("invalid auth")},
		{"Forbidden", 403, []byte("you are not admin")},
		{"Not Found", 404, []byte("not found")},
		{"I'm a teapot", 418, []byte("lol")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()
			gock.New("https://api.getfriend.ly").
				Get("/ping").
				Reply(tc.code).
				Body(io.NopCloser(strings.NewReader(string(tc.body))))

			client := NewClient()
			err := client.do(context.Background(), nil, "GET", "/ping", nil, nil)

			var apiError APIError
			require.ErrorAs(t, err, &apiError)
			require.Equal(t, tc.code, apiError.Code)
			require.Equal(t, tc.body, apiError.Body)
		})
	}
}

func TestDo_InvalidResponse(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/ping").
		Reply(200).
		BodyString("bad")

	var resp any
	client := NewClient()
	err := client.do(context.Background(), nil, "GET", "/ping", nil, &resp)
	require.Error(t, err)
}

type errorReader struct{}

func (r *errorReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("something went wrong...")
}

func TestDo_FailedToReadError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/ping").
		Reply(418).
		Map(func(resp *http.Response) *http.Response {
			resp.Body = io.NopCloser(&errorReader{})
			return resp
		})

	client := NewClient()
	err := client.do(context.Background(), nil, "GET", "/ping", nil, nil)
	require.Error(t, err)
}
