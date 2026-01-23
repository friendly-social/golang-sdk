package sdk

// FeedEntry represents single Entry from Feed.
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
func (c *Client) GetFeedQueue(auth *Authorization) (*FeedQueue, error) {
	var feed FeedQueue
	err := c.do("GET", "/feed/queue", auth, nil, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}
