package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
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

var ReadSumoLogicConfig = ResourceOperationConfig{
	Name:       "SumoLogicResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData { return &SumoLogicIntegration{} },
}

func resourceIntegrationSumoLogic() *schema.Resource {
	return &schema.Resource{
		Description: "Manages integration with [Sumo Logic to push sidecar logs](https://cyral.com/docs/integrations/siem/sumo-logic/).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "SumoLogicResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/sumologic", c.ControlPlane)
				},
				NewResourceData: func(_ *schema.ResourceData) ResourceData { return &SumoLogicIntegration{} },
				NewResponseData: func(_ *schema.ResourceData) ResponseData { return &IDBasedResponse{} },
			}, ReadSumoLogicConfig,
		),
		ReadContext: ReadResource(ReadSumoLogicConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "SumoLogicResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func(_ *schema.ResourceData) ResourceData { return &SumoLogicIntegration{} },
			}, ReadSumoLogicConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
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
			State: schema.ImportStatePassthrough,
		},
	}
}
