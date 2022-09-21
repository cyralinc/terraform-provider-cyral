package cyral

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

func listSidecars(c *client.Client) ([]IdentifiedSidecarInfo, error) {
	log.Printf("[DEBUG] Init listSidecars")
	url := fmt.Sprintf("https://%s/v1/sidecars", c.ControlPlane)
	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var sidecarsInfo []IdentifiedSidecarInfo
	if err := json.Unmarshal(body, &sidecarsInfo); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", sidecarsInfo)
	log.Printf("[DEBUG] End listSidecars")

	return sidecarsInfo, nil
}

func listRoles(c *client.Client) (*GetUserGroupsResponse, error) {
	log.Printf("[DEBUG] Init listRoles")

	url := fmt.Sprintf("https://%s/v1/users/groups", c.ControlPlane)
	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	resp := &GetUserGroupsResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", resp)
	log.Printf("[DEBUG] End listRoles")

	return resp, nil
}

func listIdPIntegrations(c *client.Client) (*IdPIntegrations, error) {
	log.Printf("[DEBUG] Init listIdPIntegrations")

	url := fmt.Sprintf("https://%s/v1/integrations/saml", c.ControlPlane)
	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	resp := &IdPIntegrations{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", resp)
	log.Printf("[DEBUG] End listIdPIntegrations")

	return resp, nil
}
