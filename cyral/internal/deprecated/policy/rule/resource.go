package rule

import (
	"context"
	"fmt"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceContextHandler = core.HTTPContextHandler{
	ResourceName:                  resourceName,
	ResourceType:                  resourcetype.Resource,
	SchemaReaderFactory:           func() core.SchemaReader { return &PolicyRule{} },
	SchemaWriterFactoryGetMethod:  func(_ *schema.ResourceData) core.SchemaWriter { return &PolicyRule{} },
	SchemaWriterFactoryPostMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &PolicyRuleIDBasedResponse{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/policies/%s/rules", c.ControlPlane, d.Get("policy_id").(string))
	},
	ReadUpdateDeleteURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		policyID, policyRuleID := unmarshalPolicyRuleID(d)
		return fmt.Sprintf("https://%s/v1/policies/%s/rules/%s",
			c.ControlPlane, policyID, policyRuleID)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "For control planes `>= v4.15`, use resource `cyral_policy_v2` instead.",
		Description: "Manages [policy rules](https://cyral.com/docs/policy/#rules). " +
			"See also the [`cyral_policy`](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/resources/policy) " +
			"resource and the [Policy Guide](https://cyral.com/docs/policy#the-rules-block-of-a-policy)." +
			"\n\n-> 1. Unless you create a default rule, users and groups only have the rights you explicitly grant them." +
			"<br> 2. Each contexted rule comprises these fields: `data`, `rows`, `severity` `additional_checks`, `dataset_rewrites`. The only required fields are `data` and `rows`." +
			"<br> 3. The rules block does not need to include all three operation types (reads, updates and deletes); actions you omit are disallowed." +
			"<br>4. If you do not include a hosts block, Cyral does not enforce limits based on the connecting client's host address.",

		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext: resourceContextHandler.ReadContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "not found in policy",
			OperationType:  operationtype.Read,
		}),
		UpdateContext: resourceContextHandler.UpdateContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "not found in policy",
			OperationType:  operationtype.Update,
		}, nil),
		DeleteContext: resourceContextHandler.DeleteContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "not found in policy",
			OperationType:  operationtype.Delete,
		}),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: policyRuleResourceSchemaV0().
					CoreConfigSchema().ImpliedType(),
				Upgrade: UpgradePolicyRuleV0,
			},
		},

		Schema: policyRuleResourceSchemaV0().Schema,

		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, err
				}
				policyID := ids[0]
				d.Set("policy_id", policyID)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func ruleSchema(description string) *schema.Schema {
	return &schema.Schema{
		Description: description,
		Type:        schema.TypeList,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"additional_checks": {
					Description: "Constraints on the data access specified in " +
						"[Rego](https://www.openpolicyagent.org/docs/latest/policy-language/). " +
						"See [Additional checks](https://cyral.com/docs/policy/rules/#additional-checks).",
					Type:     schema.TypeString,
					Optional: true,
				},
				"data": {
					Type: schema.TypeList,
					Description: "The data locations protected by this rule. " +
						"Use `*` if you want to define `any` data location. " +
						"For more information, see the " +
						"[policy rules](https://cyral.com/docs/policy/rules#contexted-rules) documentation.",
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"dataset_rewrites": {
					Description: "Defines how requests should be rewritten in the case of " +
						"policy violations. See [Request rewriting](https://cyral.com/docs/policy/rules/#request-rewriting).",
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"dataset": {
								Description: "The dataset that should be rewritten." +
									"In the case of Snowflake, this denotes a fully qualified table name in the form: " +
									"`<database>.<schema>.<table>`",
								Type:     schema.TypeString,
								Required: true,
							},
							"repo": {
								Description: "The name of the repository that the rewrite applies to.",
								Type:        schema.TypeString,
								Required:    true,
							},
							"substitution": {
								Description: "The request used to substitute references to the dataset.",
								Type:        schema.TypeString,
								Required:    true,
							},
							"parameters": {
								Description: "The set of parameters used in the substitution request, " +
									"these are references to fields in the activity log as described in " +
									"the [Additional Checks section](https://cyral.com/docs/policy/rules/#additional-checks).",
								Type:     schema.TypeList,
								Required: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"rows": {
					Description: "The number of records (for example, rows or documents) that can " +
						"be accessed/affected in a single statement. Use positive integer " +
						"numbers to define how many records. If you want to define `any` " +
						"number of records, set to `-1`.",
					Type:     schema.TypeInt,
					Required: true,
				},
				"severity": {
					Description: "severity level that's recorded when someone violate this rule. " +
						"This is an informational value. Settings: (`low` | `medium` | `high`). " +
						"If not specified, the severity is considered to be low.",
					Type:     schema.TypeString,
					Optional: true,
					Default:  "low",
				},
				"rate_limit": {
					Description: "Rate Limit specifies the limit of calls that a user can make within a given time period.",
					Type:        schema.TypeInt,
					Optional:    true,
				},
			},
		},
	}
}

func policyRuleResourceSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Description: "The ID of the policy you are adding this rule to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"deletes": ruleSchema("A contexted rule for accesses of the type `delete`."),
			"reads":   ruleSchema("A contexted rule for accesses of the type `read`."),
			"updates": ruleSchema("A contexted rule for accesses of the type `update`."),
			"hosts": {
				Description: "Hosts specification that limits access to only those users connecting from a certain network location.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"identities": {
				Description: "Identities specifies the people, applications, " +
					"or groups this rule applies to. Every rule except your default rule has one. " +
					"It can have 4 fields: `db_roles`, `groups`, `users` and `services`.",
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_roles": {
							Description: "Database roles that this rule will apply to.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"groups": {
							Description: "Groups that this rule will apply to.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"services": {
							Description: "Services that this rule will apply to.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"users": {
							Description: "Users that this rule will apply to.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			// Computed arguments
			"policy_rule_id": {
				Description: "The ID of the policy rule.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// Previously, the ID of `cyral_policy_rule` was simply the policy rule ID. That
// is not ideal, because to realize the resource we also need the policy ID
// itself. That creates an inconsistency between the ID syntax used in
// `terraform import` and the computed id for the resource. The goal of this
// upgrade is to set the `id` attribute to have the format
// `{policy_id}/{policy_rule_id}`.
func UpgradePolicyRuleV0(
	_ context.Context,
	rawState map[string]interface{},
	_ interface{},
) (map[string]interface{}, error) {
	policyRuleID := rawState["id"].(string)

	// Make sure there is not `/` in the previous ID, otherwise the new
	// resource state ID will become inconsistent.
	idComponents := strings.Split(policyRuleID, "/")
	if len(idComponents) > 1 {
		return rawState, fmt.Errorf("unexpected format for resource ID: " +
			`found more than one field separated by "/"`)
	}

	policyID := rawState["policy_id"].(string)
	newID := utils.MarshalComposedID([]string{policyID, policyRuleID}, "/")
	rawState["id"] = newID

	return rawState, nil
}
