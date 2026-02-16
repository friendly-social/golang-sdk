package sdk

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestGetNetworkDetails_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/network/details").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"friends":[{"id":2,"accessHash":"hash2","nickname":"tr3ble","description":"something2","interests":["mac"],"avatar":{"id":3,"accessHash":"hash3"}}]}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	network, err := client.GetNetworkDetails(context.Background(), auth)

	require.NoError(t, err)
	require.Equal(t, &NetworkDetails{
		Friends: []UserDetails{{
			Id:          MockUserId(2),
			AccessHash:  MockUserAccessHash("hash2"),
			Nickname:    MockNickname("tr3ble"),
			Description: MockUserDescription("something2"),
			Interests: MockInterests([]Interest{
				MockInterest("mac"),
			}),
			Avatar: &FileDescriptor{
				Id:         MockFileId(3),
				AccessHash: MockFileAccessHash("hash3"),
			},
		}},
	}, network)
}

func TestGetNetworkDetails_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/network/details").
		Reply(400)

	client := NewClient()
	_, err := client.GetNetworkDetails(context.Background(), nil)
	require.Error(t, err)
}
