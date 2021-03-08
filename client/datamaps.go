package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetDatamapByLabel - Returns the datamap filtered by a specific label
func (c *Client) GetDatamapByLabel(label string) (*Datamap, error) {
	url := fmt.Sprintf("https://%s/v1/datamaps?label=%s", c.ControlPlane, label)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	datamap := Datamap{}
	err = json.Unmarshal(body, &datamap)
	if err != nil {
		return nil, err
	}

	return &datamap, nil
}

// GetDatamap - Returns the complete datamap
func (c *Client) GetDatamap() (*Datamap, error) {
	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	datamap := Datamap{}
	err = json.Unmarshal(body, &datamap)
	if err != nil {
		return nil, err
	}

	return &datamap, nil
}

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

// UpdateDatamap - Updates a datamap
func (c *Client) UpdateDatamap(datamapLabels []DatamapLabel) (*Datamap, error) {
	payloadBytes, err := json.Marshal(datamapLabels)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	datamap := Datamap{}
	err = json.Unmarshal(body, &datamap)
	if err != nil {
		return nil, err
	}

	return &datamap, nil
}
