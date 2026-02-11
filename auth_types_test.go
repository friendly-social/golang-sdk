package sdk

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockUserId(i int64) UserId {
	return UserId{value: i}
}

func MockUserAccessHash(s string) UserAccessHash {
	return UserAccessHash{value: s}
}

func MockToken(s string) Token {
	return Token{value: s}
}

func TestAuthTypes(t *testing.T) {
	t.Run("UserId", func(t *testing.T) {
		id := NewUserId(123)
		require.Equal(t, UserId{value: 123}, id)

		data, err := json.Marshal(id)
		require.NoError(t, err)
		require.Equal(t, `123`, string(data))

		var loadedId UserId
		err = json.Unmarshal(data, &loadedId)
		require.NoError(t, err)
		require.Equal(t, id, loadedId)
	})

	t.Run("Token", func(t *testing.T) {
		token, err := NewToken(strings.Repeat("1", 256))
		require.Equal(t, MockToken(strings.Repeat("1", 256)), token)
		require.NoError(t, err)

		_, err = NewToken("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTokenLengthMustBe256)

		data, err := json.Marshal(token)
		require.NoError(t, err)
		require.Equal(t, `"`+strings.Repeat("1", 256)+`"`, string(data))

		var loadedToken Token
		err = json.Unmarshal(data, &loadedToken)
		require.NoError(t, err)
		require.Equal(t, token, loadedToken)
	})

	t.Run("UserAccesssHash", func(t *testing.T) {
		hash, err := NewUserAccessHash(strings.Repeat("1", 256))
		require.Equal(t, MockUserAccessHash(strings.Repeat("1", 256)), hash)
		require.NoError(t, err)

		_, err = NewUserAccessHash("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserAccessHashLengthMustBe256)

		data, err := json.Marshal(hash)
		require.NoError(t, err)
		require.Equal(t, `"`+strings.Repeat("1", 256)+`"`, string(data))

		var loadedHash UserAccessHash
		err = json.Unmarshal(data, &loadedHash)
		require.NoError(t, err)
		require.Equal(t, hash, loadedHash)
	})
}
