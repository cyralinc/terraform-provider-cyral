package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateLogstashIntegrationResponse struct {
	ID string `json:"ID"`
}

func (response CreateLogstashIntegrationResponse) WriteResourceData(d *schema.ResourceData) {
	d.SetId(response.ID)
}

func (response *CreateLogstashIntegrationResponse) ReadResourceData(d *schema.ResourceData) {
	response.ID = d.Id()
}

type LogstashIntegrationData struct {
	Endpoint                   string `json:"endpoint"`
	Name                       string `json:"name"`
	UseMutualAuthentication    bool   `json:"useMutualAuthentication"`
	UsePrivateCertificateChain bool   `json:"usePrivateCertificateChain"`
	UseTLS                     bool   `json:"useTLS"`
}

func (data LogstashIntegrationData) WriteResourceData(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("endpoint", data.Endpoint)
	d.Set("use_mutual_authentication", data.UseMutualAuthentication)
	d.Set("use_private_certificate_chain", data.UsePrivateCertificateChain)
	d.Set("use_tls", data.UseTLS)
}

func (data *LogstashIntegrationData) ReadResourceData(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.Endpoint = d.Get("endpoint").(string)
	data.UseMutualAuthentication = d.Get("use_mutual_authentication").(bool)
	data.UsePrivateCertificateChain = d.Get("use_private_certificate_chain").(bool)
	data.UseTLS = d.Get("use_tls").(bool)
}

var ReadLogstashFunctionConfig = FunctionConfig{
	Name:       "LogstashResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/logstash/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &LogstashIntegrationData{},
}

func resourceIntegrationLogstash() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			FunctionConfig{
				Name:       "LogstashResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/logstash", c.ControlPlane)
				},
				ResourceData: &LogstashIntegrationData{},
				ResponseData: &CreateLogstashIntegrationResponse{},
			}, ReadLogstashFunctionConfig,
		),
		ReadContext: ReadResource(ReadLogstashFunctionConfig),
		UpdateContext: UpdateResource(
			FunctionConfig{
				Name:       "LogstashResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/logstash/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &LogstashIntegrationData{},
			}, ReadLogstashFunctionConfig,
		),
		DeleteContext: DeleteResource(
			FunctionConfig{
				Name:       "LogstashResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/logstash/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"use_mutual_authentication": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"use_private_certificate_chain": {
				Type:     schema.TypeBool,
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
