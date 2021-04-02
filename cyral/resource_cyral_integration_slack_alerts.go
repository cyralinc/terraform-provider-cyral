package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateSlackAlertsIntegrationResponse struct {
	ID string `json:"id"`
}

func (response CreateSlackAlertsIntegrationResponse) WriteResourceData(d *schema.ResourceData) {
	d.SetId(response.ID)
}

func (response *CreateSlackAlertsIntegrationResponse) ReadResourceData(d *schema.ResourceData) {
	response.ID = d.Id()
}

type SlackAlertsIntegrationData struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (data SlackAlertsIntegrationData) WriteResourceData(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("url", data.URL)
}

func (data *SlackAlertsIntegrationData) ReadResourceData(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.URL = d.Get("url").(string)
}

var ReadSlackAlertsFunctionConfig = FunctionConfig{
	Name:       "SlackAlertsResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &SlackAlertsIntegrationData{},
}

func resourceIntegrationSlackAlerts() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			FunctionConfig{
				Name:       "SlackAlertsResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack", c.ControlPlane)
				},
				ResourceData: &SlackAlertsIntegrationData{},
				ResponseData: &CreateSlackAlertsIntegrationResponse{},
			}, ReadSlackAlertsFunctionConfig,
		),
		ReadContext: ReadResource(ReadSlackAlertsFunctionConfig),
		UpdateContext: UpdateResource(
			FunctionConfig{
				Name:       "SlackAlertsResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &SlackAlertsIntegrationData{},
			}, ReadSlackAlertsFunctionConfig,
		),
		DeleteContext: DeleteResource(
			FunctionConfig{
				Name:       "SlackAlertsResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
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
