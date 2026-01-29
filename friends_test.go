package sdk

import (
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

func TestGenerateFriendToken(t *testing.T) {
	type testCase struct {
		name                string
		auth                *Authorization
		mockStatus          int
		mockResponse        string
		expectedFriendToken FriendToken
		expectError         bool
	}

	cases := []testCase{
		{
			name:                "Success",
			auth:                &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")},
			mockStatus:          200,
			mockResponse:        `{"token":"token2"}`,
			expectedFriendToken: FriendToken("token2"),
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
				Post("/friends/generate").
				MatchHeader("Content-Type", "application/json").
				MatchHeader("X-User-Id", "1").
				MatchHeader("X-Token", "token").
				Reply(tc.mockStatus).
				JSON(tc.mockResponse)

			client := NewClient("https://getfriend.ly")
			token, err := client.GenerateFriendToken(tc.auth)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedFriendToken, token)
			}
		})
	}
}

func TestAddFriend(t *testing.T) {
	type input struct {
		auth  *Authorization
		id    UserId
		token FriendToken
	}

	type testCase struct {
		name          string
		input         input
		mockStatus    int
		mockResponse  string
		expectedBody  string
		expectError   bool
		expectedError error
	}

	cases := []testCase{
		{
			name: "Success",
			input: input{
				auth:  &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")},
				id:    2,
				token: FriendToken("token2"),
			},
			mockStatus:   200,
			mockResponse: `{"type":"nice"}`,
			expectedBody: `{"userId":2, "token":"token2"}`,
		},
		{
			name:          "Friend Token Expired",
			input:         input{auth: &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}},
			mockResponse:  `{"type":"FriendTokenExpired"}`,
			mockStatus:    200,
			expectedError: ErrFriendTokenExpired,
		},
		{
			name:        "API Error",
			input:       input{auth: &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}},
			mockStatus:  400,
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("https://getfriend.ly").
				Post("/friends/add").
				MatchHeader("Content-Type", "application/json").
				MatchHeader("X-User-Id", "1").
				MatchHeader("X-Token", "token").
				JSON(tc.expectedBody).
				Reply(tc.mockStatus).
				JSON(tc.mockResponse)

			client := NewClient("https://getfriend.ly")
			err := client.AddFriend(tc.input.auth, tc.input.token, tc.input.id)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendFriendRequest(t *testing.T) {
	type input struct {
		auth *Authorization
		id   UserId
		hash UserAccessHash
	}

	type testCase struct {
		name          string
		input         input
		mockStatus    int
		expectedBody  string
		expectError   bool
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
			expectedBody: `{"userId":2, "userAccessHash":"hash2"}`,
		},
		{
			name:        "API Error",
			input:       input{auth: &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}},
			mockStatus:  400,
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("https://getfriend.ly").
				Post("/friends/request").
				MatchHeader("Content-Type", "application/json").
				MatchHeader("X-User-Id", "1").
				MatchHeader("X-Token", "token").
				JSON(tc.expectedBody).
				Reply(tc.mockStatus)

			client := NewClient("https://getfriend.ly")
			err := client.SendFriendRequest(tc.input.auth, tc.input.id, tc.input.hash)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeclineFriendRequest(t *testing.T) {
	type input struct {
		auth *Authorization
		id   UserId
		hash UserAccessHash
	}

	type testCase struct {
		name          string
		input         input
		mockStatus    int
		expectedBody  string
		expectError   bool
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
			expectedBody: `{"userId":2, "userAccessHash":"hash2"}`,
		},
		{
			name:        "API Error",
			input:       input{auth: &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")}},
			mockStatus:  400,
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("https://getfriend.ly").
				Post("/friends/decline").
				MatchHeader("Content-Type", "application/json").
				MatchHeader("X-User-Id", "1").
				MatchHeader("X-Token", "token").
				JSON(tc.expectedBody).
				Reply(tc.mockStatus)

			client := NewClient("https://getfriend.ly")
			err := client.DeclineFriendRequest(tc.input.auth, tc.input.id, tc.input.hash)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
