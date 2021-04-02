package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateELKIntegrationResponse struct {
	ID string `json:"ID"`
}

func (response CreateELKIntegrationResponse) WriteResourceData(d *schema.ResourceData) {
	d.SetId(response.ID)
}

func (response *CreateELKIntegrationResponse) ReadResourceData(d *schema.ResourceData) {
	response.ID = d.Id()
}

type ELKIntegrationData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	KibanaURL string `json:"kibanaUrl"`
	ESURL     string `json:"esUrl"`
}

func (data ELKIntegrationData) WriteResourceData(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("kibana_url", data.KibanaURL)
	d.Set("es_url", data.ESURL)
}

func (data *ELKIntegrationData) ReadResourceData(d *schema.ResourceData) {
	data.ID = d.Id()
	data.Name = d.Get("name").(string)
	data.KibanaURL = d.Get("kibana_url").(string)
	data.ESURL = d.Get("es_url").(string)
}

var CreateELKFunctionConfig = FunctionConfig{
	Name:       "ELKResourceCreate",
	HttpMethod: http.MethodPost,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/elk", c.ControlPlane)
	},
	ResourceData:       &ELKIntegrationData{},
	ResponseData:       &CreateELKIntegrationResponse{},
	ReadFunctionConfig: &ReadELKFunctionConfig,
}

var ReadELKFunctionConfig = FunctionConfig{
	Name:       "ELKResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &ELKIntegrationData{},
}

var UpdateELKFunctionConfig = FunctionConfig{
	Name:       "ELKResourceUpdate",
	HttpMethod: http.MethodPut,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())
	},
	ResourceData:       &ELKIntegrationData{},
	ReadFunctionConfig: &ReadELKFunctionConfig,
}

var DeleteELKFunctionConfig = FunctionConfig{
	Name:       "ELKResourceDelete",
	HttpMethod: http.MethodDelete,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())
	},
}

func resourceIntegrationELK() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateELKFunctionConfig.Create,
		ReadContext:   ReadELKFunctionConfig.Read,
		UpdateContext: UpdateELKFunctionConfig.Update,
		DeleteContext: DeleteELKFunctionConfig.Delete,

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
