package sdk

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
)

// FileId represents the unique identifier of any file.
type FileId int64

// FileAccessHash represents the unique hash associated with file. Works in pair with FileId.
type FileAccessHash string

// FileDescriptor is a helper structure for composing file's ID and AccessHash.
type FileDescriptor struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

type uploadFileResponse struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

var (
	ErrFileAccessHashMustBe256CharactersLength = fmt.Errorf("file access hash must be 256 characters lenght")
)

// NewFileAccessHash creates new FileAccessHash or returns an error if hash length isn't 256.
func NewFileAccessHash(s string) (FileAccessHash, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("length is %d: %w", len(s), ErrFileAccessHashMustBe256CharactersLength)
	}

	return FileAccessHash(s), nil
}

// GetFileURL returns file access URL for corresponding descriptor.
func (c *Client) GetFileURL(descriptor *FileDescriptor) string {
	return fmt.Sprintf("%s/files/download/%d/%s", c.url, descriptor.Id, descriptor.AccessHash)
}

// UploadFile uploads file from disk to the server and returns corresponding descriptor.
func (c *Client) UploadFile(filename string, reader io.Reader) (*FileDescriptor, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	filename = path.Base(filename)

	go func() {
		var err error
		defer func() {
			if cerr := writer.Close(); cerr != nil && err == nil {
				err = cerr
			}
			_ = pw.CloseWithError(err)
		}()

		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			return
		}

		_, err = io.Copy(part, reader)
	}()

	completePath, err := url.JoinPath(c.url, "/files/upload")
	if err != nil {
		return nil, fmt.Errorf("%w: %s + %s", ErrInvalidPath, c.url, "/files/upload")
	}

	req, err := http.NewRequest("POST", completePath, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var resp uploadFileResponse
	err = c.execute(req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return &FileDescriptor{
		Id:         resp.Id,
		AccessHash: resp.AccessHash,
	}, nil
}
