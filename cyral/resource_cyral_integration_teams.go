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

func (data MsTeamsIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("name", data.Name)
	d.Set("url", data.URL)
	return nil
}

func (data *MsTeamsIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.Name = d.Get("name").(string)
	data.URL = d.Get("url").(string)
	return nil
}

var ReadMsTeamsConfig = ResourceOperationConfig{
	Name:       "MsTeamsResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData { return &MsTeamsIntegration{} },
}

func resourceIntegrationMsTeams() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with Microsoft Teams](https://cyral.com/docs/integrations/messaging/microsoft-teams/).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "MsTeamsResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/teams", c.ControlPlane)
				},
				NewResourceData: func(_ *schema.ResourceData) ResourceData { return &MsTeamsIntegration{} },
				NewResponseData: func(_ *schema.ResourceData) ResponseData { return &IDBasedResponse{} },
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
				NewResourceData: func(_ *schema.ResourceData) ResourceData { return &MsTeamsIntegration{} },
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
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Integration name that will be used internally in the control plane.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "Microsoft Teams webhook URL.",
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
