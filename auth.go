package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
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
}

type generateResponse struct {
	Id         UserId         `json:"id"`
	AccessHash UserAccessHash `json:"accessHash"`
	Token      Token          `json:"token"`
}

// NewToken creates new Token or returns an error if token's length isn't 256.
func NewToken(s string) (Token, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("token must be 256 characters, got %d", len(s))
	}

	return Token(s), nil
}

// NewUserAccessHash creates new UserAccessHash or returns an error if hash length isn't 256.
func NewUserAccessHash(s string) (UserAccessHash, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("access hash must be 256 characters, got %d", len(s))
	}

	return UserAccessHash(s), nil
}

// Generate makes request for creating account using provided data and returns Authorization structure.
func (c *Client) Generate(nickname Nickname, description UserDescription, interests []Interest, avatar *FileDescriptor) (*Authorization, error) {
	req := generateRequest{
		Nickname:    nickname,
		Description: description,
		Interests:   interests,
		Avatar:      avatar,
	}

	resp, err := c.do("POST", "/auth/generate", nil, req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("generate failed: status %d", resp.StatusCode)
	}

	var genResp generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &Authorization{
		Id:         genResp.Id,
		AccessHash: genResp.AccessHash,
		Token:      genResp.Token,
	}, nil
}
