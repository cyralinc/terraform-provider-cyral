package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SumoLogicIntegrationData struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (data SumoLogicIntegrationData) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("address", data.Address)

}

func (data *SumoLogicIntegrationData) ReadFromSchema(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.Address = d.Get("address").(string)
}

var ReadSumoLogicConfig = ResourceOperationConfig{
	Name:       "SumoLogicResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &SumoLogicIntegrationData{},
}

func resourceIntegrationSumoLogic() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "SumoLogicResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/sumologic", c.ControlPlane)
				},
				ResourceData: &SumoLogicIntegrationData{},
				ResponseData: &IDBasedResponse{},
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
				ResourceData: &SumoLogicIntegrationData{},
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
