package sdk

import (
	"context"
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

	client := NewClient("https://api.getfriend.ly")
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

	client := NewClient("https://api.getfriend.ly")
	err := client.do(ctx, nil, "GET", "/ping", nil, nil)

	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestDo_FailedToMarshalBody(t *testing.T) {
	client := NewClient("https://api.getfriend.ly")
	err := client.do(context.Background(), nil, "GET", "/ping", func() {}, nil)
	require.Error(t, err)
}

func TestDo_InvalidURL(t *testing.T) {
	client := NewClient("::invalid")
	err := client.do(context.Background(), nil, "GET", "/ping", nil, nil)
	require.Error(t, err)
}

func TestDo_InvalidRequest(t *testing.T) {
	client := NewClient("https://api.getfriend.ly")
	err := client.do(context.Background(), nil, "SOMETHING BAD!!!", "/ping", nil, nil)
	require.Error(t, err)
}

func TestDo_StatusCodes(t *testing.T) {
	cases := []struct {
		name    string
		code    int
		wantErr error
	}{
		{"Unauthorized", 401, ErrUnauthorized},
		{"Forbidden", 403, ErrForbidden},
		{"Not Found", 404, ErrNotFound},
		{"I'm a teapot", 418, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()
			gock.New("https://api.getfriend.ly").
				Get("/ping").
				Reply(tc.code)

			client := NewClient("https://api.getfriend.ly")
			err := client.do(context.Background(), nil, "GET", "/ping", nil, nil)

			require.Error(t, err)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}

func TestDo_InvalidResponse(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/ping").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString("bad")

	var resp any
	client := NewClient("https://api.getfriend.ly")
	err := client.do(context.Background(), nil, "GET", "/ping", nil, &resp)
	require.Error(t, err)
}
