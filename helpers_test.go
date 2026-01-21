package sdk

import (
	"errors"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

type ValueTypeTestCase[T, Y any] struct {
	name    string
	new     func(value T) (Y, error)
	must    func(value T) Y
	valid   T
	invalid T
}

type APITestCase[T any] struct {
	name             string
	setup            func()
	call             func(c *Client) (*T, error)
	expectError      bool
	expectedResponse *T
}

func RunValueTypeTest[T, Y any](t *testing.T, test ValueTypeTestCase[T, Y]) {
	t.Helper()
	t.Run(test.name, func(t *testing.T) {
		v, err := test.new(test.valid)
		require.NoError(t, err)
		require.EqualValues(t, test.valid, v)

		_, err = test.new(test.invalid)
		require.Error(t, err)

		if test.must != nil {
			require.Panics(t, func() { _ = test.must(test.invalid) })
			require.NotPanics(t, func() {
				v = test.must(test.valid)
				require.EqualValues(t, test.valid, v)
			})
		}
	})
}

func RunAPITests[T any](t *testing.T, tests []APITestCase[T]) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.setup()

			client := NewClient("https://getfriend.ly")
			gock.InterceptClient(client.http)
			resp, err := tt.call(client)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedResponse, resp)
		})
	}
}

func CommonCasesTests[T any](f func() *gock.Request, call func(c *Client) (*T, error)) []APITestCase[T] {
	return []APITestCase[T]{
		{
			name:        "Invalid JSON",
			setup:       func() { f().Reply(200).BodyString("something went wrong") },
			call:        call,
			expectError: true,
		},
		{
			name:        "Unauthorized",
			setup:       func() { f().Reply(401) },
			call:        call,
			expectError: true,
		},
		{
			name:        "Server Error",
			setup:       func() { f().Reply(500) },
			call:        call,
			expectError: true,
		},
		{
			name:        "Network Fail",
			setup:       func() { f().ReplyError(errors.New("something went wrong")) },
			call:        call,
			expectError: true,
		},
	}
}
