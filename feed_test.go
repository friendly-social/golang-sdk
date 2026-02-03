package sdk

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestGetFeedQueue_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/feed/queue").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"entries":[{"isExtendedNetwork":true,"commonFriends":[],"details":{"id":2},"isRequest":true}]}`)

	client := NewClient("https://getfriend.ly")
	feed, err := client.GetFeedQueue(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})

	require.NoError(t, err)
	require.Equal(t, &FeedQueue{
		Entries: []FeedEntry{{
			IsRequest:         true,
			IsExtendedNetwork: true,
			CommonFriends:     []UserDetails{},
			Details:           UserDetails{Id: 2},
		}},
	}, feed)
}

func TestGetFeedQueue_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/feed/queue").
		Reply(400)

	client := NewClient("https://getfriend.ly")
	_, err := client.GetFeedQueue(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})
	require.Error(t, err)
}
