package policyv2

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
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &PolicyV2{} },
	ReadUpdateDeleteURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/%s/%s", c.ControlPlane, getAPIPath(d.Get("type").(string)), d.Get("id").(string))
	},
}

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "This data source provides information about a policy.",
		ReadContext: dsContextHandler.ReadContext(),
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Identifier for the policy, unique within the policy type.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "Type of the policy, one of [`local`, `global`, `approval`]",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Name of the policy.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the policy.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Indicates if the policy is enabled.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"tags": {
				Description: "Tags associated with the policy for categorization.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"scope": {
				Description: "Scope of the policy. If empty or omitted, all repositories are in scope.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_ids": {
							Description: "List of repository IDs that are in scope.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
					},
				},
			},
			"valid_from": {
				Description: "Time when the policy comes into effect. If omitted, the policy is in effect immediately.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"valid_until": {
				Description: "Time after which the policy is no longer in effect. If omitted, the policy is in effect indefinitely.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"document": {
				Description: "The actual policy document in JSON format. It must conform to the schema for the policy type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated": {
				Description: "Information about when and by whom the policy was last updated.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created": {
				Description: "Information about when and by whom the policy was created.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enforced": {
				Description: "Indicates if the policy is enforced. If not enforced, no action is taken based on the policy, but alerts are triggered for violations.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}
