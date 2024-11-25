package policywizardv1

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
)

var dsContextHandler = core.DefaultContextHandler{
	ResourceName:                 dataSourceName,
	ResourceType:                 resourcetype.DataSource,
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &PolicySet{} },
	ReadUpdateDeleteURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/%s/%s", c.ControlPlane, apiPathPolicySet, d.Get("id").(string))
	},
}

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "This data source provides information about a policy set.",
		ReadContext: dsContextHandler.ReadContext(),
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Identifier for the policy set.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"wizard_id": {
				Description: "The ID of the policy wizard used to create this policy set.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of the policy set.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the policy set.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tags": {
				Description: "Tags associated with the policy set.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"scope": {
				Description: "Scope of the policy set.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_ids": {
							Description: "List of repository IDs that are in scope. Empty list means all repositories are in scope.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
					},
				},
			},
			"wizard_parameters": {
				Description: "Parameters passed to the wizard while creating the policy set.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Indicates if the policy set is enabled.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"policies": {
				Description: "List of policies that comprise the policy set.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Type of the policy.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "ID of the policy.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"last_updated": {
				Description: "Information about when and by whom the policy set was last updated.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"created": {
				Description: "Information about when and by whom the policy set was created.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
