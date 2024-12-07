package policyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

var resourceContextHandler = core.ContextHandler{
	ResourceName: resourceName,
	ResourceType: resourcetype.Resource,
	Create:       createPolicy,
	Read:         readPolicy,
	Update:       updatePolicy,
	Delete:       deletePolicy,
}

func PolicyTypes() []string {
	return []string{"POLICY_TYPE_GLOBAL", "global", "POLICY_TYPE_LOCAL", "local", "POLICY_TYPE_APPROVAL", "approval"}
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource allows management of various types of policies in the Cyral platform. Policies can be used to define access controls, data governance rules to ensure compliance and security within your database environment.",
		CreateContext: resourceContextHandler.CreateContext,
		ReadContext:   resourceContextHandler.ReadContext,
		UpdateContext: resourceContextHandler.UpdateContext,
		DeleteContext: resourceContextHandler.DeleteContext,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyV2StateContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Identifier for the policy, unique within the policy type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Indicates if the policy is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"scope": {
				Description: "Scope of the policy. If empty or omitted, all repositories are in scope.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_ids": {
							Description: "List of repository IDs that are in scope.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
					},
				},
			},
			"tags": {
				Description: "Tags associated with the policy to categorize it.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"valid_from": {
				Description: "Time when the policy comes into effect. If omitted, the policy is in effect immediately.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"valid_until": {
				Description: "Time after which the policy is no longer in effect. If omitted, the policy is in effect indefinitely.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"document": {
				Description: "The actual policy document in JSON format. It must conform to the schema for the policy type.",
				Type:        schema.TypeString,
				Required:    true,
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
				Optional:    true,
			},
			"type": {
				Description:  "Type of the policy, one of [`local`, `global`]",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(append(PolicyTypes(), ""), false),
			},
		},
	}
}
func importPolicyV2StateContext(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
	if err != nil {
		return nil, err
	}
	policyType := ids[0]
	policyID := ids[1]
	_ = d.Set("type", policyType)
	d.SetId(policyID)
	return []*schema.ResourceData{d}, nil
}
