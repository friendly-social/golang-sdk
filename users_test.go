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

func TestGetSelfDetails(t *testing.T) {
	type testCase struct {
		name            string
		auth            *Authorization
		mockStatus      int
		mockResponse    string
		expectedDetails *UserDetails
		expectError     bool
	}

	cases := []testCase{
		{
			name:         "Success",
			auth:         &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")},
			mockStatus:   200,
			mockResponse: `{"id":1,"accessHash":"hash","nickname":"atennop","description":"something","interests":["vim"],"avatar":{"id":2,"accessHash":"hash2"}}`,
			expectedDetails: &UserDetails{
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
			},
		},
		{
			name:        "API Error",
			auth:        &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")},
			mockStatus:  500,
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("https://getfriend.ly").
				Get("/users/details").
				MatchHeader("Content-Type", "application/json").
				MatchHeader("X-User-Id", "1").
				MatchHeader("X-Token", "token").
				Reply(tc.mockStatus).
				JSON(tc.mockResponse)

			client := NewClient("https://getfriend.ly")
			self, err := client.GetSelfDetails(context.Background(), tc.auth)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedDetails, self)
			}
		})
	}
}

func TestGetUserDetails(t *testing.T) {
	type input struct {
		auth *Authorization
		id   UserId
		hash UserAccessHash
	}

	type testCase struct {
		name            string
		input           input
		mockStatus      int
		mockResponse    string
		expectedDetails *UserDetails
		expectError     bool
	}

	cases := []testCase{
		{
			name: "Success",
			input: input{
				auth: &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")},
				id:   2,
				hash: UserAccessHash("hash2"),
			},
			mockStatus:   200,
			mockResponse: `{"id":2,"accessHash":"hash2","nickname":"tr3ble","description":"something2","interests":["mac"],"avatar":{"id":3,"accessHash":"hash3"}}`,
			expectedDetails: &UserDetails{
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
			},
		},
		{
			name: "API Error",
			input: input{
				auth: &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")},
				id:   2,
				hash: UserAccessHash("hash2"),
			},
			mockStatus:  500,
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("https://getfriend.ly").
				Get("/users/details/2/hash2").
				MatchHeader("Content-Type", "application/json").
				MatchHeader("X-User-Id", "1").
				MatchHeader("X-Token", "token").
				Reply(tc.mockStatus).
				JSON(tc.mockResponse)

			client := NewClient("https://getfriend.ly")
			user, err := client.GetUserDetails(context.Background(), tc.input.auth, tc.input.id, tc.input.hash)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedDetails, user)
			}
		})
	}
}
