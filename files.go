package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

// NewFileAccessHash creates new FileAccessHash or returns an error if hash length isn't 256.
func NewFileAccessHash(s string) (FileAccessHash, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("file access hash must be 256 characters, got %d", len(s))
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

	go func() {
		var err error
		defer func() {
			writer.Close() //nolint:errcheck
			pw.CloseWithError(err)
		}()

		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			return
		}

		_, err = io.Copy(part, reader)
	}()

	req, err := http.NewRequest("POST", c.url+"/files/upload", pr)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed: status %d", resp.StatusCode)
	}

	var uploadResp uploadFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, err
	}

	return &FileDescriptor{
		Id:         uploadResp.Id,
		AccessHash: uploadResp.AccessHash,
	}, nil
}
