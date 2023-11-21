package rule

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolicyRule struct {
	Deletes    []Rule    `json:"deletes,omitempty"`
	Hosts      []string  `json:"hosts,omitempty"`
	Identities *Identity `json:"identities,omitempty"`
	Reads      []Rule    `json:"reads,omitempty"`
	RuleID     string    `json:"ruleId"`
	Updates    []Rule    `json:"updates,omitempty"`
}

type Rule struct {
	AdditionalChecks string           `json:"additionalChecks"`
	Data             []string         `json:"data,omitempty"`
	DatasetRewrites  []DatasetRewrite `json:"datasetRewrites,omitempty"`
	Rows             int              `json:"rows"`
	Severity         string           `json:"severity"`
	RateLimit        int              `json:"rateLimit"`
}

type DatasetRewrite struct {
	Dataset      string   `json:"dataset"`
	Parameters   []string `json:"parameters,omitempty"`
	Repo         string   `json:"repo"`
	Substitution string   `json:"substitution"`
}

type Identity struct {
	DBRoles  []string `json:"dbRoles,omitempty"`
	Groups   []string `json:"groups,omitempty"`
	Services []string `json:"services,omitempty"`
	Users    []string `json:"users,omitempty"`
}

func unmarshalPolicyRuleID(d *schema.ResourceData) (policyID string, policyRuleID string) {
	// We must be careful when dealing with the ID. Specially in the Read
	// operation, due to state upgrade from v0 of this resource's schema to
	// v1. In v0, there exists only one field (the policy rule
	// ID). Therefore, if we assume there are two, the first `terraform
	// refresh` done when upgrading will fail.
	ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
	if err == nil {
		// This is the new way to organize the IDs (v1).
		policyID = ids[0]
		policyRuleID = ids[1]
	} else {
		// This conditional branch is here to treat legacy resources (v0).
		policyID = d.Get("policy_id").(string)
		policyRuleID = d.Id()
	}
	return
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

func ResourcePolicyRule() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [policy rules](https://cyral.com/docs/reference/policy/#rules). " +
			"See also the [`cyral_policy`](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/resources/policy) " +
			"resource and the [Policy Guide](https://cyral.com/docs/policy#the-rules-block-of-a-policy)." +
			"\n\n-> 1. Unless you create a default rule, users and groups only have the rights you explicitly grant them." +
			"<br> 2. Each contexted rule comprises these fields: `data`, `rows`, `severity` `additional_checks`, `dataset_rewrites`. The only required fields are `data` and `rows`." +
			"<br> 3. The rules block does not need to include all three operation types (reads, updates and deletes); actions you omit are disallowed." +
			"<br>4. If you do not include a hosts block, Cyral does not enforce limits based on the connecting client's host address.",

		CreateContext: resourcePolicyRuleCreate,
		ReadContext:   resourcePolicyRuleRead,
		UpdateContext: resourcePolicyRuleUpdate,
		DeleteContext: core.DeleteResource(policyRuleDeleteConfig()),

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

func resourcePolicyRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourcePolicyRuleCreate")
	c := m.(*client.Client)

	policyID := d.Get("policy_id").(string)
	resourceData := getPolicyRuleInfoFromResource(ctx, d)
	tflog.Debug(ctx, fmt.Sprintf("resourcePolicyRuleCreate - policyRule: %#v", resourceData))

	url := fmt.Sprintf("https://%s/v1/policies/%s/rules", c.ControlPlane, policyID)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return utils.CreateError("Unable to create policy rule", fmt.Sprintf("%v", err))
	}

	response := core.IDBasedResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	d.SetId(utils.MarshalComposedID([]string{
		policyID,
		response.ID},
		"/"))

	tflog.Debug(ctx, "End resourcePolicyRuleCreate")

	return resourcePolicyRuleRead(ctx, d, m)
}

func resourcePolicyRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourcePolicyRuleRead")
	c := m.(*client.Client)

	policyID, policyRuleID := unmarshalPolicyRuleID(d)
	url := fmt.Sprintf("https://%s/v1/policies/%s/rules/%s",
		c.ControlPlane, policyID, policyRuleID)

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return utils.CreateError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	response := PolicyRule{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	deletes := flattenRulesList(ctx, response.Deletes)
	tflog.Debug(ctx, fmt.Sprintf("flattened deletes %#v", deletes))
	if err := d.Set("deletes", deletes); err != nil {
		return utils.CreateError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	reads := flattenRulesList(ctx, response.Reads)
	tflog.Debug(ctx, fmt.Sprintf("flattened reads %#v", reads))
	if err := d.Set("reads", reads); err != nil {
		return utils.CreateError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	updates := flattenRulesList(ctx, response.Updates)
	tflog.Debug(ctx, fmt.Sprintf("flattened updates %#v", updates))
	if err := d.Set("updates", updates); err != nil {
		return utils.CreateError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	if response.Identities != nil {
		if response.Identities.DBRoles != nil || response.Identities.Users != nil ||
			response.Identities.Groups != nil || response.Identities.Services != nil {
			identities := flattenIdentities(ctx, response.Identities)
			tflog.Debug(ctx, fmt.Sprintf("flattened identities %#v", identities))
			if err := d.Set("identities", identities); err != nil {
				return utils.CreateError("Unable to read policy rule", fmt.Sprintf("%v", err))
			}
		}
	}

	d.Set("hosts", response.Hosts)

	// Computed arguments
	d.Set("policy_rule_id", policyRuleID)

	tflog.Debug(ctx, "resourcePolicyRuleRead")
	return diag.Diagnostics{}
}

func resourcePolicyRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "resourcePolicyRuleUpdate")
	c := m.(*client.Client)

	policyRule := getPolicyRuleInfoFromResource(ctx, d)

	policyID, policyRuleID := unmarshalPolicyRuleID(d)
	url := fmt.Sprintf("https://%s/v1/policies/%s/rules/%s", c.ControlPlane,
		policyID, policyRuleID)

	_, err := c.DoRequest(url, http.MethodPut, policyRule)
	if err != nil {
		return utils.CreateError("Unable to update policy rule", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, "End resourcePolicyRuleUpdate")

	return resourcePolicyRuleRead(ctx, d, m)
}

func policyRuleDeleteConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "PolicyRuleDelete",
		Type:         operationtype.Delete,
		HttpMethod:   http.MethodDelete,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			policyID, policyRuleID := unmarshalPolicyRuleID(d)
			return fmt.Sprintf("https://%s/v1/policies/%s/rules/%s",
				c.ControlPlane, policyID, policyRuleID)
		},
		RequestErrorHandler: &core.DeleteIgnoreHttpNotFound{ResName: "Policy Rule"},
	}
}

func getStrListFromInterfaceList(ctx context.Context, interfaceList []interface{}) []string {
	tflog.Debug(ctx, "Init getStrListFromInterfaceList")

	strList := []string{}

	for _, i := range interfaceList {
		strList = append(strList, i.(string))
	}

	tflog.Debug(ctx, "End getStrListFromInterfaceList")

	return strList
}

func getDatasetRewrites(ctx context.Context, datasetList []interface{}) []DatasetRewrite {
	tflog.Debug(ctx, "Init getDatasetRewrites")

	datasetRewrites := make([]DatasetRewrite, 0, len(datasetList))

	for _, d := range datasetList {
		datasetMap := d.(map[string]interface{})

		datasetRewrite := DatasetRewrite{
			Dataset:      datasetMap["dataset"].(string),
			Repo:         datasetMap["repo"].(string),
			Substitution: datasetMap["substitution"].(string),
			Parameters:   getStrListFromInterfaceList(ctx, datasetMap["parameters"].([]interface{})),
		}

		datasetRewrites = append(datasetRewrites, datasetRewrite)
	}

	tflog.Debug(ctx, "End getDatasetRewrites")

	return datasetRewrites
}

