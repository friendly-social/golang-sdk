package sdk

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestGenerateFriendToken_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/friends/generate").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"token":"token2"}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	token, err := client.GenerateFriendToken(context.Background(), auth)

	require.NoError(t, err)
	require.Equal(t, MockFriendToken("token2"), token)
}

func TestGenerateFriendToken_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/friends/generate").
		Reply(400)

	client := NewClient()
	_, err := client.GenerateFriendToken(context.Background(), nil)
	require.Error(t, err)
}

func TestAddFriend_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/friends/add").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"userId":2, "token":"token2"}`).
		Reply(200).
		JSON(`{"type":"nice"}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := client.AddFriend(context.Background(), auth, MockFriendToken("token2"), MockUserId(2))
	require.NoError(t, err)
}

func TestAddFriend_TokenExpired(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/friends/add").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"userId":2, "token":"token2"}`).
		Reply(200).
		JSON(`{"type":"FriendTokenExpired"}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := client.AddFriend(context.Background(), auth, MockFriendToken("token2"), MockUserId(2))

	require.Error(t, err)
	require.ErrorIs(t, err, ErrFriendTokenExpired)
}

func TestAddFriend_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/friends/add").
		Reply(400)

	client := NewClient()
	err := client.AddFriend(context.Background(), nil, MockFriendToken("token2"), MockUserId(2))
	require.Error(t, err)
}

func TestSendFriendRequest_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/friends/request").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"userId":2, "userAccessHash":"hash2"}`).
		Reply(200)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := client.SendFriendRequest(context.Background(), auth, MockUserId(2), MockUserAccessHash("hash2"))
	require.NoError(t, err)
}

func TestSendFriendRequest_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/friends/request").
		Reply(400)

	client := NewClient()
	err := client.SendFriendRequest(context.Background(), nil, MockUserId(2), MockUserAccessHash("hash2"))
	require.Error(t, err)
}

func TestDeclineFriendRequest_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/friends/decline").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"userId":2, "userAccessHash":"hash2"}`).
		Reply(200)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := client.DeclineFriendRequest(context.Background(), auth, MockUserId(2), MockUserAccessHash("hash2"))
	require.NoError(t, err)
}

func TestDeclineFriendRequest_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/friends/decline").
		Reply(400)

	client := NewClient()
	err := client.DeclineFriendRequest(context.Background(), nil, MockUserId(2), MockUserAccessHash("hash2"))
	require.Error(t, err)
}
