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
		expectedError error
	}

	cases := []testCase{
		{
			name:         "Success",
			auth:         &Authorization{Id: 1, Token: "token", AccessHash: "hash"},
			mockStatus:   200,
			mockResponse: `{"entries":[{"isExtendedNetwork":true,"commonFriends":[],"details":{"id":2},"isRequest":true}]}`,
			expectedFeed: &FeedQueue{
				Entries: []FeedEntry{{
						IsRequest: true,
						IsExtendedNetwork: true,
						CommonFriends:     []UserDetails{},
						Details:           UserDetails{Id: 2},
					}},
			},
		},
		{
			name:          "API Error",
			auth:          &Authorization{Id: 1, Token: "token", AccessHash: "hash"},
			mockStatus:    500,
			expectedError: ErrInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("https://getfriend.ly").
				Get("/feed/queue").
				Reply(tc.mockStatus).
				SetHeader("Content-Type", "application/json").
				SetHeader("X-User-Id", "1").
				SetHeader("X-Token", "token").
				JSON(tc.mockResponse)

			client := NewClient("https://getfriend.ly")
			feed, err := client.GetFeedQueue(tc.auth)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedFeed, feed)
			}
		})
	}
}
