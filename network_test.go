package sdk

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestGetNetworkDetails_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/network/details").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"friends":[{"id":2,"accessHash":"hash2","nickname":"tr3ble","description":"something2","interests":["mac"],"avatar":{"id":3,"accessHash":"hash3"}}]}`)

	client := NewClient("https://getfriend.ly")
	network, err := client.GetNetworkDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})

	require.NoError(t, err)
	require.Equal(t, &NetworkDetails{
		Friends: []UserDetails{{
			Id:          2,
			AccessHash:  UserAccessHash("hash2"),
			Nickname:    Nickname("tr3ble"),
			Description: UserDescription("something2"),
			Interests: []Interest{
				Interest("mac"),
			},
			Avatar: &FileDescriptor{
				Id:         3,
				AccessHash: FileAccessHash("hash3"),
			},
		}},
	}, network)
}

func TestGetNetworkDetails_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/network/details").
		Reply(400)

	client := NewClient("https://getfriend.ly")
	_, err := client.GetNetworkDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})
	require.Error(t, err)
}

