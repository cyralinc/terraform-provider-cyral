package cyral

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WriteToSchema writes the pager duty information contained in a PagerDutyIntegration to the TF schema
func (data PagerDutyIntegration) WriteToSchema(d *schema.ResourceData) error {
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

	return nil
}

// ReadFromSchema reads the pager duty information from a given schema
func (data *PagerDutyIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.ID = d.Get("id").(string)
	data.Name = d.Get("name").(string)

	var token struct {
		APIToken string `json:"apiToken"`
	}
	token.APIToken = d.Get("api_token").(string)
	msl, _ := json.Marshal(token)
	data.Parameters = string(msl)

	data.Purpose = "authorization"

	// API information for typing
	data.TemplateType = "pagerduty"
	data.Category = "builtin"

	return nil
}

// ReadPagerDutyIntegrationConfig is the configuration to read from a pager duty resource contained in a Control Plane
var ReadPagerDutyIntegrationConfig = ResourceOperationConfig{
	Name:       "PagerDutyIntegrationResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(
			"https://%s/v1/integrations/confExtensions/instances/authorization/%s",
			c.ControlPlane, d.Id(),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData { return &PagerDutyIntegration{} },
}

func resourceIntegrationPagerDuty() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with PagerDuty](https://cyral.com/docs/integrations/incident-response/pagerduty/#in-cyral).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "PagerDutyIntegrationResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/integrations/confExtensions/instances", c.ControlPlane,
					)
				},
				NewResourceData: func() ResourceData { return &PagerDutyIntegration{} },
				NewResponseData: func(_ *schema.ResourceData) ResponseData { return &IDBasedResponse{} },
			}, ReadPagerDutyIntegrationConfig,
		),
		ReadContext: ReadResource(ReadPagerDutyIntegrationConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "PagerDutyIntegrationResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/integrations/confExtensions/instances/%s", c.ControlPlane, d.Id(),
					)
				},
				NewResourceData: func() ResourceData { return &PagerDutyIntegration{} },
			}, ReadPagerDutyIntegrationConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "PagerDutyIntegrationResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/integrations/confExtensions/instances/authorization/%s",
						c.ControlPlane, d.Id(),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Integration name that will be used internally in the control plane.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"api_token": {
				Description: "API token for the PagerDuty integration.",
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
