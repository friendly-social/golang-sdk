package sdk

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilesValueTypes(t *testing.T) {
	t.Run("Valid FileAccesssHash", func(t *testing.T) {
		hash, err := NewFileAccessHash(strings.Repeat("1", 256))
		require.EqualValues(t, FileAccessHash(strings.Repeat("1", 256)), hash)
		require.NoError(t, err)
	})

	t.Run("Invalid FileAccessHash", func(t *testing.T) {
		_, err := NewFileAccessHash("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrFileAccessHashLengthMustBe256)
	})
}

func TestGetFileURL(t *testing.T) {
	descriptor := &FileDescriptor{Id: 1, AccessHash: FileAccessHash("hash")}
	client := NewClient("https://getfriend.ly")
	url := client.GetFileURL(descriptor)
	require.Equal(t, "https://getfriend.ly/files/download/1/hash", url)
}
