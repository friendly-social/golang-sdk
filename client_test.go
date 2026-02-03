package sdk

import (
	"context"
	"fmt"
	"io"
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
		require.Equal(t, client.http.Timeout, 30*time.Second)
	})

	t.Run("Localhost Client", func(t *testing.T) {
		client := NewLocalhostClient(8080)
		require.Equal(t, client.url, "http://localhost:8080")
		require.Equal(t, client.http.Timeout, 30*time.Second)
	})

	t.Run("Production Client", func(t *testing.T) {
		client := NewProductionClient()
		require.Equal(t, client.url, "https://api.getfriend.ly/")
		require.Equal(t, client.http.Timeout, 30*time.Second)
	})
}

func TestDo_SuccessJSON(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/ping").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"Field": "something interesting"}`).
		Reply(200).
		JSON(`{"Trollge": "trollge"}`)

	client := NewClient("https://getfriend.ly")
	auth := &Authorization{
		Id:         1,
		Token:      Token("token"),
		AccessHash: UserAccessHash("hash"),
	}

	var resp struct{ Trollge string }
	err := client.do(context.Background(), auth, "GET", "/ping", struct{ Field string }{"something interesting"}, &resp)

	require.NoError(t, err)
	require.Equal(t, struct{ Trollge string }{"trollge"}, resp)
}

func TestDo_SuccessPlaintext(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/ping").
		Reply(200).
		SetHeader("Content-Type", "text/plain").
		BodyString("trollge")

	var resp string
	client := NewClient("https://getfriend.ly")
	err := client.do(context.Background(), nil, "GET", "/ping", nil, &resp)

	require.NoError(t, err)
	require.Equal(t, "trollge", resp)
}

func TestDo_Cancel(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/ping").
		Reply(200).
		Delay(100 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	client := NewClient("https://getfriend.ly")
	err := client.do(ctx, nil, "GET", "/ping", nil, nil)

	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestDo_FailedToMarshalBody(t *testing.T) {
	client := NewClient("https://getfriend.ly")
	err := client.do(context.Background(), nil, "GET", "/ping", func() {}, nil)
	require.Error(t, err)
}

func TestDo_InvalidURL(t *testing.T) {
	client := NewClient("::invalid")
	err := client.do(context.Background(), nil, "GET", "/ping", nil, nil)
	require.Error(t, err)
}

func TestDo_InvalidRequest(t *testing.T) {
	client := NewClient("https://getfriend.ly")
	err := client.do(context.Background(), nil, "SOMETHING BAD!!!", "/ping", nil, nil)
	require.Error(t, err)
}

func TestDo_StatusCodes(t *testing.T) {
	cases := []struct {
		code    int
		wantErr error
	}{
		{401, ErrUnauthorized},
		{403, ErrForbidden},
		{404, ErrNotFound},
		{418, nil},
	}

	for _, tc := range cases {
		defer gock.Off()
		gock.New("https://getfriend.ly").
			Get("/ping").
			Reply(tc.code)

		client := NewClient("https://getfriend.ly")
		err := client.do(context.Background(), nil, "GET", "/ping", nil, nil)

		require.Error(t, err)
		if tc.wantErr != nil {
			require.ErrorIs(t, err, tc.wantErr)
		}
	}
}

func TestDo_InvalidResponseJSON(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/ping").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString("bad")

	var resp any
	client := NewClient("https://getfriend.ly")
	err := client.do(context.Background(), nil, "GET", "/ping", nil, &resp)
	require.Error(t, err)
}

func TestDo_InvalidResponsePlaintext(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/ping").
		Reply(200).
		SetHeader("Content-Type", "text/plain").
		BodyString("bad")

	client := NewClient("https://getfriend.ly")
	err := client.do(context.Background(), nil, "GET", "/ping", nil, "bad")
	require.Error(t, err)
}

func TestDo_InvalidContentType(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/ping").
		Reply(200).
		SetHeader("Content-Type", "bad")

	client := NewClient("https://getfriend.ly")
	err := client.do(context.Background(), nil, "GET", "/ping", nil, "bad")
	require.Error(t, err)
}

type badReader struct{}

func (e *badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("boom")
}

type badRoundTripper struct{}

func (badRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(&badReader{}),
	}, nil
}

func TestDo_FailedToReadBody(t *testing.T) {
	client := NewClient("https://getfriend.ly")
	client.http.Transport = badRoundTripper{}

	err := client.do(context.Background(), nil, "GET", "/ping", nil, nil)
	require.Error(t, err)
}
