package sdk

import (
	"context"
	"fmt"
)

// Authorization is a helper structure for composing user's ID, AccessHash and Token for authorization.
type Authorization struct {
	Id         UserId         `json:"id"`
	AccessHash UserAccessHash `json:"accessHash"`
	Token      Token          `json:"token"`
}

type registerRequest struct {
	Nickname    Nickname        `json:"nickname"`
	Description UserDescription `json:"description"`
	Interests   Interests       `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
	SocialLink  SocialLink      `json:"socialLink"`
}

type registerResponse struct {
	Id         UserId         `json:"id"`
	AccessHash UserAccessHash `json:"accessHash"`
	Token      Token          `json:"token"`
}

type sendLoginRequest struct {
	Email Email `json:"email"`
}

type confirmLoginRequest struct {
	Email Email     `json:"email"`
	Code  EmailCode `json:"code"`
}

// Register makes request for creating account using provided data and returns Authorization structure.
func (c *Client) Register(ctx context.Context, nickname Nickname, description UserDescription, interests Interests, avatar *FileDescriptor, link SocialLink) (*Authorization, error) {
	req := registerRequest{
		Nickname:    nickname,
		Description: description,
		Interests:   interests,
		Avatar:      avatar,
		SocialLink:  link,
	}

	var resp registerResponse
	err := c.do(ctx, nil, "POST", "/auth/generate", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return &Authorization{
		Id:         resp.Id,
		AccessHash: resp.AccessHash,
		Token:      resp.Token,
	}, nil
}

// SendLoginRequest sends login code to provided e-mail.
func (c *Client) SendLoginRequest(ctx context.Context, email Email) error {
	req := sendLoginRequest{
		Email: email,
	}

	err := c.do(ctx, nil, "POST", "/auth/email", req, nil)
	if err != nil {
		return fmt.Errorf("failed to send login request: %w", err)
	}

	return nil
}

// ConfirmLogin returns new Authorization for linked user if provided code is correct.
func (c *Client) ConfirmLogin(ctx context.Context, email Email, code EmailCode) (*Authorization, error) {
	req := confirmLoginRequest{
		Email: email,
		Code:  code,
	}

	var resp registerResponse
	err := c.do(ctx, nil, "POST", "/auth/login", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return &Authorization{
		Id:         resp.Id,
		AccessHash: resp.AccessHash,
		Token:      resp.Token,
	}, nil
}
