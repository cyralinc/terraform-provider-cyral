package cyral

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WriteToSchema writes the pager duty information contained in a PagerDutyIntegration to the TF schema
func (data PagerDutyIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("id", data.ID)
	d.Set("name", data.Name)

	// the API takes a marshalled json as a field for defining specific variables for a given auth policy integration
	var token struct {
		APIToken string `json:"apiToken"`
	}

	err := json.Unmarshal([]byte(data.Parameters), &token)
	if err != nil {
		panic(err)
	}

	d.Set("api_token", token.APIToken)
}

// ReadFromSchema reads the pager duty information from a given schema
func (data *PagerDutyIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.ID = d.Get("id").(string)
	data.Name = d.Get("name").(string)

	var token struct {
		APIToken string `json:"apiToken"`
	}
	token.APIToken = d.Get("api_token").(string)
	msl, _ := json.Marshal(token)
	data.Parameters = string(msl)

	// API information for typing
	data.BuiltInType = "pagerduty"
	data.Category = "builtin"
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
