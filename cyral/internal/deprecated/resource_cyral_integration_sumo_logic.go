package deprecated

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SumoLogicIntegration struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (data SumoLogicIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("name", data.Name)
	d.Set("address", data.Address)
	return nil
}

func (data *SumoLogicIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.Name = d.Get("name").(string)
	data.Address = d.Get("address").(string)
	return nil
}

func ResourceIntegrationSumoLogic() *schema.Resource {
	contextHandler := core.HTTPContextHandler{
		ResourceName:                 "SumoLogic Integration",
		ResourceType:                 resourcetype.Resource,
		SchemaReaderFactory:          func() core.SchemaReader { return &SumoLogicIntegration{} },
		SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &SumoLogicIntegration{} },
		BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/sumologic", c.ControlPlane)
		},
	}
	return &schema.Resource{
		DeprecationMessage: "Use resource `cyral_integration_logging` instead.",
		Description:        "Manages integration with [Sumo Logic to push sidecar logs](https://cyral.com/docs/integrations/siem/sumo-logic/).",
		CreateContext:      contextHandler.CreateContext(),
		ReadContext:        contextHandler.ReadContext(),
		UpdateContext:      contextHandler.UpdateContext(),
		DeleteContext:      contextHandler.DeleteContext(),
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
			"address": {
				Description: "Sumo Logic address.",
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
