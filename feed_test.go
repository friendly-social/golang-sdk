package sdk

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestGetFeedQueue(t *testing.T) {
	type testCase struct {
		name          string
		auth          *Authorization
		mockStatus    int
		mockResponse  string
		expectedFeed  *FeedQueue
		expectError   bool
	}

	cases := []testCase{
		{
			name:         "Success",
			auth:         &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")},
			mockStatus:   200,
			mockResponse: `{"entries":[{"isExtendedNetwork":true,"commonFriends":[],"details":{"id":2},"isRequest":true}]}`,
			expectedFeed: &FeedQueue{
				Entries: []FeedEntry{{
					IsRequest:         true,
					IsExtendedNetwork: true,
					CommonFriends:     []UserDetails{},
					Details:           UserDetails{Id: 2},
				}},
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
				Get("/feed/queue").
				MatchHeader("Content-Type", "application/json").
				MatchHeader("X-User-Id", "1").
				MatchHeader("X-Token", "token").
				Reply(tc.mockStatus).
				JSON(tc.mockResponse)

			client := NewClient("https://getfriend.ly")
			feed, err := client.GetFeedQueue(tc.auth)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedFeed, feed)
			}
		})
	}
}
