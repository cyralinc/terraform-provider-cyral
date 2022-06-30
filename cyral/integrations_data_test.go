package cyral

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestIntegrationsData_GetValue_Default(t *testing.T) {
	integrationsData := NewDefaultIntegrationsData()
	expected := NewDefaultIntegrationsData().Value.(string)
	actual, err := integrationsData.GetValue()
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestIntegrationsData_GetValue_Splunk(t *testing.T) {
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
	actual, err := splunkIntegrationsData.GetValue()
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}