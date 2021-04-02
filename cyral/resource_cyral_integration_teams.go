package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateMsTeamsIntegrationResponse struct {
	ID string `json:"id"`
}

func (response CreateMsTeamsIntegrationResponse) WriteResourceData(d *schema.ResourceData) {
	d.SetId(response.ID)
}

func (response *CreateMsTeamsIntegrationResponse) ReadResourceData(d *schema.ResourceData) {
	response.ID = d.Id()
}

type MsTeamsIntegrationData struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (data MsTeamsIntegrationData) WriteResourceData(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("url", data.URL)
}

func (data *MsTeamsIntegrationData) ReadResourceData(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.URL = d.Get("url").(string)
}

var ReadMsTeamsFunctionConfig = FunctionConfig{
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
			FunctionConfig{
				Name:       "MsTeamsResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/teams", c.ControlPlane)
				},
				ResourceData: &MsTeamsIntegrationData{},
				ResponseData: &CreateMsTeamsIntegrationResponse{},
			}, ReadMsTeamsFunctionConfig,
		),
		ReadContext: ReadResource(ReadMsTeamsFunctionConfig),
		UpdateContext: UpdateResource(
			FunctionConfig{
				Name:       "MsTeamsResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &MsTeamsIntegrationData{},
			}, ReadMsTeamsFunctionConfig,
		),
		DeleteContext: DeleteResource(
			FunctionConfig{
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
