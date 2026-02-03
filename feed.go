package sdk

import (
	"context"
	"fmt"
)

// FeedEntry represents single entry from the Feed.
type FeedEntry struct {
	IsRequest         bool          `json:"isRequest"`
	IsExtendedNetwork bool          `json:"isExtendedNetwork"`
	CommonFriends     []UserDetails `json:"commonFriends"`
	Details           UserDetails   `json:"details"`
}

// FeedQueue represents queue of feed entries which must be showed.
type FeedQueue struct {
	Entries []FeedEntry `json:"entries"`
}

// GetFeedQueue returns FeedQueue structure for provided Authorization.
func (c *Client) GetFeedQueue(ctx context.Context, auth *Authorization) (*FeedQueue, error) {
	var feed FeedQueue
	err := c.do(ctx, auth, "GET", "/feed/queue", nil, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed queue: %w", err)
	}

	return &feed, nil
}
