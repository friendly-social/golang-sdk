package sdk

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockFileId(i int64) FileId {
	return FileId{value: i}
}

func MockFileAccessHash(s string) FileAccessHash {
	return FileAccessHash{value: s}
}

func TestFileTypes(t *testing.T) {
	t.Run("FileId", func(t *testing.T) {
		id := NewFileId(123)
		require.Equal(t, FileId{value: 123}, id)

		data, err := json.Marshal(id)
		require.NoError(t, err)
		require.Equal(t, `123`, string(data))

		var loadedId FileId
		err = json.Unmarshal(data, &loadedId)
		require.NoError(t, err)
		require.Equal(t, id, loadedId)
	})

	t.Run("FileAccessHash", func(t *testing.T) {
		hash, err := NewFileAccessHash(strings.Repeat("1", 256))
		require.NoError(t, err)
		require.Equal(t, MockFileAccessHash(strings.Repeat("1", 256)), hash)

		_, err = NewFileAccessHash("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrFileAccessHashLengthMustBe256)

		data, err := json.Marshal(hash)
		require.NoError(t, err)
		require.Equal(t, `"`+strings.Repeat("1", 256)+`"`, string(data))

		var loadedHash FileAccessHash
		err = json.Unmarshal(data, &loadedHash)
		require.NoError(t, err)
		require.Equal(t, hash, loadedHash)
	})
}
