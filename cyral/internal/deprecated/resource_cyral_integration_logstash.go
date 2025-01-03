package deprecated

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LogstashIntegration struct {
	Endpoint                   string `json:"endpoint"`
	Name                       string `json:"name"`
	UseMutualAuthentication    bool   `json:"useMutualAuthentication"`
	UsePrivateCertificateChain bool   `json:"usePrivateCertificateChain"`
	UseTLS                     bool   `json:"useTLS"`
}

func (data LogstashIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("name", data.Name)
	d.Set("endpoint", data.Endpoint)
	d.Set("use_mutual_authentication", data.UseMutualAuthentication)
	d.Set("use_private_certificate_chain", data.UsePrivateCertificateChain)
	d.Set("use_tls", data.UseTLS)
	return nil
}

func (data *LogstashIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.Name = d.Get("name").(string)
	data.Endpoint = d.Get("endpoint").(string)
	data.UseMutualAuthentication = d.Get("use_mutual_authentication").(bool)
	data.UsePrivateCertificateChain = d.Get("use_private_certificate_chain").(bool)
	data.UseTLS = d.Get("use_tls").(bool)
	return nil
}

func ResourceIntegrationLogstash() *schema.Resource {
	contextHandler := core.HTTPContextHandler{
		ResourceName:                 "Logstash Integration",
		ResourceType:                 resourcetype.Resource,
		SchemaReaderFactory:          func() core.SchemaReader { return &LogstashIntegration{} },
		SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &LogstashIntegration{} },
		BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/logstash", c.ControlPlane)
		},
	}
	return &schema.Resource{
		DeprecationMessage: "Use resource `cyral_integration_logging` instead.",
		Description:        "Manages integration with Logstash.",
		CreateContext:      contextHandler.CreateContext(),
		ReadContext:        contextHandler.ReadContext(),
		UpdateContext:      contextHandler.UpdateContext(),
		DeleteContext:      contextHandler.DeleteContext(),
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
			"endpoint": {
				Description: "The endpoint used to connect to Logstash.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"use_mutual_authentication": {
				Description: "Logstash configured to use mutual authentication.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"use_private_certificate_chain": {
				Description: "Logstash configured to use private certificate chain.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"use_tls": {
				Description: "Logstash configured to use mutual TLS.",
				Type:        schema.TypeBool,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
