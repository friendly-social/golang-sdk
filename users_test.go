package sdk

import (
	"context"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestUsersValueTypes(t *testing.T) {
	t.Run("Valid Nickname", func(t *testing.T) {
		nickname, err := NewNickname("atennop")
		require.EqualValues(t, Nickname("atennop"), nickname)
		require.NoError(t, err)
	})

	t.Run("Invalid Nickname", func(t *testing.T) {
		_, err := NewNickname(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNicknameLengthMustBeLessThan256)
	})

	t.Run("Valid Description", func(t *testing.T) {
		desc, err := NewUserDescription("something")
		require.EqualValues(t, UserDescription("something"), desc)
		require.NoError(t, err)
	})

	t.Run("Invalid Description", func(t *testing.T) {
		_, err := NewUserDescription(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserDescriptionLengthMustBeLessThan1024)
	})

	t.Run("Valid Interest", func(t *testing.T) {
		interest, err := NewInterest("vim")
		require.EqualValues(t, Interest("vim"), interest)
		require.NoError(t, err)
	})

	t.Run("Invalid Interest", func(t *testing.T) {
		_, err := NewInterest(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInterestLengthMustBeLessThan64)
	})

	t.Run("Valid SocialLink", func(t *testing.T) {
		link, err := NewSocialLink("https://github.com/Atennop1")
		require.EqualValues(t, SocialLink("https://github.com/Atennop1"), link)
		require.NoError(t, err)
	})

	t.Run("Invalid SocialLink", func(t *testing.T) {
		_, err := NewSocialLink(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrSocialLinkLengthMustBeLessThan2048)
	})
}

func TestGetSelfDetails_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/users/details").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"id":1,"accessHash":"hash","nickname":"atennop","description":"something","interests":["vim"],"avatar":{"id":2,"accessHash":"hash2"}}`)

	client := NewClient("https://getfriend.ly")
	self, err := client.GetSelfDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})

	require.NoError(t, err)
	require.Equal(t, &UserDetails{
		Id:          1,
		AccessHash:  UserAccessHash("hash"),
		Nickname:    Nickname("atennop"),
		Description: UserDescription("something"),
		Interests: []Interest{
			Interest("vim"),
		},
		Avatar: &FileDescriptor{
			Id:         2,
			AccessHash: FileAccessHash("hash2"),
		},
	}, self)
}

func TestGetSelfDetails_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/users/details").
		Reply(400)

	client := NewClient("https://getfriend.ly")
	_, err := client.GetSelfDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})
	require.Error(t, err)
}

func TestGetUserDetails_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/users/details/2/hash2").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"id":2,"accessHash":"hash2","nickname":"tr3ble","description":"something2","interests":["mac"],"avatar":{"id":3,"accessHash":"hash3"}}`)

	client := NewClient("https://getfriend.ly")
	user, err := client.GetUserDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, 2, UserAccessHash("hash2"))

	require.NoError(t, err)
	require.Equal(t, &UserDetails{
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
	}, user)
}

func TestGetUserDetails_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/users/details/2/hash2").
		Reply(400)

	client := NewClient("https://getfriend.ly")
	_, err := client.GetUserDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, 2, UserAccessHash("hash2"))
	require.Error(t, err)
}
