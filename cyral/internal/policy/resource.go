package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:                 resourceName,
	ResourceType:                 resourcetype.Resource,
	SchemaReaderFactory:          func() core.SchemaReader { return &Policy{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &Policy{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/policies", c.ControlPlane)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [policies](https://cyral.com/docs/reference/policy). See also: " +
			"[Policy Rule](./policy_rule.md). For more information, see the " +
			"[Policy Guide](https://cyral.com/docs/policy/overview).",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext: resourceContextHandler.ReadContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "policy not found",
			OperationType:  operationtype.Read,
		}),
		UpdateContext: resourceContextHandler.UpdateContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "policy not found",
			OperationType:  operationtype.Update,
		}, nil),
		DeleteContext: resourceContextHandler.DeleteContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "policy not found",
			OperationType:  operationtype.Delete,
		}),
		Schema: map[string]*schema.Schema{
			"created": {
				Description: "Timestamp for the policy creation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"data": {
				Description: "List that specify which data fields a policy manages. Each field is represented by the LABEL " +
					"you established for it in your data map. The actual location of that data (the names of fields, columns, " +
					"or databases that hold it) is listed in the data map.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data_label_tags": {
				Description: "List of tags that represent sets of data labels (established in your data map) that " +
					"are used to specify the collections of data labels that the policy manages. For more information, " +
					"see [The tags block of a policy](https://cyral.com/docs/policy/policy-structure#the-tags-block-of-a-policy)",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Description: "String that describes the policy (ex: `your_policy_description`).",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"enabled": {
				Description: "Boolean that causes a policy to be enabled or disabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"last_updated": {
				Description: "Timestamp for the last update performed in this policy.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Policy name that will be used internally in Control Plane (ex: `your_policy_name`).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tags": {
				Deprecated: "Use `metadata_tags` instead. This will be removed in the next major version of the provider.",
				Description: "Metadata tags that can be used to organize and/or classify your policies " +
					"(ex: `[your_tag1, your_tag2]`).",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"metadata_tags"},
			},
			"metadata_tags": {
				Description: "Metadata tags that can be used to organize and/or classify your policies " +
					"(ex: `[your_tag1, your_tag2]`).",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"tags"},
			},
			"type": {
				Description: "Policy type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"version": {
				Description: "Incremental counter for every update on the policy.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},

		CustomizeDiff: func(ctx context.Context, resourceDiff *schema.ResourceDiff, i interface{}) error {
			computedKeysToChange := []string{"last_updated", "version"}
			utils.SetKeysAsNewComputedIfPlanHasChanges(resourceDiff, computedKeysToChange)
			return nil
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func ListPolicies(c *client.Client) ([]Policy, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Init ListPolicies")

	url := fmt.Sprintf("https://%s/v1/policies", c.ControlPlane)
	resp, err := c.DoRequest(ctx, url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var listResp PolicyListResponse
	if err := json.Unmarshal(resp, &listResp); err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", listResp))

	var policies []Policy
	for _, policyID := range listResp.Policies {
		url := fmt.Sprintf("https://%s/v1/policies/%s",
			c.ControlPlane, policyID)
		resp, err := c.DoRequest(ctx, url, http.MethodGet, nil)
		if err != nil {
			return nil, err
		}

		var policy Policy
		if err := json.Unmarshal(resp, &policy); err != nil {
			return nil, err
		}
		tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", policy))

		policies = append(policies, policy)
	}

	tflog.Debug(ctx, "End ListPolicies")
	return policies, nil
}
