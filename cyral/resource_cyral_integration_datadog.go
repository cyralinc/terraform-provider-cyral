package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateDatadogIntegrationResponse struct {
	ID string `json:"ID"`
}

func (response CreateDatadogIntegrationResponse) WriteResourceData(d *schema.ResourceData) {
	d.SetId(response.ID)
}

func (response *CreateDatadogIntegrationResponse) ReadResourceData(d *schema.ResourceData) {
	response.ID = d.Id()
}

type DatadogIntegrationData struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

func (data DatadogIntegrationData) WriteResourceData(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("api_key", data.APIKey)
}

func (data *DatadogIntegrationData) ReadResourceData(d *schema.ResourceData) {
	data.ID = d.Id()
	data.Name = d.Get("name").(string)
	data.APIKey = d.Get("api_key").(string)
}

var ReadDatadogFunctionConfig = FunctionConfig{
	Name:       "DatadogResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &DatadogIntegrationData{},
}

func resourceIntegrationDatadog() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			FunctionConfig{
				Name:       "DatadogResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/datadog", c.ControlPlane)
				},
				ResourceData: &DatadogIntegrationData{},
				ResponseData: &CreateDatadogIntegrationResponse{},
			}, ReadDatadogFunctionConfig,
		),
		ReadContext: ReadResource(ReadDatadogFunctionConfig),
		UpdateContext: UpdateResource(
			FunctionConfig{
				Name:       "DatadogResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &DatadogIntegrationData{},
			}, ReadDatadogFunctionConfig,
		),
		DeleteContext: DeleteResource(
			FunctionConfig{
				Name:       "DatadogResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())
				},
			},
		),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
