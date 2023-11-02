package deprecated

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	integrationTypeSplunk = "splunk"
)

type IntegrationsData struct {
	Id    string      `json:"id"`
	Type  string      `json:"type"`
	Name  string      `json:"name"`
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

func NewDefaultIntegrationsData() *IntegrationsData {
	return &IntegrationsData{
		Id:    "id",
		Type:  "default",
		Name:  "default",
		Label: "default",
		Value: "default",
	}
}

func (isd *IntegrationsData) GetValue() (string, error) {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("unable to get integration value for "+
				"type '%s': %w", isd.Type, err)
		}
	}()
	switch isd.Type {
	case integrationTypeSplunk:
		if value, ok := isd.Value.(SplunkIntegration); ok {
			bytesval, jsonErr := json.Marshal(value)
			if jsonErr == nil {
				return string(bytesval), nil
			}
			err = fmt.Errorf("error marshalling splunk "+
				"integration: %w", jsonErr)
		}
	default:
		if value, ok := isd.Value.(string); ok {
			return value, nil
		}
		err = errors.New("value is not of string type")
	}
	return "", err
}
