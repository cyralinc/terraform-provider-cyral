package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// CreateDatamap creates a datamap
func (c *Client) CreateDatamap(datamapLabels []DatamapLabel) (*Datamap, error) {
	payloadBytes, err := json.Marshal(datamapLabels)
	if err != nil {
		return nil, fmt.Errorf("failed to encode 'create datamap' payload; err: %v", err)
	}

	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("unable to create 'create datamap' request; err: %v", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, nil); err != nil {
		return nil, fmt.Errorf("unable to unmarshall json; err: %v", err)
	}

	return nil, nil
}
