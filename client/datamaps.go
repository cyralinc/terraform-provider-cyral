package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// GetDatamapByLabel - Returns the datamap filtered by a specific label
func (c *Client) GetDatamapByLabel(label string) (*DataMap, error) {
	url := fmt.Sprintf("https://%s/v1/datamaps?label=%s", c.ControlPlane, label)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	datamap := DataMap{}
	if err := json.Unmarshal(body, &datamap); err != nil {
		return nil, err
	}

	return &datamap, nil
}

// GetDatamap - Returns the complete datamap
func (c *Client) GetDatamap() (*DataMap, error) {
	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sensitiveData := make(SensitiveData)
	if string(body) != "{\"Message\":\"OK\"}" {
		log.Printf("[DEBUG] GetDatamap IF ----------------------")
		log.Printf("[DEBUG] GetDatamap ------ body: %s", string(body))
		if err := json.Unmarshal(body, &sensitiveData); err != nil {
			log.Printf("[DEBUG] GetDatamap ------ ERRO UNMARSHAL: %v", err)
			return nil, err
		}
	}

	datamap := DataMap{SensitiveData: sensitiveData}
	log.Printf("[DEBUG] GetDatamap --- datamap: %#v", sensitiveData)
	sd := MakeStrForSD(sensitiveData)
	log.Printf("[DEBUG] GetDatamap --- sensitiveData: %s", sd)
	return &datamap, nil
}

// CreateDatamap creates a datamap
func (c *Client) CreateDatamap(sensitiveData SensitiveData) error {
	payloadBytes, err := json.Marshal(sensitiveData)

	if err != nil {
		return fmt.Errorf("failed to encode 'create datamap' payload; err: %v", err)
	}

	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Errorf("unable to create 'create datamap' request; err: %v", err)
	}

	if _, err := c.doRequest(req); err != nil {
		return err
	}

	return nil
}

// UpdateDatamap - Updates a datamap
// func (c *Client) UpdateDatamap(datamapLabels []DatamapLabel) (*DataMap, error) {
// 	payloadBytes, err := json.Marshal(datamapLabels)
// 	if err != nil {
// 		return nil, err
// 	}

// 	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

// 	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(payloadBytes)))
// 	if err != nil {
// 		return nil, err
// 	}

// 	body, err := c.doRequest(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	datamap := Datamap{}
// 	if err := json.Unmarshal(body, &datamap); err != nil {
// 		return nil, err
// 	}

// 	return &datamap, nil
// }

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
