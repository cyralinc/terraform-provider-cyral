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

var CreateSlackAlertsFunctionConfig = FunctionConfig{
	Name:       "SlackAlertsResourceCreate",
	HttpMethod: http.MethodPost,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/slack", c.ControlPlane)
	},
	ResourceData:       &SlackAlertsIntegrationData{},
	ResponseData:       &CreateSlackAlertsIntegrationResponse{},
	ReadFunctionConfig: &ReadSlackAlertsFunctionConfig,
}

var ReadSlackAlertsFunctionConfig = FunctionConfig{
	Name:       "SlackAlertsResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &SlackAlertsIntegrationData{},
}

var UpdateSlackAlertsFunctionConfig = FunctionConfig{
	Name:       "SlackAlertsResourceUpdate",
	HttpMethod: http.MethodPut,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
	},
	ResourceData:       &SlackAlertsIntegrationData{},
	ReadFunctionConfig: &ReadSlackAlertsFunctionConfig,
}

var DeleteSlackAlertsFunctionConfig = FunctionConfig{
	Name:       "SlackAlertsResourceDelete",
	HttpMethod: http.MethodDelete,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
	},
}

func resourceIntegrationSlackAlerts() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateSlackAlertsFunctionConfig.Create,
		ReadContext:   ReadSlackAlertsFunctionConfig.Read,
		UpdateContext: UpdateSlackAlertsFunctionConfig.Update,
		DeleteContext: DeleteSlackAlertsFunctionConfig.Delete,

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
