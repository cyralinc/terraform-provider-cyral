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
	if err := json.Unmarshal(body, &datamap); err != nil {
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
	if err := json.Unmarshal(body, &datamap); err != nil {
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

	datamap := Datamap{}
	if err := json.Unmarshal(body, &datamap); err != nil {
		return nil, fmt.Errorf("unable to unmarshall json; err: %v", err)
	}

	return &datamap, nil
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
	if err := json.Unmarshal(body, &datamap); err != nil {
		return nil, err
	}

	return &datamap, nil
}

// DeleteDatamapLabel - Deletes a certain label of the datamap
func (c *Client) DeleteDatamapLabel(label string) error {
	url := fmt.Sprintf("https://%s/v1/datamaps?label=%s", c.ControlPlane, label)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	if _, err := c.doRequest(req); err != nil {
		return err
	}

	// NOTE: I think this does not apply here,
	// 		 since the delete api does not return a response body
	//
	// if string(body) != "Deleted datamap" {
	// 	return errors.New(string(body))
	// }

	return nil
}

// DeleteDatamap - Deletes the complete datamap
func (c *Client) DeleteDatamap() error {
	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	if _, err := c.doRequest(req); err != nil {
		return err
	}

	// NOTE: I think this does not apply here,
	// 		 since the delete api does not return a response body
	//
	// if string(body) != "Deleted datamap" {
	// 	return errors.New(string(body))
	// }

	return nil
}
