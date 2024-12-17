package wizard

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
)

var dsContextHandler = core.ContextHandler{
	ResourceName: dataSourceName,
	ResourceType: resourcetype.DataSource,
	Read:         readPolicyWizards,
}

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "This data source provides information policy wizards",
		ReadContext: dsContextHandler.ReadContext,
		Schema: map[string]*schema.Schema{
			"wizard_id": {
				Description: "id of the policy wizard of interest.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"wizards": {
				Description: "Set of supported policy wizards.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Description: "Information about a policy wizard.",
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Identifier for the policy wizard, use as the value of wizard_id parameter in the policy set resource.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "Name of the policy wizard.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "Description of the policy wizard.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"tags": {
							Description: "Tags associated with the policy wizard.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"parameter_schema": {
							Description: "JSON schema for the policy wizard parameters.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