func getRuleListFromResource(ctx context.Context, d *schema.ResourceData, name string) []Rule {
	tflog.Debug(ctx, "Init getRuleListFromResource")
	ruleInfoList := d.Get(name).([]interface{})
	ruleList := make([]Rule, 0, len(ruleInfoList))

	for _, ruleInterface := range ruleInfoList {
		ruleMap := ruleInterface.(map[string]interface{})

		rule := Rule{
			AdditionalChecks: ruleMap["additional_checks"].(string),
			Data:             getStrListFromInterfaceList(ctx, ruleMap["data"].([]interface{})),
			DatasetRewrites:  getDatasetRewrites(ctx, ruleMap["dataset_rewrites"].([]interface{})),
			Rows:             ruleMap["rows"].(int),
			Severity:         ruleMap["severity"].(string),
			RateLimit:        ruleMap["rate_limit"].(int),
		}

		ruleList = append(ruleList, rule)
	}
	tflog.Debug(ctx, "End getRuleListFromResource")

	return ruleList
}

func getPolicyRuleInfoFromResource(ctx context.Context, d *schema.ResourceData) PolicyRule {
	tflog.Debug(ctx, "Init getPolicyRuleInfoFromResource")
	hosts := utils.GetStrListFromSchemaField(d, "hosts")

	identity := d.Get("identities").([]interface{})

	var identities *Identity
	for _, id := range identity {
		idMap := id.(map[string]interface{})

		identities = &Identity{
			DBRoles:  getStrListFromInterfaceList(ctx, idMap["db_roles"].([]interface{})),
			Groups:   getStrListFromInterfaceList(ctx, idMap["groups"].([]interface{})),
			Services: getStrListFromInterfaceList(ctx, idMap["services"].([]interface{})),
			Users:    getStrListFromInterfaceList(ctx, idMap["users"].([]interface{})),
		}
	}

	policyRule := PolicyRule{
		Deletes:    getRuleListFromResource(ctx, d, "deletes"),
		Hosts:      hosts,
		Identities: identities,
		Reads:      getRuleListFromResource(ctx, d, "reads"),
		Updates:    getRuleListFromResource(ctx, d, "updates"),
	}

	tflog.Debug(ctx, "End getPolicyRuleInfoFromResource")

	return policyRule
}

func flattenIdentities(ctx context.Context, identities *Identity) []interface{} {
	tflog.Debug(ctx, "Init flattenIdentities")
	tflog.Debug(ctx, fmt.Sprintf("identities %#v", identities))
	identityMap := make(map[string]interface{})

	identityMap["db_roles"] = identities.DBRoles
	identityMap["groups"] = identities.Groups
	identityMap["services"] = identities.Services
	identityMap["users"] = identities.Users

	tflog.Debug(ctx, "End flattenIdentities")
	return []interface{}{identityMap}
}

func flattenRulesList(ctx context.Context, rulesList []Rule) []interface{} {
	tflog.Debug(ctx, "Init flattenRulesList")
	if rulesList != nil {
		rules := make([]interface{}, len(rulesList), len(rulesList))

		for i, rule := range rulesList {
			ruleMap := make(map[string]interface{})

			datasetRewriteList := make([]interface{}, len(rule.DatasetRewrites), len(rule.DatasetRewrites))

			for j, datasetRewrite := range rule.DatasetRewrites {
				drMap := make(map[string]interface{})
				drMap["dataset"] = datasetRewrite.Dataset
				drMap["repo"] = datasetRewrite.Repo
				drMap["substitution"] = datasetRewrite.Substitution
				drMap["parameters"] = datasetRewrite.Parameters

				datasetRewriteList[j] = drMap
			}

			ruleMap["additional_checks"] = rule.AdditionalChecks
			ruleMap["data"] = rule.Data
			ruleMap["dataset_rewrites"] = datasetRewriteList
			ruleMap["rows"] = rule.Rows
			ruleMap["severity"] = rule.Severity
			ruleMap["rate_limit"] = rule.RateLimit

			rules[i] = ruleMap
		}

		return rules
	}
	tflog.Debug(ctx, "End flattenRulesList")

	return make([]interface{}, 0)
}
