package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FeedEntry represents single Entry from Feed.
type FeedEntry struct {
	IsExtendedNetwork bool          `json:"isExtendedNetwork"`
	CommonFriends     []UserDetails `json:"commonFriends"`
	Details           UserDetails   `json:"details"`
}

// FeedQueue represents queue of feed entries which must be showed.
type FeedQueue struct {
	Entries []FeedEntry `json:"entries"`
}

// GetFeedQueue returns FeedQueue structure for provided Authorization.
func (c *Client) GetFeedQueue(auth *Authorization) (*FeedQueue, error) {
	resp, err := c.do("GET", "/feed/queue", auth, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get feed failed: status %d", resp.StatusCode)
	}

	var feed FeedQueue
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	return &feed, nil
}
