package sdk

import (
	"strings"
	"testing"

	"github.com/h2non/gock"
)

func TestAuthValueTypes(t *testing.T) {
	RunValueTypeTest(t, ValueTypeTestCase[string, Token]{
		name:    "Token",
		new:     NewToken,
		valid:   strings.Repeat("1", 256),
		invalid: "1",
	})

	RunValueTypeTest(t, ValueTypeTestCase[string, UserAccessHash]{
		name:    "UserAccessHash",
		new:     NewUserAccessHash,
		valid:   strings.Repeat("1", 256),
		invalid: "1",
	})
}

func TestGenerate(t *testing.T) {
	f := func() *gock.Request {
		return gock.New("https://getfriend.ly").Post("/auth/generate")
	}

	call := func(c *Client) (*Authorization, error) {
		return c.Generate(MustNickname("atennop"), MustUserDescription("something"), []Interest{}, nil)
	}

	RunAPITests(t, append(CommonCasesTests(f, call),
		APITestCase[Authorization]{
			name:             "Success",
			setup:            func() { f().Reply(200).JSON(generateResponse{Id: 1, AccessHash: "1", Token: "1"}) },
			call:             call,
			expectError:      false,
			expectedResponse: &Authorization{Id: 1, AccessHash: "1", Token: "1"},
		},
	))
}
