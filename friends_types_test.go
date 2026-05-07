package sdk

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockFriendToken(s string) FriendToken {
	return FriendToken{value: s}
}

func TestFriendsTypes(t *testing.T) {
	t.Run("FriendToken", func(t *testing.T) {
		token, err := NewFriendToken(strings.Repeat("1", 256))
		require.NoError(t, err)
		require.Equal(t, MockFriendToken(strings.Repeat("1", 256)), token)
		require.Equal(t, strings.Repeat("1", 256), token.Value())

		_, err = NewFriendToken("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrFriendTokenLengthMustBe256)

		data, err := json.Marshal(token)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(`"%s"`, strings.Repeat("1", 256)), string(data))

		var loadedToken FriendToken
		err = json.Unmarshal(data, &loadedToken)
		require.NoError(t, err)
		require.Equal(t, token, loadedToken)
	})
}
