package policyv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:                 resourceName,
	ResourceType:                 resourcetype.Resource,
	SchemaReaderFactory:          func() core.SchemaReader { return &PolicyV2{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &PolicyV2{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/%s", c.ControlPlane, getAPIPath(d.Get("type").(string)))
	},
	ReadUpdateDeleteURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/%s/%s",
			c.ControlPlane,
			getAPIPath(d.Get("type").(string)),
			d.Get("id").(string),
		)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "Some description.",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyV2StateContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "...",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "...",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Indicates if the policy is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"scope": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
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
				Description: "Time when the policy comes into effect.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"valid_until": {
				Description: "Time after which the policy is no longer in effect.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"document": {
				Description: "The actual policy document in JSON format.",
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
				Description: "Indicates if the policy is enforced.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"type": {
				Description: "Type of the policy.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validTypes := []string{"POLICY_TYPE_UNSPECIFIED", "POLICY_TYPE_GLOBAL", "global", "POLICY_TYPE_LOCAL", "local", "POLICY_TYPE_APPROVAL", "approval"}
					for _, validType := range validTypes {
						if v == validType {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validTypes, v))
					return
				},
			},
		},
	}
}
func importPolicyV2StateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
	if err != nil {
		return nil, err
	}
	policyType := ids[0]
	policyID := ids[1]
	d.Set("type", policyType)
	d.SetId(policyID)
	return []*schema.ResourceData{d}, nil
}
