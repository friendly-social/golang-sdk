package sdk

import (
	"encoding/json"
	"fmt"
)

var (
	ErrTooLongNickname        = fmt.Errorf("nickname must be less than 256 characters length")
	ErrTooLongUserDescription = fmt.Errorf("user description must be less than 1024 characters length")
	ErrTooLongInterest        = fmt.Errorf("interest must be less than 64 characters length")
	ErrTooMuchInterests       = fmt.Errorf("maximum amount of interests is 100")
	ErrTooLongSocialLink      = fmt.Errorf("social link must be less than 2048 characters length")

	ErrEmptyUserDescription = fmt.Errorf("user description can't be empty string")
	ErrEmptyNickname        = fmt.Errorf("nickname can't be empty string")
	ErrEmptyInterest        = fmt.Errorf("interest can't be empty string")
	ErrEmptyInterests       = fmt.Errorf("interests can't be empty or nil slice")
	ErrEmptySocialLink      = fmt.Errorf("social link can't be empty string")
)

// --- NICKNAME ---

// Nickname represents user's name (not unique).
type Nickname struct {
	value string
}

// NewNickname creates new Nickname or returns an error if length is more than 256.
func NewNickname(s string) (Nickname, error) {
	if s == "" {
		return Nickname{}, ErrEmptyNickname
	}

	if len(s) > 256 {
		return Nickname{}, fmt.Errorf("length is %d: %w", len(s), ErrTooLongNickname)
	}

	return Nickname{value: s}, nil
}

func (n Nickname) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.value)
}

func (n *Nickname) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &n.value)
}

// --- USER DESCRIPTION ---

// UserDescription represents user's description.
type UserDescription struct {
	value string
}

// NewUserDescription creates new UserDescription or returns an error if description is more than 1024.
func NewUserDescription(s string) (UserDescription, error) {
	if s == "" {
		return UserDescription{}, ErrEmptyUserDescription
	}

	if len(s) > 1024 {
		return UserDescription{}, fmt.Errorf("length is %d: %w", len(s), ErrTooLongUserDescription)
	}

	return UserDescription{value: s}, nil
}

func (d UserDescription) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.value)
}

func (d *UserDescription) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &d.value)
}

// --- INTEREST ---

// Interest represents some user's interest.
type Interest struct {
	value string
}

// NewInterest creates new Interest or returns an error if length is more than 64.
func NewInterest(s string) (Interest, error) {
	if s == "" {
		return Interest{}, ErrEmptyInterest
	}

	if len(s) > 64 {
		return Interest{}, fmt.Errorf("length is %d: %w", len(s), ErrTooLongInterest)
	}

	return Interest{value: s}, nil
}

func (i Interest) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

func (i *Interest) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &i.value)
}

// --- INTERESTS ---

// Interests represents complete list of user's interests.
type Interests struct {
	value []Interest
}

// NewInterests creates new Interests or returns an error if their amount is more than 100.
func NewInterests(i []Interest) (Interests, error) {
	if len(i) == 0 {
		return Interests{}, ErrEmptyInterests
	}

	if len(i) > 100 {
		return Interests{}, fmt.Errorf("amount is %d: %w", len(i), ErrTooMuchInterests)
	}

	return Interests{value: i}, nil
}

func (i Interests) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

func (i *Interests) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &i.value)
}

// --- SOCIAL LINK ---

// SocialLink represents link to user's external social network.
type SocialLink struct {
	value string
}

// NewSocialLink creates new SocialLink or returns an error if length is more than 2048.
func NewSocialLink(s string) (SocialLink, error) {
	if s == "" {
		return SocialLink{}, ErrEmptySocialLink
	}

	if len(s) > 2048 {
		return SocialLink{}, fmt.Errorf("length is %d: %w", len(s), ErrTooLongSocialLink)
	}

	return SocialLink{value: s}, nil
}

func (l SocialLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.value)
}

func (l *SocialLink) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &l.value)
}
