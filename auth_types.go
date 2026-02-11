package sdk

import (
	"encoding/json"
	"fmt"
)

var (
	ErrTokenLengthMustBe256          = fmt.Errorf("token must be 256 characters long")
	ErrUserAccessHashLengthMustBe256 = fmt.Errorf("user access hash must be 256 characters long")
)

// --- USER ID ---

// UserId represents the unique identifier of user.
type UserId struct {
	value int64
}

// NewUserId creates new UserId from int64.
func NewUserId(i int64) UserId {
	return UserId{value: i}
}

func (i UserId) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

func (i *UserId) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &i.value)
}

// --- ACCESS HASH ---

// UserAccessHash represents the unique hash associated with user. Works in trio with UserId and Token.
type UserAccessHash struct {
	value string
}

// NewUserAccessHash creates new UserAccessHash or returns an error if hash length isn't 256.
func NewUserAccessHash(s string) (UserAccessHash, error) {
	if len(s) != 256 {
		return UserAccessHash{}, fmt.Errorf("length is %d: %w", len(s), ErrUserAccessHashLengthMustBe256)
	}

	return UserAccessHash{value: s}, nil
}

func (h UserAccessHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.value)
}

func (h *UserAccessHash) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &h.value)
}

// --- TOKEN ---

// Token represents access token for the user. Works in trio with UserId and UserAccessHash.
type Token struct {
	value string
}

// NewToken creates new Token or returns an error if token's length isn't 256.
func NewToken(s string) (Token, error) {
	if len(s) != 256 {
		return Token{}, fmt.Errorf("length is %d: %w", len(s), ErrTokenLengthMustBe256)
	}

	return Token{value: s}, nil
}

func (t Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func (t *Token) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &t.value)
}
