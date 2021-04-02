package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MsTeamsIntegrationData struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (data MsTeamsIntegrationData) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("url", data.URL)
}

func (data *MsTeamsIntegrationData) ReadFromSchema(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.URL = d.Get("url").(string)
}

var ReadMsTeamsConfig = ResourceOperationConfig{
	Name:       "MsTeamsResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &MsTeamsIntegrationData{},
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
				ResourceData: &MsTeamsIntegrationData{},
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
				ResourceData: &MsTeamsIntegrationData{},
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
