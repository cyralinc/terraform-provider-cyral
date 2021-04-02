package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ELKIntegration struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	KibanaURL string `json:"kibanaUrl"`
	ESURL     string `json:"esUrl"`
}

func (data ELKIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("kibana_url", data.KibanaURL)
	d.Set("es_url", data.ESURL)
}

func (data *ELKIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.ID = d.Id()
	data.Name = d.Get("name").(string)
	data.KibanaURL = d.Get("kibana_url").(string)
	data.ESURL = d.Get("es_url").(string)
}

var ReadELKConfig = ResourceOperationConfig{
	Name:       "ELKResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &ELKIntegration{},
}

func resourceIntegrationELK() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "ELKResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/elk", c.ControlPlane)
				},
				ResourceData: &ELKIntegration{},
				ResponseData: &IDBasedResponse{},
			}, ReadELKConfig,
		),
		ReadContext: ReadResource(ReadELKConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "ELKResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &ELKIntegration{},
			}, ReadELKConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "ELKResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kibana_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"es_url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
