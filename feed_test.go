package sdk

import (
	"testing"

	"github.com/h2non/gock"
)

func TestGetFeed(t *testing.T) {
	f := func() *gock.Request {
		return gock.New("https://getfriend.ly").Get("/feed/queue")
	}

	call := func(c *Client) (*FeedQueue, error) {
		return c.GetFeedQueue(&Authorization{Id: 1, AccessHash: "1", Token: "1"})
	}

	RunAPITests(t, append(CommonCasesTests(f, call),
		APITestCase[FeedQueue]{
			name:             "Success",
			setup:            func() { f().Reply(200).JSON(FeedQueue{Entries: []FeedEntry{}}) },
			call:             call,
			expectError:      false,
			expectedResponse: &FeedQueue{Entries: []FeedEntry{}},
		},
	))
}
