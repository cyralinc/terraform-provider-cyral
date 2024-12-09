package deprecated

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
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

func ResourceIntegrationDatadog() *schema.Resource {
	contextHandler := core.HTTPContextHandler{
		ResourceName:                 "Datadog Integration",
		ResourceType:                 resourcetype.Resource,
		SchemaReaderFactory:          func() core.SchemaReader { return &DatadogIntegration{} },
		SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &DatadogIntegration{} },
		BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/datadog", c.ControlPlane)
		},
	}
	return &schema.Resource{
		DeprecationMessage: "If configuring Datadog for logging purposes, use resource `cyral_integration_logging` instead.",
		Description: "Manages [integration with DataDog](https://cyral.com/docs/integrations/siem/datadog-logs) " +
			"to push sidecar logs and/or metrics.",
		CreateContext: contextHandler.CreateContext(),
		ReadContext:   contextHandler.ReadContext(),
		UpdateContext: contextHandler.UpdateContext(),
		DeleteContext: contextHandler.DeleteContext(),
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
