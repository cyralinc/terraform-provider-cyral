package deprecated

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LookerIntegration struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	URL          string `json:"url"`
}

func (data LookerIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("client_secret", data.ClientSecret)
	d.Set("client_id", data.ClientId)
	d.Set("url", data.URL)
	return nil
}

func (data *LookerIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.ClientSecret = d.Get("client_secret").(string)
	data.ClientId = d.Get("client_id").(string)
	data.URL = d.Get("url").(string)
	return nil
}

func ResourceIntegrationLooker() *schema.Resource {
	contextHandler := core.DefaultContextHandler{
		ResourceName:                 "Looker Integration",
		ResourceType:                 resourcetype.Resource,
		SchemaReaderFactory:          func() core.SchemaReader { return &LookerIntegration{} },
		SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &LookerIntegration{} },
		PostURLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/looker", c.ControlPlane)
		},
	}
	return &schema.Resource{
		DeprecationMessage: "Integration no longer supported.",
		Description:        "Manages integration with Looker.",
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
			"client_id": {
				Description: "Looker client id.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"client_secret": {
				Description: "Looker client secret.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"url": {
				Description: "Looker integration url.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
