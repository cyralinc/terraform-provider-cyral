package slack

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
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

var ReadSlackAlertsConfig = core.ResourceOperationConfig{
	Name:       "SlackAlertsResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
	},
	NewResponseData:     func(_ *schema.ResourceData) core.ResponseData { return &SlackAlertsIntegration{} },
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Integration Slack"},
}

func ResourceIntegrationSlackAlerts() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with Slack to push alerts](https://cyral.com/docs/integrations/messaging/slack).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "SlackAlertsResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack", c.ControlPlane)
				},
				NewResourceData: func() core.ResourceData { return &SlackAlertsIntegration{} },
				NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &core.IDBasedResponse{} },
			}, ReadSlackAlertsConfig,
		),
		ReadContext: core.ReadResource(ReadSlackAlertsConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "SlackAlertsResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() core.ResourceData { return &SlackAlertsIntegration{} },
			}, ReadSlackAlertsConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
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
