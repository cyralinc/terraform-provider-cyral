package cyral

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	sampleSplunkIntegrationsDataStr = `{
  "id": "id1",
  "name": "name1",
  "label": "label1",
  "value": {
    "name": "name1",
    "host": "host1",
    "hecPort": 0,
    "accessToken": "accessToken1",
    "index": "index1",
    "useTLS": false,
    "cyralActivityLogsEnabled": false
  },
  "type": "splunk"
}
`
)

func sampleSplunkIntegrationsData() *IntegrationsData {
	return &IntegrationsData{
		Id:    "id1",
		Type:  "splunk",
		Name:  "name1",
		Label: "label1",
		Value: SplunkIntegration{
			Name:                     "name1",
			AccessToken:              "accessToken1",
			Port:                     0,
			Host:                     "host1",
			Index:                    "index1",
			UseTLS:                   false,
			CyralActivityLogsEnabled: false,
		},
	}
}

func TestGetValue_Default(t *testing.T) {
	integrationsData := NewDefaultIntegrationsData()
	expected := NewDefaultIntegrationsData().Value.(string)
	actual, err := integrationsData.getValue()
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestGetValue_Splunk(t *testing.T) {
	splunkIntegrationsData := sampleSplunkIntegrationsData()

	expectedBytes, err := json.Marshal(SplunkIntegration{
		Name:                     "name1",
		AccessToken:              "accessToken1",
		Port:                     0,
		Host:                     "host1",
		Index:                    "index1",
		UseTLS:                   false,
		CyralActivityLogsEnabled: false,
	})
	require.NoError(t, err)
	expected := string(expectedBytes)
	actual, err := splunkIntegrationsData.getValue()
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}
