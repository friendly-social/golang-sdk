package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Nickname represents user's name (not unique).
type Nickname string

// UserDescription represents user's description.
type UserDescription string

// Interest represents some user's interest.
type Interest string

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

// GetSelfDetails returns UserDetails structure for provided Authorization data.
func (c *Client) GetSelfDetails(auth *Authorization) (*UserDetails, error) {
	resp, err := c.do("GET", "/users/details", auth, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get details failed: status %d", resp.StatusCode)
	}

	var details UserDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}

// GetUserDetails returns UserDetails for provided user's ID and AccessHash from provided Authorization's perspective.
func (c *Client) GetUserDetails(auth *Authorization, userId UserId, accessHash UserAccessHash) (*UserDetails, error) {
	path := fmt.Sprintf("/users/details/%d/%s", userId, accessHash)
	resp, err := c.do("GET", path, auth, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user details failed: status %d", resp.StatusCode)
	}

	var details UserDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}
