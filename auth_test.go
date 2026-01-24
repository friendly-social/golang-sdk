package sdk

import (
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
		require.ErrorIs(t, err, ErrTokenMustBe256CharactersLength)
	})

	t.Run("Valid UserAccesssHash", func(t *testing.T) {
		hash, err := NewUserAccessHash(strings.Repeat("1", 256))
		require.EqualValues(t, UserAccessHash(strings.Repeat("1", 256)), hash)
		require.NoError(t, err)
	})

	t.Run("Invalid UserAccessHash", func(t *testing.T) {
		_, err := NewUserAccessHash("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserAccessHashMustBe256CharactersLength)
	})
}

func TestGenerate(t *testing.T) {
	type input struct {
		nickname    Nickname
		description UserDescription
		interests   []Interest
		avatar      *FileDescriptor
		link        SocialLink
	}

	type testCase struct {
		name          string
		input         input
		mockStatus    int
		mockResponse  string
		expectedBody  string
		expectedAuth  *Authorization
		expectedError error
	}

	cases := []testCase{
		{
			name: "Success",
			input: input{
				nickname:    Nickname("atennop"),
				description: UserDescription("bio"),
				interests:   []Interest{Interest("programming")},
				avatar:      &FileDescriptor{Id: 10, AccessHash: FileAccessHash("hash")},
				link:        SocialLink("https://github.com/Atennop1"),
			},
			mockStatus:   200,
			mockResponse: `{"id":1,"token":"token","accessHash":"hash"}`,
			expectedBody: `{"nickname":"atennop", "description":"bio","interests":["programming"],"avatar":{"id":10,"accessHash":"hash"},"socialLink":"https://github.com/Atennop1"}`,
			expectedAuth: &Authorization{
				Id:         1,
				Token:      Token("token"),
				AccessHash: UserAccessHash("hash"),
			},
		},
		{
			name:          "API Error",
			mockStatus:    400,
			expectedError: ErrRequestFailed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("https://getfriend.ly").
				Post("/auth/generate").
				JSON(tc.expectedBody).
				Reply(tc.mockStatus).
				JSON(tc.mockResponse)

			client := NewClient("https://getfriend.ly")
			auth, err := client.Generate(tc.input.nickname, tc.input.description, tc.input.interests, tc.input.avatar, tc.input.link)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedAuth, auth)
			}
		})
	}
}
