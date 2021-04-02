package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateLookerIntegrationResponse struct {
	ID string `json:"ID"`
}

func (response CreateLookerIntegrationResponse) WriteResourceData(d *schema.ResourceData) {
	d.SetId(response.ID)
}

func (response *CreateLookerIntegrationResponse) ReadResourceData(d *schema.ResourceData) {
	response.ID = d.Id()
}

type LookerIntegrationData struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	URL          string `json:"url"`
}

func (data LookerIntegrationData) WriteResourceData(d *schema.ResourceData) {
	d.Set("client_secret", data.ClientSecret)
	d.Set("client_id", data.ClientId)
	d.Set("url", data.URL)
}

func (data *LookerIntegrationData) ReadResourceData(d *schema.ResourceData) {
	data.ClientSecret = d.Get("client_secret").(string)
	data.ClientId = d.Get("client_id").(string)
	data.URL = d.Get("url").(string)
}

var CreateLookerFunctionConfig = FunctionConfig{
	Name:       "LookerResourceCreate",
	HttpMethod: http.MethodPost,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/looker", c.ControlPlane)
	},
	ResourceData:       &LookerIntegrationData{},
	ResponseData:       &CreateLookerIntegrationResponse{},
	ReadFunctionConfig: &ReadLookerFunctionConfig,
}

var ReadLookerFunctionConfig = FunctionConfig{
	Name:       "LookerResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &LookerIntegrationData{},
}

var UpdateLookerFunctionConfig = FunctionConfig{
	Name:       "LookerResourceUpdate",
	HttpMethod: http.MethodPut,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())
	},
	ResourceData:       &LookerIntegrationData{},
	ReadFunctionConfig: &ReadLookerFunctionConfig,
}

var DeleteLookerFunctionConfig = FunctionConfig{
	Name:       "LookerResourceDelete",
	HttpMethod: http.MethodDelete,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())
	},
}

func resourceIntegrationLooker() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateLookerFunctionConfig.Create,
		ReadContext:   ReadLookerFunctionConfig.Read,
		UpdateContext: UpdateLookerFunctionConfig.Update,
		DeleteContext: DeleteLookerFunctionConfig.Delete,

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
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
