package sdk

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

var emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

var (
	ErrInvalidEmail           = fmt.Errorf("invalid e-mail format: must be example@example.com")
	ErrTooLongEmail           = fmt.Errorf("invalid e-mail: length must be less than 2048")
	ErrEmailCodeLengthMustBe8 = fmt.Errorf("email verification code must be exactly 8 digits length")
	ErrInvalidEmailCode       = fmt.Errorf("invalid e-mail code: must be in format 11111111")
)

// --- EMAIL ---

// Email represents e-mail linked to account.
type Email struct {
	value string
}

// NewEmail creates new Email object ensuring valid format and length.
func NewEmail(s string) (Email, error) {
	if len(s) > 2048 {
		return Email{}, fmt.Errorf("failed to create e-mail: %w", ErrTooLongEmail)
	}

	if !emailRegexp.MatchString(s) {
		return Email{}, fmt.Errorf("failed to create e-mail: %w", ErrInvalidEmail)
	}

	return Email{value: s}, nil
}

func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.value)
}

func (e *Email) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &e.value)
}

// --- EMAIL CODE ---

// EmailCode represents e-mail verification code.
type EmailCode struct {
	value int64
}

// NewEmailCode creates new EmailCode object ensuring valid length.
func NewEmailCode(s string) (EmailCode, error) {
	if len(s) != 8 {
		return EmailCode{}, fmt.Errorf("failed to create e-mail code %s: %w", s, ErrEmailCodeLengthMustBe8)
	}

	code, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return EmailCode{}, fmt.Errorf("failed to create e-mail code %s: %w", s, ErrInvalidEmailCode)
	}

	return EmailCode{value: code}, nil
}

func (e EmailCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.value)
}

func (e *EmailCode) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &e.value)
}
