package sdk

import (
	"fmt"
)

// FriendToken is a token by which other users can add Token's owner to their friend list.
type FriendToken string

type addFriendRequest struct {
	Token  FriendToken `json:"token"`
	UserId UserId      `json:"userId"`
}

type friendRequestRequest struct {
	UserId     UserId         `json:"userId"`
	AccessHash UserAccessHash `json:"userAccessHash"`
}

type generateFriendTokenResponse struct {
	Token FriendToken `json:"token"`
}

type addFriendResponse struct {
	Type string `json:"type"`
}

// NewFriendToken creates new FriendToken or returns an error if tokens length isn't 256.
func NewFriendToken(s string) (FriendToken, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("friend token must be 256 characters, got %d", len(s))
	}

	return FriendToken(s), nil
}

// GenerateFriendToken creates token for Authorization's user by which another users can add them.
func (c *Client) GenerateFriendToken(auth *Authorization) (FriendToken, error) {
	var resp generateFriendTokenResponse
	err := c.do("POST", "/friends/generate", auth, nil, &resp)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

// AddFriend makes request to add user with provided FriendToken and ID to Authorization's friends list.
func (c *Client) AddFriend(auth *Authorization, token FriendToken, userId UserId) error {
	req := addFriendRequest{
		Token:  token,
		UserId: userId,
	}

	var resp addFriendResponse
	err := c.do("POST", "/friends/add", auth, req, nil)
	if err != nil {
		return err
	}

	if resp.Type == "FriendTokenExpired" {
		return fmt.Errorf("friend token expired")
	}

	return nil
}

// SendFriendRequest sends friend request from Authorization to user with provided ID and AccessHash.
func (c *Client) SendFriendRequest(auth *Authorization, userId UserId, accessHash UserAccessHash) error {
	req := friendRequestRequest{
		UserId:     userId,
		AccessHash: accessHash,
	}

	err := c.do("POST", "/friends/request", auth, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeclineFriendRequest declines Authorization's request from user with provided ID and AccessHash.
func (c *Client) DeclineFriendRequest(auth *Authorization, userId UserId, accessHash UserAccessHash) error {
	req := friendRequestRequest{
		UserId:     userId,
		AccessHash: accessHash,
	}

	err := c.do("POST", "/friends/decline", auth, req, nil)
	if err != nil {
		return err
	}

	return nil
}
