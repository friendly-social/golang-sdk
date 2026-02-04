package sdk

import (
	"context"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestAuthValueTypes(t *testing.T) {
	t.Run("Valid Token", func(t *testing.T) {
		token, err := NewToken(strings.Repeat("1", 256))
		require.EqualValues(t, Token(strings.Repeat("1", 256)), token)
		require.NoError(t, err)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := NewToken("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTokenLengthMustBe256)
	})

	t.Run("Valid UserAccesssHash", func(t *testing.T) {
		hash, err := NewUserAccessHash(strings.Repeat("1", 256))
		require.EqualValues(t, UserAccessHash(strings.Repeat("1", 256)), hash)
		require.NoError(t, err)
	})

	t.Run("Invalid UserAccessHash", func(t *testing.T) {
		_, err := NewUserAccessHash("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserAccessHashLengthMustBe256)
	})
}

func TestGenerate_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/auth/generate").
		JSON(`{"nickname":"atennop", "description":"bio","interests":["programming"],"avatar":{"id":10,"accessHash":"hash"},"socialLink":"https://github.com/Atennop1"}`).
		Reply(200).
		JSON(`{"id":1,"token":"token","accessHash":"hash"}`)

	client := NewClient("https://getfriend.ly")
	auth, err := client.Generate(context.Background(), "atennop", "bio", []Interest{"programming"}, &FileDescriptor{Id: 10, AccessHash: "hash"}, "https://github.com/Atennop1")

	require.NoError(t, err)
	require.Equal(t, &Authorization{
		Id:         1,
		Token:      Token("token"),
		AccessHash: UserAccessHash("hash"),
	}, auth)
}

func TestGenerate_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/auth/generate").
		Reply(400)

	client := NewClient("https://getfriend.ly")
	_, err := client.Generate(context.Background(), "atennop", "bio", []Interest{"programming"}, &FileDescriptor{Id: 10, AccessHash: "hash"}, "https://github.com/Atennop1")
	require.Error(t, err)
}
