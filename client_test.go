package sdk

import (
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestClients(t *testing.T) {
	expectedTimeout := 30 * time.Second

	table := []struct {
		name string
		want string
		got  *Client
	}{
		{"manual", "https://example.com", NewClient("https://example.com")},
		{"localhost", "http://localhost:8080", NewLocalhostClient(8080)},
		{"meetacy", "https://meetacy.app/friendly", NewMeetacyClient()},
	}

	for _, row := range table {
		t.Run(row.name, func(t *testing.T) {
			require.Equal(t, row.want, row.got.url)
			require.Equal(t, expectedTimeout, row.got.http.Timeout)
		})
	}
}

func TestDo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://getfriend.ly").
			Get("/test").
			MatchHeader("X-User-Id", "1").
			MatchHeader("X-Token", "1").
			Reply(200)

		client := NewClient("https://getfriend.ly")
		resp, err := client.do("GET", "/test", &Authorization{Id: 1, AccessHash: "1", Token: "1"}, nil)

		require.NotNil(t, resp)
		require.NoError(t, err)
	})

	t.Run("Invalid Body", func(t *testing.T) {
		client := NewClient("https://getfriend.ly")
		resp, err := client.do("GET", "/test", nil, struct{ Test func() }{func() {}})
		require.Nil(t, resp)
		require.Error(t, err)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		client := NewClient("://")
		resp, err := client.do("GET", "/test", nil, nil)
		require.Nil(t, resp)
		require.Error(t, err)
	})
}
