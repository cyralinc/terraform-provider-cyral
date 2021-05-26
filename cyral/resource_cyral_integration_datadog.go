package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DatadogIntegration struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

func (data DatadogIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("api_key", data.APIKey)
}

func (data *DatadogIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.ID = d.Id()
	data.Name = d.Get("name").(string)
	data.APIKey = d.Get("api_key").(string)
}

var ReadDatadogConfig = ResourceOperationConfig{
	Name:       "DatadogResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &DatadogIntegration{},
}

func resourceIntegrationDatadog() *schema.Resource {
	return &schema.Resource{
		Description: "CRUD operations for Datadog integration",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "DatadogResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/datadog", c.ControlPlane)
				},
				ResourceData: &DatadogIntegration{},
				ResponseData: &IDBasedResponse{},
			}, ReadDatadogConfig,
		),
		ReadContext: ReadResource(ReadDatadogConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "DatadogResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &DatadogIntegration{},
			}, ReadDatadogConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "DatadogResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())
				},
			},
		),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Integration name that will be used internally in Control Plane",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Datadog API key",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
