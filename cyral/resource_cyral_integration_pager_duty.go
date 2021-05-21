package cyral

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// UnmarshalJSON unmarshalls the received data from the API into a pagerDutyIntegration struct. Since the API is not
// 1 to 1 with the resources, we need this extra steps to recover the api token from the server.
func (data *PagerDutyIntegration) UnmarshalJSON(b []byte) error {
	tokenInfo := struct {
		APIToken string `json:"apiToken"`
	}{}

	// the rest of the information is unecessary
	response := struct {
		Name       string `json:"name"`
		Parameters string `json:"parameters"`
	}{}

	err := json.Unmarshal(b, &response)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(response.Parameters), &tokenInfo)
	if err != nil {
		return err
	}

	data.Name = response.Name
	data.APIToken = tokenInfo.APIToken
	return nil
}

// MarshalJSON marshalls the structure into a json. Since the API is not
// 1 to 1 with the resources, we need this extra steps to recover the api token from the server.
func (data PagerDutyIntegration) MarshalJSON() ([]byte, error) {
	// need to marshal APIToken twice because field on the API is a marshalled JSON containing necessary information
	tokenBytes, _ := json.Marshal(struct {
		Token string `json:"apiToken"`
	}{
		data.APIToken,
	})

	return json.Marshal(struct {
		Name        string `json:"name"`
		Category    string `json:"category"`
		BuiltInType string `json:"builtInType"`
		Parameters  string `json:"parameters"`
	}{
		Name:        data.Name,
		Category:    "builtin",
		BuiltInType: "pagerduty",
		Parameters:  string(tokenBytes),
	})
}

// WriteToSchema writes the pager duty information contained in a PagerDutyIntegration to the TF schema
func (data PagerDutyIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("id", data.ID)
	d.Set("name", data.Name)
	d.Set("api_token", data.APIToken)
}

// ReadFromSchema reads the pager duty information from a given schema
func (data *PagerDutyIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.ID = d.Get("id").(string)
	data.Name = d.Get("name").(string)
	data.APIToken = d.Get("api_token").(string)
}

// ReadPagerDutyIntegrationConfig is the configuration to read from a pager duty resource contained in a Control Plane
var ReadPagerDutyIntegrationConfig = ResourceOperationConfig{
	Name:       "PagerDutyIntegrationResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/authorizationPolicies/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &PagerDutyIntegration{},
}

func resourceIntegrationPagerDuty() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "PagerDutyIntegrationResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/authorizationPolicies", c.ControlPlane)
				},
				ResourceData: &PagerDutyIntegration{},
				ResponseData: &IDBasedResponse{},
			}, ReadPagerDutyIntegrationConfig,
		),
		ReadContext: ReadResource(ReadPagerDutyIntegrationConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "PagerDutyIntegrationResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/authorizationPolicies/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &PagerDutyIntegration{},
			}, ReadPagerDutyIntegrationConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "PagerDutyIntegrationResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/authorizationPolicies/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"api_token": {
				Required:  true,
				Sensitive: true,
				Type:      schema.TypeString,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
