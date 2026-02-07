package sdk

import (
	"context"
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
	ErrFileAccessHashLengthMustBe256 = fmt.Errorf("file access hash must be 256 characters length")
)

// NewFileAccessHash creates new FileAccessHash or returns an error if hash length isn't 256.
func NewFileAccessHash(s string) (FileAccessHash, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("length is %d: %w", len(s), ErrFileAccessHashLengthMustBe256)
	}

	return FileAccessHash(s), nil
}

// GetFileURL returns file access URL for corresponding descriptor.
func (c *Client) GetFileURL(fd *FileDescriptor) (string, error) {
	url, err := url.JoinPath(c.url, fmt.Sprintf("/files/download/%d/%s", fd.Id, fd.AccessHash))
	if err != nil {
		return "", fmt.Errorf("invalid path: %s + %s", c.url, fmt.Sprintf("/files/download/%d/%s", fd.Id, fd.AccessHash))
	}

	return url, nil
}

// DownloadFile opens connection for downloading file and returns corresponding io.ReadCloser.
func (c *Client) DownloadFile(ctx context.Context, fd *FileDescriptor) (io.ReadCloser, error) {
	url, err := c.GetFileURL(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to get file URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp.Body, nil
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("failed to download file: %w", ErrNotFound)
	default:
		return nil, fmt.Errorf("failed to download file: unexpected status code %d", resp.StatusCode)
	}
}

// UploadFile uploads file from io.Reader to the server and returns corresponding descriptor.
// It accepts real IP address of client machine (required by Cloufdlare), filename by which file will be saved on server, and reader from which file will be read.
func (c *Client) UploadFile(ctx context.Context, ip string, filename string, reader io.Reader) (*FileDescriptor, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	filename = path.Base(filename)

	go func() {
		var err error
		defer func() {
			_ = writer.Close()
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
		return nil, fmt.Errorf("invalid path: %s + %s", c.url, "/files/upload")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", completePath, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("CF-Connecting-IP", ip)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var resp uploadFileResponse
	err = c.execute(req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &FileDescriptor{
		Id:         resp.Id,
		AccessHash: resp.AccessHash,
	}, nil
}
