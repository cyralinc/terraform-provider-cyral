package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func resourcePolicyRule() *schema.Resource {
	ruleSchema := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"additional_checks": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"data": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"dataset_rewrites": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"dataset": {
								Type:     schema.TypeString,
								Required: true,
							},
							"repo": {
								Type:     schema.TypeString,
								Required: true,
							},
							"substitution": {
								Type:     schema.TypeString,
								Required: true,
							},
							"parameters": {
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
					Description: "How many rows can be return by the policy rule. Use positive integer numbers to define how many rows. If you want to define `any` number of rows, set as `-1`.",
					Type:        schema.TypeInt,
					Required:    true,
				},
				"severity": {
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

	return &schema.Resource{
		CreateContext: resourcePolicyRuleCreate,
		ReadContext:   resourcePolicyRuleRead,
		UpdateContext: resourcePolicyRuleUpdate,
		DeleteContext: resourcePolicyRuleDelete,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deletes": ruleSchema,
			"reads":   ruleSchema,
			"updates": ruleSchema,
			"hosts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"identities": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_roles": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"groups": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"services": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"users": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourcePolicyRuleCreate")

	return resourcePolicyRuleRead(ctx, d, m)
}

func resourcePolicyRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRuleRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/policies/%s/rules/%s", c.ControlPlane, d.Get("policy_id").(string), d.Id())

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

	log.Printf("[DEBUG] End resourcePolicyRuleRead")
	return diag.Diagnostics{}
}

func resourcePolicyRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRuleUpdate")
	c := m.(*client.Client)

	policyRule := getPolicyRuleInfoFromResource(d)

	url := fmt.Sprintf("https://%s/v1/policies/%s/rules/%s", c.ControlPlane, d.Get("policy_id").(string), d.Id())

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

	url := fmt.Sprintf("https://%s/v1/policies/%s/rules/%s", c.ControlPlane, d.Get("policy_id").(string), d.Id())

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
