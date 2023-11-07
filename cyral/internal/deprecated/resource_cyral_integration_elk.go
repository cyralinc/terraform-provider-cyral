package deprecated

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ELKIntegration struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	KibanaURL string `json:"kibanaUrl"`
	ESURL     string `json:"esUrl"`
}

func (data ELKIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("name", data.Name)
	d.Set("kibana_url", data.KibanaURL)
	d.Set("es_url", data.ESURL)
	return nil
}

func (data *ELKIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.ID = d.Id()
	data.Name = d.Get("name").(string)
	data.KibanaURL = d.Get("kibana_url").(string)
	data.ESURL = d.Get("es_url").(string)
	return nil
}

var ReadELKConfig = core.ResourceOperationConfig{
	Name:       "ELKResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())
	},
	NewResponseData:     func(_ *schema.ResourceData) core.SchemaWriter { return &ELKIntegration{} },
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Integration elk"},
}

func ResourceIntegrationELK() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use resource `cyral_integration_logging` instead.",
		Description:        "Manages [integration with ELK](https://cyral.com/docs/integrations/siem/elk/) to push sidecar metrics.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "ELKResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/elk", c.ControlPlane)
				},
				NewResourceData: func() core.SchemaReader { return &ELKIntegration{} },
			}, ReadELKConfig,
		),
		ReadContext: core.ReadResource(ReadELKConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "ELKResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() core.SchemaReader { return &ELKIntegration{} },
			}, ReadELKConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "ELKResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())
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
			"kibana_url": {
				Description: "Kibana URL.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"es_url": {
				Description: "Elastic Search URL.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
