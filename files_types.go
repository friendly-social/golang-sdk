package sdk

import (
	"encoding/json"
	"fmt"
)

var (
	ErrFileAccessHashLengthMustBe256 = fmt.Errorf("file access hash must be 256 characters length")
)

// --- ID ---

// FileId represents the unique identifier of any file.
type FileId struct {
	value int64
}

// NewFileId creates new FileId from int64.
func NewFileId(i int64) FileId {
	return FileId{value: i}
}

func (i FileId) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

func (i *FileId) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &i.value)
}

// --- HASH ---

// FileAccessHash represents the unique hash associated with file. Works in pair with FileId.
type FileAccessHash struct {
	value string
}

// NewFileAccessHash creates new FileAccessHash or returns an error if hash length isn't 256.
func NewFileAccessHash(s string) (FileAccessHash, error) {
	if len(s) != 256 {
		return FileAccessHash{}, fmt.Errorf("length is %d: %w", len(s), ErrFileAccessHashLengthMustBe256)
	}

	return FileAccessHash{value: s}, nil
}

func (h FileAccessHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.value)
}

func (h *FileAccessHash) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &h.value)
}
