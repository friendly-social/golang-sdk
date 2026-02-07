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

// Interests represents complete list of user's interests.
type Interests []Interest

// SocialLink represents link to user's external social network.
type SocialLink string

// UserDetails represents complete information about some user: ID, AccessHash, Nickname, Description, list of Interests and Avatar.
type UserDetails struct {
	Id          UserId          `json:"id"`
	AccessHash  UserAccessHash  `json:"accessHash"`
	Nickname    Nickname        `json:"nickname"`
	Description UserDescription `json:"description"`
	Interests   Interests       `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
}

type editAccountValue[T any] struct {
	Value T `json:"value"`
}

type editAccountRequest struct {
	Nickname    *editAccountValue[Nickname]        `json:"nickname,omitempty"`
	Description *editAccountValue[UserDescription] `json:"description,omitempty"`
	Interests   *editAccountValue[Interests]       `json:"interests,omitempty"`
	Avatar      *editAccountValue[*FileDescriptor] `json:"avatar,omitempty"`
	SocialLink  *editAccountValue[SocialLink]      `json:"socialLink,omitempty"`
}

type editAccountOption func(*editAccountRequest)

var (
	ErrNicknameLengthMustBeLessThan256         = fmt.Errorf("nickname must be less than 256 characters length")
	ErrUserDescriptionLengthMustBeLessThan1024 = fmt.Errorf("user description must be less than 1024 characters lengt")
	ErrInterestLengthMustBeLessThan64          = fmt.Errorf("interest must be less than 64 characters length")
	ErrSocialLinkLengthMustBeLessThan2048      = fmt.Errorf("social link must be less than 2048 characters length")
	ErrTooMuchInterests                        = fmt.Errorf("maximum amount of interests is 100")
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

// NewInterests creates new Interests or returns an error if thier amount is more than 100.
func NewInterests(interests ...Interest) (Interests, error) {
	if len(interests) > 100 {
		return nil, fmt.Errorf("amount is %d: %w", len(interests), ErrTooMuchInterests)
	}

	return Interests(interests), nil
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
	err := c.do(ctx, auth, "GET", "/users/details", nil, &details)
	if err != nil {
		return nil, fmt.Errorf("failed to get self details: %w", err)
	}

	return &details, nil
}

// GetUserDetails returns UserDetails for provided user's ID and AccessHash from provided Authorization's perspective.
func (c *Client) GetUserDetails(ctx context.Context, auth *Authorization, userId UserId, accessHash UserAccessHash) (*UserDetails, error) {
	var details UserDetails
	err := c.do(ctx, auth, "GET", fmt.Sprintf("/users/details/%d/%s", userId, accessHash), nil, &details)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	return &details, nil
}

// WithUserNickname applies new Nickname for editing account request.
func WithUserNickname(nickname Nickname) editAccountOption {
	return func(r *editAccountRequest) {
		if nickname != "" {
			r.Nickname = &editAccountValue[Nickname]{nickname}
		}
	}
}

// WithUserDescription applies new UserDescription for editing account request.
func WithUserDescription(description UserDescription) editAccountOption {
	return func(r *editAccountRequest) {
		if description != "" {
			r.Description = &editAccountValue[UserDescription]{description}
		}
	}
}

// WithUserInterests applies new Interests for editing account request.
func WithUserInterests(interests Interests) editAccountOption {
	return func(r *editAccountRequest) {
		if len(interests) != 0 {
			r.Interests = &editAccountValue[Interests]{interests}
		}
	}
}

// WithUserAvatar applies new Avatar for editing account request.
func WithUserAvatar(avatar *FileDescriptor) editAccountOption {
	return func(r *editAccountRequest) {
		if avatar != nil {
			r.Avatar = &editAccountValue[*FileDescriptor]{avatar}
		}
	}
}

// WithUserSocialLink applies new SocialLink for editing account request.
func WithUserSocialLink(link SocialLink) editAccountOption {
	return func(r *editAccountRequest) {
		if link != "" {
			r.SocialLink = &editAccountValue[SocialLink]{link}
		}
	}
}

// EditAccount makes request for editing account using provided options.
func (c *Client) EditAccount(ctx context.Context, opts ...editAccountOption) error {
	if len(opts) == 0 {
		return nil
	}

	req := editAccountRequest{}
	for _, opt := range opts {
		opt(&req)
	}

	err := c.do(ctx, nil, "PATCH", "/users/edit", req, nil)
	if err != nil {
		return fmt.Errorf("failed to edit account: %w", err)
	}

	return nil
}
