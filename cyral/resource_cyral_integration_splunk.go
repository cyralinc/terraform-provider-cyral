package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateSplunkIntegrationResponse struct {
	ID string `json:"id"`
}

func (response CreateSplunkIntegrationResponse) WriteResourceData(d *schema.ResourceData) {
	d.SetId(response.ID)
}

func (response *CreateSplunkIntegrationResponse) ReadResourceData(d *schema.ResourceData) {
	response.ID = d.Id()
}

type SplunkIntegrationData struct {
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
	Port        int    `json:"hecPort,string"`
	Host        string `json:"host"`
	Index       string `json:"index"`
	UseTLS      bool   `json:"useTLS"`
}

func (data SplunkIntegrationData) WriteResourceData(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("access_token", data.AccessToken)
	d.Set("port", data.Port)
	d.Set("host", data.Host)
	d.Set("index", data.Index)
	d.Set("use_tls", data.UseTLS)

}

func (data *SplunkIntegrationData) ReadResourceData(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.AccessToken = d.Get("access_token").(string)
	data.Port = d.Get("port").(int)
	data.Host = d.Get("host").(string)
	data.Index = d.Get("index").(string)
	data.UseTLS = d.Get("use_tls").(bool)
}

var ReadSplunkFunctionConfig = FunctionConfig{
	Name:       "SplunkResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &SplunkIntegrationData{},
}

func resourceIntegrationSplunk() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			FunctionConfig{
				Name:       "SplunkResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/splunk", c.ControlPlane)
				},
				ResourceData: &SplunkIntegrationData{},
				ResponseData: &CreateSplunkIntegrationResponse{},
			}, ReadSplunkFunctionConfig,
		),
		ReadContext: ReadResource(ReadSplunkFunctionConfig),
		UpdateContext: UpdateResource(
			FunctionConfig{
				Name:       "SplunkResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &SplunkIntegrationData{},
			}, ReadSplunkFunctionConfig,
		),
		DeleteContext: DeleteResource(
			FunctionConfig{
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
				Type:     schema.TypeInt,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"index": {
				Type:     schema.TypeString,
				Required: true,
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
