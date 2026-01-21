package sdk

import (
	"testing"

	"github.com/h2non/gock"
)

func TestGetNetwork(t *testing.T) {
	f := func() *gock.Request {
		return gock.New("https://getfriend.ly").Get("/network/details")
	}

	call := func(c *Client) (*NetworkDetails, error) {
		return c.GetNetworkDetails(&Authorization{})
	}

	RunAPITests(t, append(CommonCasesTests(f, call),
		APITestCase[NetworkDetails]{
			name:             "Success",
			setup:            func() { f().Reply(200).JSON(NetworkDetails{Friends: []UserDetails{{Nickname: "atennop"}}}) },
			call:             call,
			expectError:      false,
			expectedResponse: &NetworkDetails{Friends: []UserDetails{{Nickname: "atennop"}}},
		},
	))
}
