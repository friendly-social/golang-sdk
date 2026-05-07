package sdk

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestLinkEmail_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/link").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		MatchHeader("X-Locale", "ru").
		JSON(`{"email":"example@example.com"}`).
		Reply(200)

	c := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := c.LinkEmail(context.Background(), auth, MockEmail("example@example.com"), MockEmailLocale("ru"))
	require.NoError(t, err)
}

func TestLinkEmail_Taken(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/link").
		Reply(409)

	c := NewClient()
	err := c.LinkEmail(context.Background(), nil, MockEmail("example@example.com"), MockEmailLocale("ru"))
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEmailTaken)
}

func TestLinkEmail_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/link").
		Reply(400)

	c := NewClient()
	err := c.LinkEmail(context.Background(), nil, MockEmail("example@example.com"), MockEmailLocale("ru"))
	require.Error(t, err)
}

func TestLinkEmail_NewRequestFailed(t *testing.T) {
	c := NewClient()
	err := c.LinkEmail(nil, nil, MockEmail("example@example.com"), MockEmailLocale("ru")) //nolint:staticcheck
	require.Error(t, err)
}

func TestConfirmEmail_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/confirm").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"code":11111111}`).
		Reply(200)

	c := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := c.ConfirmEmail(context.Background(), auth, MockEmailCode(11111111))
	require.NoError(t, err)
}

func TestConfirmEmail_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/confirm").
		Reply(400)

	c := NewClient()
	err := c.ConfirmEmail(context.Background(), nil, MockEmailCode(11111111))
	require.Error(t, err)
}

func TestUnlinkEmail_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/unlink").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200)

	c := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := c.UnlinkEmail(context.Background(), auth)
	require.NoError(t, err)
}

func TestUnlinkEmail_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/unlink").
		Reply(400)

	c := NewClient()
	err := c.UnlinkEmail(context.Background(), nil)
	require.Error(t, err)
}
