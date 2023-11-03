package deprecated

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DatadogIntegration struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

func (data DatadogIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("name", data.Name)
	d.Set("api_key", data.APIKey)
	return nil
}

func (data *DatadogIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.ID = d.Id()
	data.Name = d.Get("name").(string)
	data.APIKey = d.Get("api_key").(string)
	return nil
}

var ReadDatadogConfig = core.ResourceOperationConfig{
	Name:       "DatadogResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())
	},
	NewResponseData:     func(_ *schema.ResourceData) core.ResponseData { return &DatadogIntegration{} },
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Integration datadog"},
}

func ResourceIntegrationDatadog() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "If configuring Datadog for logging purposes, use resource `cyral_integration_logging` instead.",
		Description: "Manages [integration with DataDog](https://cyral.com/docs/integrations/apm/datadog/) " +
			"to push sidecar logs and/or metrics.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "DatadogResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/datadog", c.ControlPlane)
				},
				NewResourceData: func() core.ResourceData { return &DatadogIntegration{} },
				NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &core.IDBasedResponse{} },
			}, ReadDatadogConfig,
		),
		ReadContext: core.ReadResource(ReadDatadogConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "DatadogResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() core.ResourceData { return &DatadogIntegration{} },
			}, ReadDatadogConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "DatadogResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())
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
				Type:        schema.TypeString,
				Required:    true,
			},
			"api_key": {
				Description: "Datadog API key.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
