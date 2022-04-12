package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LogstashIntegration struct {
	Endpoint                   string `json:"endpoint"`
	Name                       string `json:"name"`
	UseMutualAuthentication    bool   `json:"useMutualAuthentication"`
	UsePrivateCertificateChain bool   `json:"usePrivateCertificateChain"`
	UseTLS                     bool   `json:"useTLS"`
}

func (data LogstashIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("endpoint", data.Endpoint)
	d.Set("use_mutual_authentication", data.UseMutualAuthentication)
	d.Set("use_private_certificate_chain", data.UsePrivateCertificateChain)
	d.Set("use_tls", data.UseTLS)
}

func (data *LogstashIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.Endpoint = d.Get("endpoint").(string)
	data.UseMutualAuthentication = d.Get("use_mutual_authentication").(bool)
	data.UsePrivateCertificateChain = d.Get("use_private_certificate_chain").(bool)
	data.UseTLS = d.Get("use_tls").(bool)
}

var ReadLogstashConfig = ResourceOperationConfig{
	Name:       "LogstashResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/logstash/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &LogstashIntegration{},
}

func resourceIntegrationLogstash() *schema.Resource {
	return &schema.Resource{
		Description: "Manages integration with Logstash.",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "LogstashResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/logstash", c.ControlPlane)
				},
				ResourceData: &LogstashIntegration{},
				ResponseData: &IDBasedResponse{},
			}, ReadLogstashConfig,
		),
		ReadContext: ReadResource(ReadLogstashConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "LogstashResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/logstash/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &LogstashIntegration{},
			}, ReadLogstashConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "LogstashResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/logstash/%s", c.ControlPlane, d.Id())
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
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Description: "The endpoint used to connect to Logstash.",
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"use_mutual_authentication": {
				Description: "Logstash configured to use mutual authentication.",
				Type:     schema.TypeBool,
				Required: true,
			},
			"use_private_certificate_chain": {
				Description: "Logstash configured to use private certificate chain.",
				Type:     schema.TypeBool,
				Required: true,
			},
			"use_tls": {
				Description: "Logstash configured to use mutual TLS.",
				Type:     schema.TypeBool,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
