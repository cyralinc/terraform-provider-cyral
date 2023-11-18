package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [policies](https://cyral.com/docs/reference/policy). See also: " +
			"[Policy Rule](./policy_rule.md). For more information, see the " +
			"[Policy Guide](https://cyral.com/docs/policy/overview).",
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
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

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourcePolicyCreate")
	c := m.(*client.Client)

	d.Set("type", "terraform")
	policy := getPolicyInfoFromResource(d)

	url := fmt.Sprintf("https://%s/v1/policies", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, policy)
	if err != nil {
		return utils.CreateError("Unable to create policy", fmt.Sprintf("%v", err))
	}

	response := core.IDBasedResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	d.SetId(response.ID)

	tflog.Debug(ctx, "End resourcePolicyCreate")

	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourcePolicyRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return utils.CreateError("Unable to read policy", fmt.Sprintf("%v", err))
	}

	response := Policy{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	d.Set("created", response.Meta.Created.String())
	d.Set("data", response.Data)
	d.Set("data_label_tags", response.Tags)
	d.Set("description", response.Meta.Description)
	d.Set("enabled", response.Meta.Enabled)
	d.Set("last_updated", response.Meta.LastUpdated.String())
	d.Set("name", response.Meta.Name)
	d.Set("type", response.Meta.Type)
	d.Set("version", response.Meta.Version)
	// Once the `tags` field is removed, this conditional logic should also be
	// removed and only the `metadata_tags` should be set.
	_, isDeprecatedFieldSet := d.GetOk("tags")
	if isDeprecatedFieldSet {
		d.Set("tags", response.Meta.Tags)
	} else {
		d.Set("metadata_tags", response.Meta.Tags)
	}

	tflog.Debug(ctx, "End resourcePolicyRead")
	return diag.Diagnostics{}
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourcePolicyUpdate")
	c := m.(*client.Client)

	d.Set("type", "terraform")
	policy := getPolicyInfoFromResource(d)

	url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

	_, err := c.DoRequest(url, http.MethodPut, policy)
	if err != nil {
		return utils.CreateError("Unable to update policy", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, "End resourcePolicyUpdate")

	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourcePolicyDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return utils.CreateError("Unable to delete policy", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, "End resourcePolicyDelete")

	return diag.Diagnostics{}
}

func getPolicyInfoFromResource(d *schema.ResourceData) Policy {
	data := utils.GetStrListFromSchemaField(d, "data")
	dataTags := utils.GetStrListFromSchemaField(d, "data_label_tags")
	metadataTags := utils.GetStrListFromSchemaField(d, "metadata_tags")
	if len(metadataTags) == 0 {
		metadataTags = utils.GetStrListFromSchemaField(d, "tags")
	}

	policy := Policy{
		Data: data,
		Tags: dataTags,
		Meta: &PolicyMetadata{
			Tags: metadataTags,
		},
	}

	if v, ok := d.Get("name").(string); ok {
		policy.Meta.Name = v
	}

	if v, ok := d.Get("version").(string); ok {
		policy.Meta.Version = v
	}

	if v, ok := d.Get("type").(string); ok {
		policy.Meta.Type = v
	}

	if v, ok := d.Get("enabled").(bool); ok {
		policy.Meta.Enabled = v
	}

	if v, ok := d.Get("description").(string); ok {
		policy.Meta.Description = v
	}

	return policy
}

func ListPolicies(c *client.Client) ([]Policy, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Init ListPolicies")

	url := fmt.Sprintf("https://%s/v1/policies", c.ControlPlane)
	resp, err := c.DoRequest(url, http.MethodGet, nil)
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
		resp, err := c.DoRequest(url, http.MethodGet, nil)
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
