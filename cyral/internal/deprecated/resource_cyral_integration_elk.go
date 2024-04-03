package deprecated

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
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

func ResourceIntegrationELK() *schema.Resource {
	contextHandler := core.DefaultContextHandler{
		ResourceName:                 "ELK Integration",
		ResourceType:                 resourcetype.Resource,
		SchemaReaderFactory:          func() core.SchemaReader { return &ELKIntegration{} },
		SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &ELKIntegration{} },
		PostURLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/elk", c.ControlPlane)
		},
	}
	return &schema.Resource{
		DeprecationMessage: "Use resource `cyral_integration_logging` instead.",
		Description:        "Manages [integration with ELK](https://cyral.com/docs/integrations/siem/elk/) to push sidecar metrics.",
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
