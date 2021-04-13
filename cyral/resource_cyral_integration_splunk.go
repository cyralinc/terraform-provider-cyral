package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SplunkIntegration struct {
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
	Port        int    `json:"hecPort,string"`
	Host        string `json:"host"`
	Index       string `json:"index"`
	UseTLS      bool   `json:"useTLS"`
}

func (data SplunkIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("access_token", data.AccessToken)
	d.Set("port", data.Port)
	d.Set("host", data.Host)
	d.Set("index", data.Index)
	d.Set("use_tls", data.UseTLS)
}

func (data *SplunkIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.AccessToken = d.Get("access_token").(string)
	data.Port = d.Get("port").(int)
	data.Host = d.Get("host").(string)
	data.Index = d.Get("index").(string)
	data.UseTLS = d.Get("use_tls").(bool)
}

var ReadSplunkConfig = ResourceOperationConfig{
	Name:       "SplunkResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &SplunkIntegration{},
}

func resourceIntegrationSplunk() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "SplunkResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/splunk", c.ControlPlane)
				},
				ResourceData: &SplunkIntegration{},
				ResponseData: &IDBasedResponse{},
			}, ReadSplunkConfig,
		),
		ReadContext: ReadResource(ReadSplunkConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "SplunkResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &SplunkIntegration{},
			}, ReadSplunkConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "SplunkResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"port": {
				Type:      schema.TypeInt,
				Required:  true,
			},
			"host": {
				Type:      schema.TypeString,
				Required:  true,
			},
			"index": {
				Type:      schema.TypeString,
				Required:  true,
			},
			"use_tls": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
