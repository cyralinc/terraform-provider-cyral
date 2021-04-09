package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LookerIntegration struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	URL          string `json:"url"`
}

func (data LookerIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("client_secret", data.ClientSecret)
	d.Set("client_id", data.ClientId)
	d.Set("url", data.URL)
}

func (data *LookerIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.ClientSecret = d.Get("client_secret").(string)
	data.ClientId = d.Get("client_id").(string)
	data.URL = d.Get("url").(string)
}

var ReadLookerConfig = ResourceOperationConfig{
	Name:       "LookerResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &LookerIntegration{},
}

func resourceIntegrationLooker() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "LookerResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/looker", c.ControlPlane)
				},
				ResourceData: &LookerIntegration{},
				ResponseData: &IDBasedResponse{},
			}, ReadLookerConfig,
		),
		ReadContext: ReadResource(ReadLookerConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "LookerResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &LookerIntegration{},
			}, ReadLookerConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "LookerResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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
