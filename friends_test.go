package sdk

import (
	"errors"
	"strings"
	"testing"

	"github.com/h2non/gock"
)

func TestFriendsValueTypes(t *testing.T) {
	RunValueTypeTest(t, ValueTypeTestCase[string, FriendToken]{
		name:    "FriendToken",
		new:     NewFriendToken,
		valid:   strings.Repeat("1", 256),
		invalid: "1",
	})
}

func TestGenerateFriendToken(t *testing.T) {
	f := func() *gock.Request {
		return gock.New("https://getfriend.ly").Post("/friends/generate")
	}

	call := func(c *Client) (*FriendToken, error) {
		token, err := c.GenerateFriendToken(&Authorization{Id: 1, AccessHash: "1", Token: "1"})
		return &token, err
	}

	RunAPITests(t, append(CommonCasesTests(f, call),
		APITestCase[FriendToken]{
			name:        "Success",
			setup:       func() { f().Reply(200).JSON(generateFriendTokenResponse{Token: "123"}) },
			call:        call,
			expectError: false,
			expectedResponse: func() *FriendToken {
				token := FriendToken("123")
				return &token
			}(),
		},
	))
}

func TestAddFriend(t *testing.T) {
	f := func() *gock.Request {
		return gock.New("https://getfriend.ly").Post("/friends/add")
	}

	call := func(c *Client) (*any, error) {
		err := c.AddFriend(&Authorization{Id: 1, AccessHash: "1", Token: "1"}, FriendToken("1"), UserId(1))
		return nil, err
	}

	RunAPITests(t, append(CommonCasesTests(f, call),
		APITestCase[any]{
			name:             "Success",
			setup:            func() { f().Reply(200).JSON(addFriendResponse{Type: "123"}) },
			call:             call,
			expectError:      false,
			expectedResponse: nil,
		},
		APITestCase[any]{
			name:             "Token Expired",
			setup:            func() { f().Reply(200).JSON(addFriendResponse{Type: "FriendTokenExpired"}) },
			call:             call,
			expectError:      true,
			expectedResponse: nil,
		},
	))
}

func TestSendFriendRequest(t *testing.T) {
	f := func() *gock.Request {
		return gock.New("https://getfriend.ly").Post("/friends/request")
	}

	call := func(c *Client) (*any, error) {
		err := c.SendFriendRequest(&Authorization{Id: 1, AccessHash: "1", Token: "1"}, UserId(1), UserAccessHash("123"))
		return nil, err
	}

	RunAPITests(t, []APITestCase[any]{
		{
			name:             "Success",
			setup:            func() { f().Reply(200) },
			call:             call,
			expectError:      false,
			expectedResponse: nil,
		},
		{
			name:        "User Not Found",
			setup:       func() { f().Reply(404) },
			call:        call,
			expectError: true,
		},
		{
			name:        "Unauthorized",
			setup:       func() { f().Reply(401) },
			call:        call,
			expectError: true,
		},
		{
			name:        "Server Error",
			setup:       func() { f().Reply(500) },
			call:        call,
			expectError: true,
		},
		{
			name:        "Network Fail",
			setup:       func() { f().ReplyError(errors.New("something went wrong")) },
			call:        call,
			expectError: true,
		},
	})
}

func TestDeclineFriendRequest(t *testing.T) {
	f := func() *gock.Request {
		return gock.New("https://getfriend.ly").Post("/friends/decline")
	}

	call := func(c *Client) (*any, error) {
		err := c.DeclineFriendRequest(&Authorization{Id: 1, AccessHash: "1", Token: "1"}, UserId(1), UserAccessHash("123"))
		return nil, err
	}

	RunAPITests(t, []APITestCase[any]{
		{
			name:             "Success",
			setup:            func() { f().Reply(200) },
			call:             call,
			expectError:      false,
			expectedResponse: nil,
		},
		{
			name:        "User Not Found",
			setup:       func() { f().Reply(404) },
			call:        call,
			expectError: true,
		},
		{
			name:        "Unauthorized",
			setup:       func() { f().Reply(401) },
			call:        call,
			expectError: true,
		},
		{
			name:        "Server Error",
			setup:       func() { f().Reply(500) },
			call:        call,
			expectError: true,
		},
		{
			name:        "Network Fail",
			setup:       func() { f().ReplyError(errors.New("something went wrong")) },
			call:        call,
			expectError: true,
		},
	})
}
