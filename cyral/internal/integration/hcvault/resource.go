package hcvault

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:        "HC Vault Integration",
	ResourceType:        resourcetype.Resource,
	SchemaReaderFactory: func() core.SchemaReader { return &HCVaultIntegration{} },
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &HCVaultIntegration{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/secretProviders/hcvault", c.ControlPlane)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages integration with Hashicorp Vault to store secrets.",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"auth_method": {
				Description: "Authentication method for the integration.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"auth_type": {
				Description: "Authentication type for the integration.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"name": {
				Description: "Integration name that will be used internally in the control plane.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"server": {
				Description: "Server on which the vault service is running.",
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
