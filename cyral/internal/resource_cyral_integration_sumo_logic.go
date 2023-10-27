package internal

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
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

var ReadSumoLogicConfig = core.ResourceOperationConfig{
	Name:       "SumoLogicResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())
	},
	NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &SumoLogicIntegration{} },
}

func ResourceIntegrationSumoLogic() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use resource `cyral_integration_logging` instead.",
		Description:        "Manages integration with [Sumo Logic to push sidecar logs](https://cyral.com/docs/integrations/siem/sumo-logic/).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "SumoLogicResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/sumologic", c.ControlPlane)
				},
				NewResourceData: func() core.ResourceData { return &SumoLogicIntegration{} },
				NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &core.IDBasedResponse{} },
			}, ReadSumoLogicConfig,
		),
		ReadContext: core.ReadResource(ReadSumoLogicConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "SumoLogicResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() core.ResourceData { return &SumoLogicIntegration{} },
			}, ReadSumoLogicConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "SumoLogicResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())
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
