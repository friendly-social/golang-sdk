package sdk

import (
	"encoding/json"
	"fmt"
)

var (
	ErrFriendTokenLengthMustBe256 = fmt.Errorf("friend token must be 256 characters length")
	ErrFriendTokenExpired         = fmt.Errorf("friend token expired")
)

// --- FRIEND TOKEN

// FriendToken is a token by which other users can add Token's owner to their friend list.
type FriendToken struct {
	value string
}

// Valu returns FriendToken as a plain string.
func (t FriendToken) Value() string {
	return t.value
}

// NewFriendToken creates new FriendToken or returns an error if length is not 256.
func NewFriendToken(s string) (FriendToken, error) {
	if len(s) != 256 {
		return FriendToken{}, fmt.Errorf("length is %d: %w", len(s), ErrFriendTokenLengthMustBe256)
	}

	return FriendToken{value: s}, nil
}

func (h FriendToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.value)
}

func (t *FriendToken) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &t.value)
}
