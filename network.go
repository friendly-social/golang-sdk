package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// NetworkDetails represents details about user's network, particularly their Friends list.
type NetworkDetails struct {
	Friends []UserDetails `json:"friends"`
}

// GetNetworkDetails returns NetworkDetails structure for provided Authorization.
func (c *Client) GetNetworkDetails(auth *Authorization) (*NetworkDetails, error) {
	resp, err := c.do("GET", "/network/details", auth, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get network failed: status %d", resp.StatusCode)
	}

	var network NetworkDetails
	if err := json.NewDecoder(resp.Body).Decode(&network); err != nil {
		return nil, err
	}

	return &network, nil
}
