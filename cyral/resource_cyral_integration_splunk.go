package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SplunkIntegration struct {
	Name                     string `json:"name"`
	AccessToken              string `json:"accessToken"`
	Port                     int    `json:"hecPort,string"`
	Host                     string `json:"host"`
	Index                    string `json:"index"`
	UseTLS                   bool   `json:"useTLS"`
	CyralActivityLogsEnabled bool   `json:"cyralActivityLogsEnabled"`
}

func (data SplunkIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("access_token", data.AccessToken)
	d.Set("port", data.Port)
	d.Set("host", data.Host)
	d.Set("index", data.Index)
	d.Set("use_tls", data.UseTLS)
	// d.Set("cyral_activity_logs_enabled", data.CyralActivityLogsEnabled)
}

func (data *SplunkIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.AccessToken = d.Get("access_token").(string)
	data.Port = d.Get("port").(int)
	data.Host = d.Get("host").(string)
	data.Index = d.Get("index").(string)
	data.UseTLS = d.Get("use_tls").(bool)
	// data.CyralActivityLogsEnabled = d.Get("cyral_activity_logs_enabled").(bool)
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
		Description: "Manages [integration with Splunk](https://cyral.com/docs/integrations/siem/splunk/#procedure).",
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
			"access_token": {
				Description: "Splunk access token.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"port": {
				Description: "Splunk host port.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"host": {
				Description: "Splunk host.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"index": {
				Description: "Splunk data index name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"use_tls": {
				Description: "Should the communication with Splunk use TLS encryption?",
				Type:        schema.TypeBool,
				Required:    true,
			},
			/* "cyral_activity_logs_enabled": {
				Description: "Should enable Cyral activity logs.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			}, */
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
