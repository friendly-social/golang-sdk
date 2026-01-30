package sdk

import (
	"context"
	"fmt"
)

// Nickname represents user's name (not unique).
type Nickname string

// UserDescription represents user's description.
type UserDescription string

// Interest represents some user's interest.
type Interest string

// SocialLink represents link to user's external social network.
type SocialLink string

// UserDetails represents complete information about some user: ID, AccessHash, Nickname, Description, list of Interests and Avatar.
type UserDetails struct {
	Id          UserId          `json:"id"`
	AccessHash  UserAccessHash  `json:"accessHash"`
	Nickname    Nickname        `json:"nickname"`
	Description UserDescription `json:"description"`
	Interests   []Interest      `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
}

var (
	ErrNicknameLengthMustBeLessThan256         = fmt.Errorf("nickname must be less than 256 characters lenght")
	ErrUserDescriptionLengthMustBeLessThan1024 = fmt.Errorf("user description must be less than 1024 characters lenght")
	ErrInterestLengthMustBeLessThan64          = fmt.Errorf("interest must be less than 64 characters lenght")
	ErrSocialLinkLengthMustBeLessThan2048      = fmt.Errorf("social link must be less than 2048 characters lenght")
)

// NewNickname creates new Nickname or returns an error if length is more than 256.
func NewNickname(s string) (Nickname, error) {
	if len(s) > 256 {
		return "", fmt.Errorf("length is %d: %w", len(s), ErrNicknameLengthMustBeLessThan256)
	}

	return Nickname(s), nil
}

// NewUserDescription creates new UserDescription or returns an error if description is more than 1024.
func NewUserDescription(s string) (UserDescription, error) {
	if len(s) > 1024 {
		return "", fmt.Errorf("length is %d: %w", len(s), ErrUserDescriptionLengthMustBeLessThan1024)
	}

	return UserDescription(s), nil
}

// NewInterest creates new Interest or returns an error if length is more than 64.
func NewInterest(s string) (Interest, error) {
	if len(s) > 64 {
		return "", fmt.Errorf("length is %d: %w", len(s), ErrInterestLengthMustBeLessThan64)
	}

	return Interest(s), nil
}

// NewSocialLink creates new SocialLink or returns an error if length is more than 2048.
func NewSocialLink(s string) (SocialLink, error) {
	if len(s) > 2048 {
		return "", fmt.Errorf("length is %d: %w", len(s), ErrSocialLinkLengthMustBeLessThan2048)
	}

	return SocialLink(s), nil
}

// GetSelfDetails returns UserDetails structure for provided Authorization data.
func (c *Client) GetSelfDetails(ctx context.Context, auth *Authorization) (*UserDetails, error) {
	var details UserDetails
	err := c.do(ctx, "GET", "/users/details", auth, nil, &details)
	if err != nil {
		return nil, fmt.Errorf("failed to get self details: %w", err)
	}

	return &details, nil
}

// GetUserDetails returns UserDetails for provided user's ID and AccessHash from provided Authorization's perspective.
func (c *Client) GetUserDetails(ctx context.Context, auth *Authorization, userId UserId, accessHash UserAccessHash) (*UserDetails, error) {
	var details UserDetails
	err := c.do(ctx, "GET", fmt.Sprintf("/users/details/%d/%s", userId, accessHash), auth, nil, &details)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	return &details, nil
}
