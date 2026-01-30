package sdk

import "context"

// NetworkDetails represents details about user's network, particularly their friends list.
type NetworkDetails struct {
	Friends []UserDetails `json:"friends"`
}

// GetNetworkDetails returns NetworkDetails structure for provided Authorization.
func (c *Client) GetNetworkDetails(ctx context.Context, auth *Authorization) (*NetworkDetails, error) {
	var network NetworkDetails
	err := c.do(ctx, "GET", "/network/details", auth, nil, &network)
	if err != nil {
		return nil, err
	}

	return &network, nil
}
