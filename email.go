package sdk

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Email represents e-mail linked to account.
type Email string

// EmailCode represents e-mail verification code.
type EmailCode string

type linkEmailRequest struct {
	Email Email `json:"email"`
}

type confirmEmailRequest struct {
	Code EmailCode `json:"code"`
}

var (
	emailRegexp               = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	ErrInvalidEmail           = fmt.Errorf("invalid e-mail format: must be example@example.com")
	ErrTooLongEmail           = fmt.Errorf("invalid e-mail: length must be less than 2048")
	ErrEmailTaken             = fmt.Errorf("invalid e-mail: already used by another user")
	ErrEmailCodeLengthMustBe8 = fmt.Errorf("email verification code must be exactly 8 digits length")
)

// NewEmail creates new Email object ensuring valid format and length.
func NewEmail(s string) (Email, error) {
	if len(s) > 2048 {
		return "", fmt.Errorf("failed to create e-mail: %w", ErrTooLongEmail)
	}

	if !emailRegexp.MatchString(s) {
		return "", fmt.Errorf("failed to create e-mail: %w", ErrInvalidEmail)
	}

	return Email(s), nil
}

// NewEmailCode creates new EmailCode object ensuring valid length.
func NewEmailCode(s string) (EmailCode, error) {
	if len(s) != 8 {
		return "", fmt.Errorf("failed to create e-mail code: %w", ErrEmailCodeLengthMustBe8)
	}

	return EmailCode(s), nil
}

// LinkEmail sends request for linking unverified e-mail to the Authorization.
func (c *Client) LinkEmail(ctx context.Context, auth *Authorization, email Email) error {
	req := linkEmailRequest{
		Email: email,
	}

	err := c.do(ctx, auth, "POST", "/email/link", req, nil)
	if err != nil {
		if strings.Contains(err.Error(), "status code 409") {
			return fmt.Errorf("%w: %w", ErrEmailTaken, err)
		}

		return err
	}

	return nil
}

// ConfirmEmail sends request for confirming e-mail for Authorization using provided EmailCode.
func (c *Client) ConfirmEmail(ctx context.Context, auth *Authorization, code EmailCode) error {
	req := confirmEmailRequest{
		Code: code,
	}

	err := c.do(ctx, auth, "POST", "/email/confirm", req, nil)
	if err != nil {
		return err
	}

	return nil
}

// UnlinkEmail sends request for unlinking currently linking e-mail of Authorization.
func (c *Client) UnlinkEmail(ctx context.Context, auth *Authorization) error {
	err := c.do(ctx, auth, "POST", "/email/unlink", nil, nil)
	if err != nil {
		return err
	}

	return nil
}
