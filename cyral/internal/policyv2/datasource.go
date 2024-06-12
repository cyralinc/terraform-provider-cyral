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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique ID for the policy. This field is automatically set.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the policy.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-friendly name for the policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description for the policy.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the policy is enabled.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tags associated with the policy for categorization.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"scope": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
			"valid_from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the policy comes into effect. An unspecified value implies the policy has no start time.",
			},
			"valid_until": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time after which the policy is no longer in effect. An unspecified value implies the policy will always apply once it comes into effect.",
			},
			"document": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The actual policy document in JSON format. It must conform to the schema for the policy type.",
			},
			"last_updated": {
				Description: "Information about when and by whom the policy was last updated.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created": {
				Description: "Information about when and by whom the policy was created.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enforced": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the policy is enforced.",
			},
		},
	}
}
