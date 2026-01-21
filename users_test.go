package sdk

import (
	"strings"
	"testing"

	"github.com/h2non/gock"
)

func TestUsersValueTypes(t *testing.T) {
	RunValueTypeTest(t, ValueTypeTestCase[string, Nickname]{
		name:    "Nickname",
		new:     NewNickname,
		must:    MustNickname,
		valid:   "atennop",
		invalid: strings.Repeat("1", 1024),
	})

	RunValueTypeTest(t, ValueTypeTestCase[string, UserDescription]{
		name:    "UserDescription",
		new:     NewUserDescription,
		must:    MustUserDescription,
		valid:   "something",
		invalid: strings.Repeat("1", 2048),
	})

	RunValueTypeTest(t, ValueTypeTestCase[string, Interest]{
		name:    "Interest",
		new:     NewInterest,
		must:    MustInterest,
		valid:   "programming",
		invalid: strings.Repeat("1", 256),
	})
}

func TestGetSelfDetails(t *testing.T) {
	f := func() *gock.Request {
		return gock.New("https://getfriend.ly").Get("/users/details")
	}

	call := func(c *Client) (*UserDetails, error) {
		return c.GetSelfDetails(&Authorization{})
	}

	RunAPITests(t, append(CommonCasesTests(f, call),
		APITestCase[UserDetails]{
			name:             "Success",
			setup:            func() { f().Reply(200).JSON(UserDetails{Id: 1, Nickname: "atennop"}) },
			call:             call,
			expectError:      false,
			expectedResponse: &UserDetails{Id: 1, Nickname: "atennop"},
		},
	))
}

func TestGetUserDetails(t *testing.T) {
	f := func() *gock.Request {
		return gock.New("https://getfriend.ly").Get("/users/details/1/222")
	}

	call := func(c *Client) (*UserDetails, error) {
		return c.GetUserDetails(&Authorization{}, 1, "222")
	}

	RunAPITests(t, append(CommonCasesTests(f, call),
		APITestCase[UserDetails]{
			name:             "Success",
			setup:            func() { f().Reply(200).JSON(UserDetails{Id: 1, Nickname: "atennop"}) },
			call:             call,
			expectError:      false,
			expectedResponse: &UserDetails{Id: 1, Nickname: "atennop"},
		},
	))
}
