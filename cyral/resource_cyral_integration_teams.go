package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MsTeamsIntegration struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (data MsTeamsIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("url", data.URL)
}

func (data *MsTeamsIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.URL = d.Get("url").(string)
}

var ReadMsTeamsConfig = ResourceOperationConfig{
	Name:       "MsTeamsResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &MsTeamsIntegration{},
}

func resourceIntegrationMsTeams() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "MsTeamsResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/teams", c.ControlPlane)
				},
				ResourceData: &MsTeamsIntegration{},
				ResponseData: &IDBasedResponse{},
			}, ReadMsTeamsConfig,
		),
		ReadContext: ReadResource(ReadMsTeamsConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "MsTeamsResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &MsTeamsIntegration{},
			}, ReadMsTeamsConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "MsTeamsResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the integration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
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
