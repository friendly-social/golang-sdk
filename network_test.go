package sdk

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestGetNetworkDetails(t *testing.T) {
	type testCase struct {
		name            string
		auth            *Authorization
		mockStatus      int
		mockResponse    string
		expectedNetwork *NetworkDetails
		expectError     bool
	}

	cases := []testCase{
		{
			name:         "Success",
			auth:         &Authorization{Id: 1, Token: Token("token"), AccessHash: UserAccessHash("hash")},
			mockStatus:   200,
			mockResponse: `{"friends":[{"id":2,"accessHash":"hash2","nickname":"tr3ble","description":"something2","interests":["mac"],"avatar":{"id":3,"accessHash":"hash3"}}]}`,
			expectedNetwork: &NetworkDetails{
				Friends: []UserDetails{{
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
				Get("/network/details").
				MatchHeader("Content-Type", "application/json").
				MatchHeader("X-User-Id", "1").
				MatchHeader("X-Token", "token").
				Reply(tc.mockStatus).
				JSON(tc.mockResponse)

			client := NewClient("https://getfriend.ly")
			network, err := client.GetNetworkDetails(tc.auth)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedNetwork, network)
			}
		})
	}
}
