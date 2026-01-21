package sdk

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/agiledragon/gomonkey"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestFilesValueTypes(t *testing.T) {
	RunValueTypeTest(t, ValueTypeTestCase[string, FileAccessHash]{
		name:    "FileAccessHash",
		new:     NewFileAccessHash,
		valid:   strings.Repeat("1", 256),
		invalid: "1",
	})
}

func TestGetFileUrl(t *testing.T) {
	client := NewClient("https://getfriend.ly")
	file := &FileDescriptor{Id: 1, AccessHash: "222"}
	require.Equal(t, client.GetFileURL(file), "https://getfriend.ly/files/download/1/222")
}

func TestUploadFile(t *testing.T) {
	f := func() *gock.Request {
		// matcher that ensures that we read an entire body
		m := gock.NewMatcher()

		m.Add(func(req *http.Request, _ *gock.Request) (bool, error) {
			_, err := io.ReadAll(req.Body)
			return err == nil, err
		})

		return gock.New("https://getfriend.ly").Post("/files/upload").SetMatcher(m)
	}

	call := func(c *Client) (*FileDescriptor, error) {
		return c.UploadFile("test.png", strings.NewReader("data"))
	}

	RunAPITests(t, append(CommonCasesTests(f, call),
		APITestCase[FileDescriptor]{
			name:             "Success",
			setup:            func() { f().Reply(200).JSON(uploadFileResponse{Id: 1, AccessHash: "1"}) },
			call:             call,
			expectError:      false,
			expectedResponse: &FileDescriptor{Id: 1, AccessHash: "1"},
		},
		APITestCase[FileDescriptor]{
			name:  "Interrupted",
			setup: func() { f().Reply(200).JSON(uploadFileResponse{Id: 1, AccessHash: "1"}) },
			call: func(c *Client) (*FileDescriptor, error) {
				return c.UploadFile("test.png", io.MultiReader(
					strings.NewReader("initial data"),
					iotest.ErrReader(errors.New("disk read error")),
				))
			},
			expectError:      true,
			expectedResponse: nil,
		},
		APITestCase[FileDescriptor]{
			name:  "Invalid URL",
			setup: func() { f().Reply(200) },
			call: func(c *Client) (*FileDescriptor, error) {
				return NewClient("://").UploadFile("test.png", strings.NewReader("initial data"))
			},
			expectError:      true,
			expectedResponse: nil,
		},
		APITestCase[FileDescriptor]{
			name:  "Invalid Form File",
			setup: func() { f().Reply(200).JSON(uploadFileResponse{Id: 1, AccessHash: "1"}) },
			call: func(c *Client) (*FileDescriptor, error) {
				patches := gomonkey.ApplyMethod(reflect.TypeFor[*multipart.Writer](), "CreateFormFile",
					func(_ *multipart.Writer, fieldname, filename string) (io.Writer, error) {
						return nil, errors.New("failure")
					})
				resp, err := c.UploadFile("test.png", strings.NewReader("data"))
				patches.Reset()
				return resp, err
			},
			expectError:      true,
			expectedResponse: nil,
		},
	))
}
