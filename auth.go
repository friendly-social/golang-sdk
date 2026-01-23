package sdk

import (
	"fmt"
)

// UserId represents the unique identifier of user.
type UserId int64

// UserAccessHash represents the unique hash associated with user. Works in pair with UserId.
type UserAccessHash string

// Token represents access token for the user. Works in pair with UserId.
type Token string

// Authorization is a helper structure for composing user's ID, AccessHash and Token for authorization.
type Authorization struct {
	Id         UserId         `json:"id"`
	AccessHash UserAccessHash `json:"accessHash"`
	Token      Token          `json:"token"`
}

type generateRequest struct {
	Nickname    Nickname        `json:"nickname"`
	Description UserDescription `json:"description"`
	Interests   []Interest      `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
	SocialLink  SocialLink      `json:"socialLink"`
}

type generateResponse struct {
	Id         UserId         `json:"id"`
	AccessHash UserAccessHash `json:"accessHash"`
	Token      Token          `json:"token"`
}

var (
	ErrTokenMustBe256CharactersLength          = fmt.Errorf("token must be 256 characters lenght")
	ErrUserAccessHashMustBe256CharactersLength = fmt.Errorf("user access hash must be 256 characters lenght")
)

// NewToken creates new Token or returns an error if token's length isn't 256.
func NewToken(s string) (Token, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("length is %d: %w", len(s), ErrTokenMustBe256CharactersLength)
	}

	return Token(s), nil
}

// NewUserAccessHash creates new UserAccessHash or returns an error if hash length isn't 256.
func NewUserAccessHash(s string) (UserAccessHash, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("length is %d: %w", len(s), ErrUserAccessHashMustBe256CharactersLength)
	}

	return UserAccessHash(s), nil
}

// Generate makes request for creating account using provided data and returns Authorization structure.
func (c *Client) Generate(nickname Nickname, description UserDescription, interests []Interest, avatar *FileDescriptor, link SocialLink) (*Authorization, error) {
	req := generateRequest{
		Nickname:    nickname,
		Description: description,
		Interests:   interests,
		Avatar:      avatar,
		SocialLink:  link,
	}

	var resp generateResponse
	err := c.do("POST", "/auth/generate", nil, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("request to /auth/generate failed: %w", err)
	}

	return &Authorization{
		Id:         resp.Id,
		AccessHash: resp.AccessHash,
		Token:      resp.Token,
	}, nil
}
