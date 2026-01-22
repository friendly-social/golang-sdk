package sdk

// NetworkDetails represents details about user's network, particularly their Friends list.
type NetworkDetails struct {
	Friends []UserDetails `json:"friends"`
}

// GetNetworkDetails returns NetworkDetails structure for provided Authorization.
func (c *Client) GetNetworkDetails(auth *Authorization) (*NetworkDetails, error) {
	var network NetworkDetails
	err := c.do("GET", "/network/details", auth, nil, &network)
	if err != nil {
		return nil, err
	}

	return &network, nil
}
