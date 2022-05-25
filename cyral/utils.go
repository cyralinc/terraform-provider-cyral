package cyral

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type CreateSidecarCredentialsRequest struct {
	SidecarID string `json:"sidecarId"`
}

func createSidecarCredentials(c *client.Client, sidecarId string) (
	*SidecarCredentialsData, error) {

	payload := CreateSidecarCredentialsRequest{sidecarId}

	url := fmt.Sprintf("https://%s/v1/users/sidecarAccounts", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, payload)
	if err != nil {
		return nil, fmt.Errorf("error when performing request: %w", err)
	}

	response := &SidecarCredentialsData{}
	if err := json.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("error when unmarshalling JSON: %w", err)
	}

	return response, nil
}

func urlQuery(kv map[string]string) string {
	queryStr := "?"
	for k, v := range kv {
		queryStr += fmt.Sprintf("&%s=%s", k, v)
	}
	return queryStr
}
