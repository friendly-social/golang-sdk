package sdk

import (
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

// NewNickname creates new Nickname or returns an error if length is more than 256.
func NewNickname(s string) (Nickname, error) {
	if len(s) > 256 {
		return "", fmt.Errorf("nickname is too long: %d > 256", len(s))
	}

	return Nickname(s), nil
}

// MustNickname wraps NewNickname and panics on error.
func MustNickname(s string) Nickname {
	n, err := NewNickname(s)
	if err != nil {
		panic(err)
	}

	return n
}

// NewUserDescription creates new UserDescription or returns an error if description is more than 1024.
func NewUserDescription(s string) (UserDescription, error) {
	if len(s) > 1024 {
		return "", fmt.Errorf("description is too long: %d > 1024", len(s))
	}

	return UserDescription(s), nil
}

// MustUserDescription wraps NewUserDescription and panics on error.
func MustUserDescription(s string) UserDescription {
	d, err := NewUserDescription(s)
	if err != nil {
		panic(err)
	}

	return d
}

// NewInterest creates new Interest or returns an error if length is more than 64.
func NewInterest(s string) (Interest, error) {
	if len(s) > 64 {
		return "", fmt.Errorf("interest is too long: %d > 64", len(s))
	}

	return Interest(s), nil
}

// MustInterest wraps NewInterest and panics on error.
func MustInterest(s string) Interest {
	i, err := NewInterest(s)
	if err != nil {
		panic(err)
	}

	return i
}

// NewSocialLink creates new SocialLink or returns an error if length is more than 2048.
func NewSocialLink(s string) (SocialLink, error) {
	if len(s) > 2048 {
		return "", fmt.Errorf("social link is too long: %d > 64", len(s))
	}

	return SocialLink(s), nil
}

// MustSocialLink wraps NewSocialLink and panics on error.
func MustSocialLink(s string) SocialLink {
	l, err := NewSocialLink(s)
	if err != nil {
		panic(err)
	}

	return l
}

// GetSelfDetails returns UserDetails structure for provided Authorization data.
func (c *Client) GetSelfDetails(auth *Authorization) (*UserDetails, error) {
	var details UserDetails
	err := c.do("GET", "/users/details", auth, nil, &details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

// GetUserDetails returns UserDetails for provided user's ID and AccessHash from provided Authorization's perspective.
func (c *Client) GetUserDetails(auth *Authorization, userId UserId, accessHash UserAccessHash) (*UserDetails, error) {
	var details UserDetails
	err := c.do("GET", fmt.Sprintf("/users/details/%d/%s", userId, accessHash), auth, nil, &details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}
