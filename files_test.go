package sdk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
	"github.com/h2non/gock"
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

func TestGetFileURL_Success(t *testing.T) {
	descriptor := &FileDescriptor{Id: 1, AccessHash: FileAccessHash("hash")}
	client := NewClient("https://getfriend.ly")

	url, err := client.GetFileURL(descriptor)
	require.NoError(t, err)
	require.Equal(t, "https://getfriend.ly/files/download/1/hash", url)
}

func TestGetFileURL_InvalidURL(t *testing.T) {
	descriptor := &FileDescriptor{Id: 1, AccessHash: FileAccessHash("hash")}
	client := NewClient("::invalid")

	_, err := client.GetFileURL(descriptor)
	require.Error(t, err)
}

func TestUploadFile_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/files/upload", r.URL.Path)
		require.Equal(t, "1.2.3.4", r.Header["Cf-Connecting-Ip"][0])

		err := r.ParseMultipartForm(1 << 20)
		require.NoError(t, err)
		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close() //nolint:errcheck

		require.Equal(t, "file.txt", header.Filename)
		data, err := io.ReadAll(file)
		require.NoError(t, err)
		require.Equal(t, "hello world", string(data))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"accessHash": "hash"
		}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	file := strings.NewReader("hello world")

	fd, err := client.UploadFile(context.Background(), "1.2.3.4", "file.txt", file)
	require.NoError(t, err)
	require.Equal(t, &FileDescriptor{Id: 123, AccessHash: "hash"}, fd)
}

func TestUploadFile_Canceled(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Post("/files/upload").
		Reply(200).
		Delay(100 * time.Hour)

	client := NewClient("https://getfriend.ly")
	file := strings.NewReader("hello world")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.UploadFile(ctx, "1.2.3.4", "file.txt", file)
	require.Error(t, err)
}

func TestUploadFile_AlreadyClosed(t *testing.T) {
	client := NewClient("https://getfriend.ly")
	file := strings.NewReader("hello world")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.UploadFile(ctx, "1.2.3.4", "file.txt", file)
	require.Error(t, err)
}

func TestUploadFile_InvalidURL(t *testing.T) {
	client := NewClient("::invalid")
	file := strings.NewReader("hello world")

	_, err := client.UploadFile(context.Background(), "1.2.3.4", "file.txt", file)
	require.Error(t, err)
}

func TestUploadFile_NewRequestFailed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyFunc(http.NewRequestWithContext, func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
		return nil, fmt.Errorf("unreachable")
	})

	client := NewClient("https://getfriend.ly")
	file := strings.NewReader("hello world")

	_, err := client.UploadFile(context.Background(), "1.2.3.4", "file.txt", file)
	require.Error(t, err)
}

func TestDownloadFile_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/files/download/123/hash", r.URL.Path)
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte("hello"))
	}))
	defer ts.Close() //nolint:errcheck

	client := NewClient(ts.URL)
	fd := &FileDescriptor{Id: 123, AccessHash: "hash"}

	reader, err := client.DownloadFile(context.Background(), fd)
	require.NoError(t, err)
	defer reader.Close() //nolint:errcheck

	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "hello", string(data))
}

func TestDownloadFile_InvalidURL(t *testing.T) {
	fd := &FileDescriptor{Id: 1, AccessHash: FileAccessHash("hash")}
	client := NewClient("::invalid")

	_, err := client.DownloadFile(context.Background(), fd)
	require.Error(t, err)
}

func TestDownloadFile_Cancel(t *testing.T) {
	defer gock.Off()

	gock.New("https://getfriend.ly").
		Get("/files/download/123/hash").
		Reply(200).
		Delay(100 * time.Hour)

	client := NewClient("https://getfriend.ly")
	fd := &FileDescriptor{Id: 123, AccessHash: "hash"}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.DownloadFile(ctx, fd)
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestDownloadFile_StatusCodes(t *testing.T) {
	cases := []struct {
		name    string
		code    int
		wantErr error
	}{
		{"Not Found", 404, ErrNotFound},
		{"I'm a teapot", 418, nil},
	}

	for _, tc := range cases {
		defer gock.Off()
		gock.New("https://getfriend.ly").
			Get("/files/download/123/hash").
			Reply(tc.code)

		client := NewClient("https://getfriend.ly")
		fd := &FileDescriptor{Id: 123, AccessHash: "hash"}
		_, err := client.DownloadFile(context.Background(), fd)

		require.Error(t, err)
		if tc.wantErr != nil {
			require.ErrorIs(t, err, tc.wantErr)
		}
	}
}

func TestDownloadFile_CreateRequestFailed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyFunc(http.NewRequestWithContext, func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
		return nil, fmt.Errorf("unreachable")
	})

	client := NewClient("https://getfriend.ly")
	fd := &FileDescriptor{Id: 123, AccessHash: "hash"}

	_, err := client.DownloadFile(context.Background(), fd)
	require.Error(t, err)
}
