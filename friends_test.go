package sdk

import (
	"context"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestFriendValueTypes(t *testing.T) {
	t.Run("Valid FriendToken", func(t *testing.T) {
		token, err := NewFriendToken(strings.Repeat("1", 256))
		require.EqualValues(t, FriendToken(strings.Repeat("1", 256)), token)
		require.NoError(t, err)
	})

	t.Run("Invalid FriendToken", func(t *testing.T) {
		_, err := NewFriendToken("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrFriendTokenLengthMustBe256)
	})
}

func TestGenerateFriendToken_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/friends/generate").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"token":"token2"}`)

	client := NewClient("https://getfriend.ly")
	token, err := client.GenerateFriendToken(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})

	require.NoError(t, err)
	require.Equal(t, FriendToken("token2"), token)
}

func TestGenerateFriendToken_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/friends/generate").
		Reply(400)

	client := NewClient("https://getfriend.ly")
	_, err := client.GenerateFriendToken(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})
	require.Error(t, err)
}

func TestAddFriend_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/friends/add").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"userId":2, "token":"token2"}`).
		Reply(200).
		JSON(`{"type":"nice"}`)

	client := NewClient("https://getfriend.ly")
	err := client.AddFriend(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, "token2", 2)
	require.NoError(t, err)
}

func TestAddFriend_TokenExpired(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/friends/add").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"userId":2, "token":"token2"}`).
		Reply(200).
		JSON(`{"type":"FriendTokenExpired"}`)

	client := NewClient("https://getfriend.ly")
	err := client.AddFriend(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, "token2", 2)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrFriendTokenExpired)
}

func TestAddFriend_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/friends/add").
		Reply(400)

	client := NewClient("https://getfriend.ly")
	err := client.AddFriend(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, "token2", 2)
	require.Error(t, err)
}

func TestSendFriendRequest_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/friends/request").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"userId":2, "userAccessHash":"hash2"}`).
		Reply(200)

	client := NewClient("https://getfriend.ly")
	err := client.SendFriendRequest(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, 2, "hash2")
	require.NoError(t, err)
}

func TestSendFriendRequest_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/friends/request").
		Reply(400)

	client := NewClient("https://getfriend.ly")
	err := client.SendFriendRequest(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, 2, "hash2")
	require.Error(t, err)
}

func TestDeclineFriendRequest_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/friends/decline").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"userId":2, "userAccessHash":"hash2"}`).
		Reply(200)

	client := NewClient("https://getfriend.ly")
	err := client.DeclineFriendRequest(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, 2, "hash2")
	require.NoError(t, err)
}

func TestDeclineFriendRequest_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/friends/decline").
		Reply(400)

	client := NewClient("https://getfriend.ly")
	err := client.DeclineFriendRequest(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, 2, "hash2")
	require.Error(t, err)
}
