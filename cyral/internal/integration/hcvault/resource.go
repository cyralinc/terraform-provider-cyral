package hcvault

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages integration with Hashicorp Vault to store secrets.",
		CreateContext: contextHandler.CreateContext(),
		ReadContext:   contextHandler.ReadContext(),
		UpdateContext: contextHandler.UpdateContext(),
		DeleteContext: contextHandler.DeleteContext(),
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
