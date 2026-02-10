package sdk

import (
	"context"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestEmailValueTypes(t *testing.T) {
	t.Run("Valid Email", func(t *testing.T) {
		email, err := NewEmail("example@example.com")
		require.EqualValues(t, "example@example.com", email)
		require.NoError(t, err)
	})

	t.Run("Invalid Email", func(t *testing.T) {
		_, err := NewEmail("junk")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidEmail)
	})

	t.Run("Too long Email", func(t *testing.T) {
		_, err := NewEmail(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongEmail)
	})

	t.Run("Valid EmailCode", func(t *testing.T) {
		email, err := NewEmailCode("11111111")
		require.EqualValues(t, "11111111", email)
		require.NoError(t, err)
	})

	t.Run("Invalid EmailCode", func(t *testing.T) {
		_, err := NewEmailCode("junk")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEmailCodeLengthMustBe8)
	})
}

func TestLinkEmail_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/link").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"email":"example@example.com"}`).
		Reply(200)

	c := NewClient("https://api.getfriend.ly")
	err := c.LinkEmail(context.Background(), &Authorization{Id: 1, Token: "token"}, "example@example.com")
	require.NoError(t, err)
}

func TestLinkEmail_Taken(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").Post("/email/link").
		Reply(409)

	c := NewClient("https://api.getfriend.ly")
	err := c.LinkEmail(context.Background(), &Authorization{Id: 1, Token: "token"}, "example@example.com")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEmailTaken)
}

func TestLinkEmail_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/link").
		Reply(400)

	c := NewClient("https://api.getfriend.ly")
	err := c.LinkEmail(context.Background(), &Authorization{Id: 1, Token: "token"}, "example@example.com")
	require.Error(t, err)
}

func TestConfirmEmail_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/confirm").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"code":"11111111"}`).
		Reply(200)

	c := NewClient("https://api.getfriend.ly")
	err := c.ConfirmEmail(context.Background(), &Authorization{Id: 1, Token: "token"}, EmailCode("11111111"))
	require.NoError(t, err)
}

func TestConfirmEmail_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/confirm").
		Reply(400)

	c := NewClient("https://api.getfriend.ly")
	err := c.ConfirmEmail(context.Background(), &Authorization{Id: 1, Token: "token"}, EmailCode("11111111"))
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

	c := NewClient("https://api.getfriend.ly")
	err := c.UnlinkEmail(context.Background(), &Authorization{Id: 1, Token: "token"})
	require.NoError(t, err)
}

func TestUnlinkEmail_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/email/unlink").
		Reply(400)

	c := NewClient("https://api.getfriend.ly")
	err := c.UnlinkEmail(context.Background(), &Authorization{Id: 1, Token: "token"})
	require.Error(t, err)
}
