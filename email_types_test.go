package sdk

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockEmail(s string) Email {
	return Email{value: s}
}

func MockEmailCode(i int64) EmailCode {
	return EmailCode{value: i}
}

func MockEmailLocale(s string) EmailLocale {
	return EmailLocale{value: s}
}

func TestEmailTypes(t *testing.T) {
	t.Run("Email", func(t *testing.T) {
		email, err := NewEmail("example@example.com")
		require.Equal(t, MockEmail("example@example.com"), email)
		require.NoError(t, err)

		_, err = NewEmail("junk")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidEmail)

		_, err = NewEmail(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongEmail)

		data, err := json.Marshal(email)
		require.NoError(t, err)
		require.Equal(t, `"example@example.com"`, string(data))

		var loadedEmail Email
		err = json.Unmarshal(data, &loadedEmail)
		require.NoError(t, err)
		require.Equal(t, email, loadedEmail)
	})

	t.Run("EmailCode", func(t *testing.T) {
		emailCode, err := NewEmailCode("11111111")
		require.Equal(t, MockEmailCode(11111111), emailCode)
		require.NoError(t, err)

		_, err = NewEmailCode("junk")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEmailCodeLengthMustBe8)

		_, err = NewEmailCode("junkjunk")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidEmailCode)

		data, err := json.Marshal(emailCode)
		require.NoError(t, err)
		require.Equal(t, `11111111`, string(data))

		var loadedEmailCode EmailCode
		err = json.Unmarshal(data, &loadedEmailCode)
		require.NoError(t, err)
		require.Equal(t, emailCode, loadedEmailCode)
	})

	t.Run("EmailLocale", func(t *testing.T) {
		locale := NewEmailLocale("ru")
		require.Equal(t, MockEmailLocale("ru"), locale)

		data, err := json.Marshal(locale)
		require.NoError(t, err)
		require.Equal(t, `"ru"`, string(data))

		var loadedEmail EmailLocale
		err = json.Unmarshal(data, &loadedEmail)
		require.NoError(t, err)
		require.Equal(t, locale, loadedEmail)
	})
}
