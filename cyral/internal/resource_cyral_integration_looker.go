package internal

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LookerIntegration struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	URL          string `json:"url"`
}

func (data LookerIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("client_secret", data.ClientSecret)
	d.Set("client_id", data.ClientId)
	d.Set("url", data.URL)
	return nil
}

func (data *LookerIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.ClientSecret = d.Get("client_secret").(string)
	data.ClientId = d.Get("client_id").(string)
	data.URL = d.Get("url").(string)
	return nil
}

var ReadLookerConfig = core.ResourceOperationConfig{
	Name:       "LookerResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())
	},
	NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &LookerIntegration{} },
}

func ResourceIntegrationLooker() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Integration no longer supported.",
		Description:        "Manages integration with Looker.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "LookerResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/looker", c.ControlPlane)
				},
				NewResourceData: func() core.ResourceData { return &LookerIntegration{} },
				NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &core.IDBasedResponse{} },
			}, ReadLookerConfig,
		),
		ReadContext: core.ReadResource(ReadLookerConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "LookerResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() core.ResourceData { return &LookerIntegration{} },
			}, ReadLookerConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "LookerResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"client_id": {
				Description: "Looker client id.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"client_secret": {
				Description: "Looker client secret.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"url": {
				Description: "Looker integration url.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
