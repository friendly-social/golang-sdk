package sdk

import (
	"context"
	"slices"
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
		require.ErrorIs(t, err, ErrTooLongNickname)
	})

	t.Run("Valid Description", func(t *testing.T) {
		desc, err := NewUserDescription("something")
		require.EqualValues(t, UserDescription("something"), desc)
		require.NoError(t, err)
	})

	t.Run("Invalid Description", func(t *testing.T) {
		_, err := NewUserDescription(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrToLongUserDescription)
	})

	t.Run("Valid Interest", func(t *testing.T) {
		interest, err := NewInterest("vim")
		require.EqualValues(t, Interest("vim"), interest)
		require.NoError(t, err)
	})

	t.Run("Invalid Interest", func(t *testing.T) {
		_, err := NewInterest(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongInterest)
	})

	t.Run("Valid Interests", func(t *testing.T) {
		interests, err := NewInterests(Interest("vim"))
		require.EqualValues(t, Interests{"vim"}, interests)
		require.NoError(t, err)
	})

	t.Run("Invalid Interests", func(t *testing.T) {
		_, err := NewInterests(slices.Repeat([]Interest{Interest("vim")}, 1000)...)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooMuchInterests)
	})

	t.Run("Valid SocialLink", func(t *testing.T) {
		link, err := NewSocialLink("https://github.com/Atennop1")
		require.EqualValues(t, SocialLink("https://github.com/Atennop1"), link)
		require.NoError(t, err)
	})

	t.Run("Invalid SocialLink", func(t *testing.T) {
		_, err := NewSocialLink(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongSocialLink)
	})
}

func TestGetSelfDetails_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/users/details").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"id":1,"accessHash":"hash","nickname":"atennop","description":"something","interests":["vim"],"avatar":{"id":2,"accessHash":"hash2"}}`)

	client := NewClient("https://api.getfriend.ly")
	self, err := client.GetSelfDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})

	require.NoError(t, err)
	require.Equal(t, &UserDetails{
		Id:          1,
		AccessHash:  UserAccessHash("hash"),
		Nickname:    Nickname("atennop"),
		Description: UserDescription("something"),
		Interests: Interests{
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

	gock.New("https://api.getfriend.ly").
		Get("/users/details").
		Reply(400)

	client := NewClient("https://api.getfriend.ly")
	_, err := client.GetSelfDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")})
	require.Error(t, err)
}

func TestGetUserDetails_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/users/details/2/hash2").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"id":2,"accessHash":"hash2","nickname":"tr3ble","description":"something2","interests":["mac"],"avatar":{"id":3,"accessHash":"hash3"}}`)

	client := NewClient("https://api.getfriend.ly")
	user, err := client.GetUserDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, 2, UserAccessHash("hash2"))

	require.NoError(t, err)
	require.Equal(t, &UserDetails{
		Id:          2,
		AccessHash:  UserAccessHash("hash2"),
		Nickname:    Nickname("tr3ble"),
		Description: UserDescription("something2"),
		Interests: Interests{
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

	gock.New("https://api.getfriend.ly").
		Get("/users/details/2/hash2").
		Reply(400)

	client := NewClient("https://api.getfriend.ly")
	_, err := client.GetUserDetails(context.Background(), &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}, 2, UserAccessHash("hash2"))
	require.Error(t, err)
}

func TestOptions(t *testing.T) {
	cases := []struct {
		name         string
		option       editAccountOption
		expectedBody string
	}{
		{
			name:         "Valid Nickname Option",
			option:       WithUserNickname("atennop"),
			expectedBody: `{"nickname":{"value":"atennop"}}`,
		},
		{
			name:         "Invalid Nickname Option",
			option:       WithUserNickname(""),
			expectedBody: `{}`,
		},
		{
			name:         "Valid Description Option",
			option:       WithUserDescription("bio"),
			expectedBody: `{"description":{"value":"bio"}}`,
		},
		{
			name:         "Invalid Description Option",
			option:       WithUserDescription(""),
			expectedBody: `{}`,
		},
		{
			name:         "Valid Interests Option",
			option:       WithUserInterests(Interests{"neovim", "coding"}),
			expectedBody: `{"interests":{"value":["neovim", "coding"]}}`,
		},
		{
			name:         "Invalid Interests Option",
			option:       WithUserInterests(nil),
			expectedBody: `{}`,
		},
		{
			name:         "Valid Avatar Option",
			option:       WithUserAvatar(&FileDescriptor{Id: 10, AccessHash: "hash"}),
			expectedBody: `{"avatar":{"value":{"id":10, "accessHash":"hash"}}}`,
		},
		{
			name:         "Invalid Avatar Option",
			option:       WithUserAvatar(nil),
			expectedBody: `{}`,
		},
		{
			name:         "Valid SocialLink Option",
			option:       WithUserSocialLink("https://example.com"),
			expectedBody: `{"socialLink":{"value":"https://example.com"}}`,
		},
		{
			name:         "Invalid SocialLink Option",
			option:       WithUserSocialLink(""),
			expectedBody: `{}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("https://api.getfriend.ly").
				JSON(tc.expectedBody).
				Reply(200)

			c := NewClient("https://api.getfriend.ly")
			err := c.EditAccount(context.Background(), &Authorization{Id: 1, Token: Token("token")}, tc.option)
			require.NoError(t, err)
		})
	}
}

func TestEditAccount_Empty(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Patch("/users/edit").
		Reply(200)

	c := NewClient("https://api.getfriend.ly")
	err := c.EditAccount(context.Background(), &Authorization{Id: 1, Token: Token("token")})

	require.NoError(t, err)
	require.False(t, gock.IsDone())
}

func TestEditAccount_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Patch("/users/edit").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"nickname":{"value":"atennop"},"description":{"value":"bio"}}`).
		Reply(200)

	c := NewClient("https://api.getfriend.ly")
	err := c.EditAccount(context.Background(), &Authorization{Id: 1, Token: Token("token")}, WithUserNickname("atennop"), WithUserDescription("bio"))
	require.NoError(t, err)
}

func TestEditAccount_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Patch("/users/edit").
		Reply(400)

	c := NewClient("https://api.getfriend.ly")
	err := c.EditAccount(context.Background(), &Authorization{Id: 1, Token: Token("token")}, WithUserNickname("atennop"), WithUserDescription("bio"))
	require.Error(t, err)
}
