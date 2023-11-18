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
	ResourceName: "SlackAlertsResourceRead",
	HttpMethod:   http.MethodGet,
	URLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
	},
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &SlackAlertsIntegration{} },
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Integration Slack"},
}

func ResourceIntegrationSlackAlerts() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with Slack to push alerts](https://cyral.com/docs/integrations/messaging/slack).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				ResourceName: "SlackAlertsResourceCreate",
				HttpMethod:   http.MethodPost,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack", c.ControlPlane)
				},
				SchemaReaderFactory: func() core.SchemaReader { return &SlackAlertsIntegration{} },
			}, ReadSlackAlertsConfig,
		),
		ReadContext: core.ReadResource(ReadSlackAlertsConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				ResourceName: "SlackAlertsResourceUpdate",
				HttpMethod:   http.MethodPut,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())
				},
				SchemaReaderFactory: func() core.SchemaReader { return &SlackAlertsIntegration{} },
			}, ReadSlackAlertsConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				ResourceName: "SlackAlertsResourceDelete",
				HttpMethod:   http.MethodDelete,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
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
