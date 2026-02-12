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

// FileDescriptor is a helper structure for composing file's ID and AccessHash.
type FileDescriptor struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

type uploadFileResponse struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

// GetFileURL returns file access URL for corresponding descriptor.
func (c *Client) GetFileURL(fd *FileDescriptor) (string, error) {
	url, err := url.JoinPath(c.url, fmt.Sprintf("/files/download/%d/%s", fd.Id.value, fd.AccessHash.value))
	if err != nil {
		return "", fmt.Errorf("invalid path: %s + %s", c.url, fmt.Sprintf("/files/download/%d/%s", fd.Id.value, fd.AccessHash.value))
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
func (c *Client) UploadFile(ctx context.Context, filename string, reader io.Reader) (*FileDescriptor, error) {
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
