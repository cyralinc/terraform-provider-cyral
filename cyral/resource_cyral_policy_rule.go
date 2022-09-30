package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolicyRule struct {
	Deletes    []Rule   `json:"deletes,omitempty"`
	Hosts      []string `json:"hosts,omitempty"`
	Identities Identity `json:"identities,omitempty"`
	Reads      []Rule   `json:"reads,omitempty"`
	RuleID     string   `json:"ruleId"`
	Updates    []Rule   `json:"updates,omitempty"`
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
	ids, err := unmarshalComposedID(d.Id(), "/", 2)
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
func upgradePolicyRuleV0(
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
	newID := marshalComposedID([]string{policyID, policyRuleID}, "/")
	rawState["id"] = newID

	return rawState, nil
}

func resourcePolicyRule() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [policy rules](https://cyral.com/docs/reference/policy/#rules). See also: [Policy](./policy.md)" +
			"\n\n> Notes:\n>" +
			"\n> 1. Unless you create a default rule, users and groups only have the rights you explicitly grant them." +
			"\n> 2. Each contexted rule comprises these fields: `data`, `rows`, `severity` `additional_checks`, `dataset_rewrites`. The only required fields are `data` and `rows`." +
			"\n> 3. The rules block does not need to include all three operation types (reads, updates and deletes); actions you omit are disallowed." +
			"\n> 4. If you do not include a hosts block, Cyral does not enforce limits based on the connecting client's host address." +
			"\n\nFor more information, see the [Policy Guide](https://cyral.com/docs/policy#the-rules-block-of-a-policy).",

		CreateContext: resourcePolicyRuleCreate,
		ReadContext:   resourcePolicyRuleRead,
		UpdateContext: resourcePolicyRuleUpdate,
		DeleteContext: resourcePolicyRuleDelete,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: policyRuleResourceSchemaV0().
					CoreConfigSchema().ImpliedType(),
				Upgrade: upgradePolicyRuleV0,
			},
		},

		Schema: policyRuleResourceSchemaV0().Schema,

		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := unmarshalComposedID(d.Id(), "/", 2)
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
	log.Printf("[DEBUG] Init resourcePolicyRuleCreate")
	c := m.(*client.Client)

	policyID := d.Get("policy_id").(string)
	resourceData := getPolicyRuleInfoFromResource(d)
	log.Printf("[DEBUG] resourcePolicyRuleCreate - policyRule: %#v", resourceData)

	url := fmt.Sprintf("https://%s/v1/policies/%s/rules", c.ControlPlane, policyID)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create policy rule", fmt.Sprintf("%v", err))
	}

	response := IDBasedResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(marshalComposedID([]string{
		policyID,
		response.ID},
		"/"))

	log.Printf("[DEBUG] End resourcePolicyRuleCreate")

	return resourcePolicyRuleRead(ctx, d, m)
}

func resourcePolicyRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRuleRead")
	c := m.(*client.Client)

	policyID, policyRuleID := unmarshalPolicyRuleID(d)
	url := fmt.Sprintf("https://%s/v1/policies/%s/rules/%s",
		c.ControlPlane, policyID, policyRuleID)

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	response := PolicyRule{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	deletes := flattenRulesList(response.Deletes)
	log.Printf("[DEBUG] flattened deletes %#v", deletes)
	if err := d.Set("deletes", deletes); err != nil {
		return createError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	reads := flattenRulesList(response.Reads)
	log.Printf("[DEBUG] flattened reads %#v", reads)
	if err := d.Set("reads", reads); err != nil {
		return createError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	updates := flattenRulesList(response.Updates)
	log.Printf("[DEBUG] flattened updates %#v", updates)
	if err := d.Set("updates", updates); err != nil {
		return createError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	if response.Identities.DBRoles != nil || response.Identities.Users != nil || response.Identities.Groups != nil || response.Identities.Services != nil {
		identities := flattenIdentities(response.Identities)
		log.Printf("[DEBUG] flattened identities %#v", identities)
		if err := d.Set("identities", identities); err != nil {
			return createError("Unable to read policy rule", fmt.Sprintf("%v", err))
		}
	}

	d.Set("hosts", response.Hosts)

	// Computed arguments
	d.Set("policy_rule_id", policyRuleID)

	log.Printf("[DEBUG] End resourcePolicyRuleRead")
	return diag.Diagnostics{}
}

func resourcePolicyRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRuleUpdate")
	c := m.(*client.Client)

	policyRule := getPolicyRuleInfoFromResource(d)

	policyID, policyRuleID := unmarshalPolicyRuleID(d)
	url := fmt.Sprintf("https://%s/v1/policies/%s/rules/%s", c.ControlPlane,
		policyID, policyRuleID)

	_, err := c.DoRequest(url, http.MethodPut, policyRule)
	if err != nil {
		return createError("Unable to update policy rule", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourcePolicyRuleUpdate")

	return resourcePolicyRuleRead(ctx, d, m)
}

func resourcePolicyRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRuleDelete")
	c := m.(*client.Client)

	policyID, policyRuleID := unmarshalPolicyRuleID(d)
	url := fmt.Sprintf("https://%s/v1/policies/%s/rules/%s",
		c.ControlPlane, policyID, policyRuleID)

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete policy rule", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourcePolicyRuleDelete")

	return diag.Diagnostics{}
}

func getStrListFromInterfaceList(interfaceList []interface{}) []string {
	log.Printf("[DEBUG] Init getStrListFromInterfaceList")

	strList := []string{}

	for _, i := range interfaceList {
		strList = append(strList, i.(string))
	}

	log.Printf("[DEBUG] End getStrListFromInterfaceList")

	return strList
}

func getDatasetRewrites(datasetList []interface{}) []DatasetRewrite {
	log.Printf("[DEBUG] Init getDatasetRewrites")

	datasetRewrites := make([]DatasetRewrite, 0, len(datasetList))

	for _, d := range datasetList {
		datasetMap := d.(map[string]interface{})

		datasetRewrite := DatasetRewrite{
			Dataset:      datasetMap["dataset"].(string),
			Repo:         datasetMap["repo"].(string),
			Substitution: datasetMap["substitution"].(string),
			Parameters:   getStrListFromInterfaceList(datasetMap["parameters"].([]interface{})),
		}

		datasetRewrites = append(datasetRewrites, datasetRewrite)
	}

	log.Printf("[DEBUG] End getDatasetRewrites")

	return datasetRewrites
}

func getRuleListFromResource(d *schema.ResourceData, name string) []Rule {
	log.Printf("[DEBUG] Init getRuleListFromResource")
	ruleInfoList := d.Get(name).([]interface{})
	ruleList := make([]Rule, 0, len(ruleInfoList))

	for _, ruleInterface := range ruleInfoList {
		ruleMap := ruleInterface.(map[string]interface{})

		rule := Rule{
			AdditionalChecks: ruleMap["additional_checks"].(string),
			Data:             getStrListFromInterfaceList(ruleMap["data"].([]interface{})),
			DatasetRewrites:  getDatasetRewrites(ruleMap["dataset_rewrites"].([]interface{})),
			Rows:             ruleMap["rows"].(int),
			Severity:         ruleMap["severity"].(string),
			RateLimit:        ruleMap["rate_limit"].(int),
		}

		ruleList = append(ruleList, rule)
	}
	log.Printf("[DEBUG] End getRuleListFromResource")

	return ruleList
}

func getPolicyRuleInfoFromResource(d *schema.ResourceData) PolicyRule {
	log.Printf("[DEBUG] Init getPolicyRuleInfoFromResource")
	hosts := getStrListFromSchemaField(d, "hosts")

	identity := d.Get("identities").([]interface{})

	var identities Identity
	for _, id := range identity {
		idMap := id.(map[string]interface{})

		identities = Identity{
			DBRoles:  getStrListFromInterfaceList(idMap["db_roles"].([]interface{})),
			Groups:   getStrListFromInterfaceList(idMap["groups"].([]interface{})),
			Services: getStrListFromInterfaceList(idMap["services"].([]interface{})),
			Users:    getStrListFromInterfaceList(idMap["users"].([]interface{})),
		}
	}

	policyRule := PolicyRule{
		Deletes:    getRuleListFromResource(d, "deletes"),
		Hosts:      hosts,
		Identities: identities,
		Reads:      getRuleListFromResource(d, "reads"),
		Updates:    getRuleListFromResource(d, "updates"),
	}

	log.Printf("[DEBUG] End getPolicyRuleInfoFromResource")

	return policyRule
}

func flattenIdentities(identities Identity) []interface{} {
	log.Printf("[DEBUG] Init flattenIdentities")
	log.Printf("[DEBUG] identities %#v", identities)
	identityMap := make(map[string]interface{})

	identityMap["db_roles"] = identities.DBRoles
	identityMap["groups"] = identities.Groups
	identityMap["services"] = identities.Services
	identityMap["users"] = identities.Users

	log.Printf("[DEBUG] End flattenIdentities")
	return []interface{}{identityMap}
}

func flattenRulesList(rulesList []Rule) []interface{} {
	log.Printf("[DEBUG] Init flattenRulesList")
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
	log.Printf("[DEBUG] End flattenRulesList")

	return make([]interface{}, 0)
}
