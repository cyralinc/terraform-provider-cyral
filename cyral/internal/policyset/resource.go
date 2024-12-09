package policyset

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
)

var resourceContextHandler = core.HTTPContextHandler{
	ResourceName:                 resourceName,
	ResourceType:                 resourcetype.Resource,
	SchemaReaderFactory:          func() core.SchemaReader { return &PolicySet{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &PolicySet{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		baseURL := &url.URL{
			Scheme: "https",
			Host:   c.ControlPlane,
			Path:   apiPathPolicySet,
		}
		return baseURL.String()
	},
	ReadUpdateDeleteURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		baseURL := &url.URL{
			Scheme: "https",
			Host:   c.ControlPlane,
			Path:   fmt.Sprintf("%s/%s", apiPathPolicySet, d.Id()),
		}
		return baseURL.String()
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource allows management of policy sets in the Cyral platform.",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),
		Importer: &schema.ResourceImporter{
			StateContext: importPolicySetStateContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Identifier for the policy set.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"wizard_id": {
				Description: "The ID of the policy wizard used to create this policy set.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Name of the policy set.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the policy set.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tags": {
				Description: "Tags associated with the policy set.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"scope": {
				Description: "Scope of the policy set.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_ids": {
							Description: "List of repository IDs that are in scope. Empty list means all repositories are in scope.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
					},
				},
			},
			"wizard_parameters": {
				Description: "Parameters passed to the wizard while creating the policy set.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Indicates if the policy set is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
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

func importPolicySetStateContext(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
