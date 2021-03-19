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

type CreatePolicyRuleResponse struct {
	ID string `json:"ID"`
}

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
	DatasetRewrites  []DatasetRewrite `json:"datasetRewrite,omitempty"`
	Rows             int64            `json:"rows"`
	Severity         string           `json:"severity"`
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
		Required: true,
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
					Type:     schema.TypeInt,
					Required: true,
				},
				"severity": {
					Type:     schema.TypeString,
					Optional: true,
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

	url := fmt.Sprintf("https://%s/v1/policies/%s/rules", c.ControlPlane, policyID)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create policy rule", fmt.Sprintf("%v", err))
	}

	response := CreatePolicyRuleResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)
	//update policy `last_updated` field?

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
	if err := d.Set("deletes", deletes); err != nil {
		return createError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	reads := flattenRulesList(response.Reads)
	if err := d.Set("reads", reads); err != nil {
		return createError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	updates := flattenRulesList(response.Updates)
	if err := d.Set("updates", updates); err != nil {
		return createError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	identities := flattenIdentities(response.Identities)
	if err := d.Set("identities", identities); err != nil {
		return createError("Unable to read policy rule", fmt.Sprintf("%v", err))
	}

	d.Set("hosts", response.Hosts)

	// log.Printf("[DEBUG] resourcePolicyRuleRead - policyRule: %#v", policyRule)

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
	strList := []string{}

	for _, i := range interfaceList {
		strList = append(strList, i.(string))
	}

	return strList
}

func getDatasetRewrites(datasetList []interface{}) []DatasetRewrite {
	datasetRewrites := make([]DatasetRewrite, len(datasetList), len(datasetList))

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

	return datasetRewrites
}

func getRuleListFromResource(d *schema.ResourceData, name string) []Rule {
	ruleInfoList := d.Get(name).([]interface{})
	ruleList := make([]Rule, len(ruleInfoList), len(ruleInfoList))

	for _, ruleInterface := range ruleInfoList {
		ruleMap := ruleInterface.(map[string]interface{})

		rule := Rule{
			AdditionalChecks: ruleMap["additional_checks"].(string),
			Data:             getStrListFromInterfaceList(ruleMap["data"].([]interface{})),
			DatasetRewrites:  getDatasetRewrites(ruleMap["dataset_rewrites"].([]interface{})),
			Rows:             ruleMap["rows"].(int64),
			Severity:         ruleMap["severity"].(string),
		}

		ruleList = append(ruleList, rule)
	}

	return ruleList
}

func getPolicyRuleInfoFromResource(d *schema.ResourceData) PolicyRule {
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

	return policyRule
}

func flattenIdentities(identities Identity) []interface{} {
	identityMap := make(map[string]interface{})

	identityMap["db_roles"] = identities.DBRoles
	identityMap["groups"] = identities.Groups
	identityMap["services"] = identities.Services
	identityMap["users"] = identities.Users

	return []interface{}{identityMap}
}

func flattenRulesList(rulesList []Rule) []interface{} {
	if rulesList != nil {
		rules := make([]interface{}, 0, len(rulesList))

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

			rules[i] = ruleMap
		}

		return rules
	}

	return make([]interface{}, 0)
}
