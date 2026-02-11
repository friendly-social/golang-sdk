package sdk

import (
	"context"
	"fmt"
	"strings"
)

type linkEmailRequest struct {
	Email Email `json:"email"`
}

type confirmEmailRequest struct {
	Code EmailCode `json:"code"`
}

var (
	ErrEmailTaken = fmt.Errorf("invalid e-mail: already used by another user")
)

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
