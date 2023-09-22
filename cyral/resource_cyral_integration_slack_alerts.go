package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SlackAlertsIntegration struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (data SlackAlertsIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("name", data.Name)
	d.Set("url", data.URL)
	return nil
}

func (data *SlackAlertsIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.Name = d.Get("name").(string)
	data.URL = d.Get("url").(string)
	return nil
}

var ReadSlackAlertsConfig = ResourceOperationConfig{
	Name:       "SlackAlertsResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
	},
	NewResponseData:     func(_ *schema.ResourceData) ResponseData { return &SlackAlertsIntegration{} },
	RequestErrorHandler: &ReadIgnoreHttpNotFound{resName: "Integration Slack"},
}

func resourceIntegrationSlackAlerts() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with Slack to push alerts](https://cyral.com/docs/integrations/messaging/slack).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "SlackAlertsResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack", c.ControlPlane)
				},
				NewResourceData: func() ResourceData { return &SlackAlertsIntegration{} },
				NewResponseData: func(_ *schema.ResourceData) ResponseData { return &IDBasedResponse{} },
			}, ReadSlackAlertsConfig,
		),
		ReadContext: ReadResource(ReadSlackAlertsConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "SlackAlertsResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() ResourceData { return &SlackAlertsIntegration{} },
			}, ReadSlackAlertsConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "SlackAlertsResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
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
				Description: "Slack Alert Webhook url.",
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
