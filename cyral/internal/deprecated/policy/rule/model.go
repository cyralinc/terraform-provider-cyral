package rule

import (
	"context"
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func (r PolicyRule) WriteToSchema(d *schema.ResourceData) error {
	ctx := context.Background()
	deletes := flattenRulesList(ctx, r.Deletes)
	tflog.Debug(ctx, fmt.Sprintf("flattened deletes %#v", deletes))
	if err := d.Set("deletes", deletes); err != nil {
		return fmt.Errorf("error setting 'deletes' field: %w", err)
	}

	reads := flattenRulesList(ctx, r.Reads)
	tflog.Debug(ctx, fmt.Sprintf("flattened reads %#v", reads))
	if err := d.Set("reads", reads); err != nil {
		return fmt.Errorf("error setting 'reads' field: %w", err)
	}

	updates := flattenRulesList(ctx, r.Updates)
	tflog.Debug(ctx, fmt.Sprintf("flattened updates %#v", updates))
	if err := d.Set("updates", updates); err != nil {
		return fmt.Errorf("error setting 'updates' field: %w", err)
	}

	if r.Identities != nil {
		if r.Identities.DBRoles != nil || r.Identities.Users != nil ||
			r.Identities.Groups != nil || r.Identities.Services != nil {
			identities := flattenIdentities(ctx, r.Identities)
			tflog.Debug(ctx, fmt.Sprintf("flattened identities %#v", identities))
			if err := d.Set("identities", identities); err != nil {
				return fmt.Errorf("error setting 'identities' field: %w", err)
			}
		}
	}

	d.Set("hosts", r.Hosts)

	_, policyRuleID := unmarshalPolicyRuleID(d)
	// Computed arguments
	d.Set("policy_rule_id", policyRuleID)

	return nil
}

func (r *PolicyRule) ReadFromSchema(d *schema.ResourceData) error {
	hosts := utils.GetStrListFromSchemaField(d, "hosts")

	identity := d.Get("identities").([]interface{})

	ctx := context.Background()

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

	r.Deletes = getRuleListFromResource(ctx, d, "deletes")
	r.Hosts = hosts
	r.Identities = identities
	r.Reads = getRuleListFromResource(ctx, d, "reads")
	r.Updates = getRuleListFromResource(ctx, d, "updates")

	return nil
}

func getRuleListFromResource(ctx context.Context, d *schema.ResourceData, name string) []Rule {
	tflog.Trace(ctx, "Init getRuleListFromResource")
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
	tflog.Trace(ctx, "End getRuleListFromResource")

	return ruleList
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

func getStrListFromInterfaceList(ctx context.Context, interfaceList []interface{}) []string {
	tflog.Trace(ctx, "Init getStrListFromInterfaceList")

	strList := []string{}

	for _, i := range interfaceList {
		strList = append(strList, i.(string))
	}

	tflog.Trace(ctx, "End getStrListFromInterfaceList")

	return strList
}

func getDatasetRewrites(ctx context.Context, datasetList []interface{}) []DatasetRewrite {
	tflog.Trace(ctx, "Init getDatasetRewrites")

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

	tflog.Trace(ctx, "End getDatasetRewrites")

	return datasetRewrites
}

func flattenIdentities(ctx context.Context, identities *Identity) []interface{} {
	tflog.Trace(ctx, "Init flattenIdentities")
	tflog.Trace(ctx, fmt.Sprintf("identities %#v", identities))
	identityMap := make(map[string]interface{})

	identityMap["db_roles"] = identities.DBRoles
	identityMap["groups"] = identities.Groups
	identityMap["services"] = identities.Services
	identityMap["users"] = identities.Users

	tflog.Trace(ctx, "End flattenIdentities")
	return []interface{}{identityMap}
}

func flattenRulesList(ctx context.Context, rulesList []Rule) []interface{} {
	tflog.Trace(ctx, "Init flattenRulesList")
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
	tflog.Trace(ctx, "End flattenRulesList")

	return make([]interface{}, 0)
}

type PolicyRuleIDBasedResponse struct {
	ID string `json:"id"`
}

func (r PolicyRuleIDBasedResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(utils.MarshalComposedID([]string{
		d.Get("policy_id").(string),
		r.ID},
		"/"))
	return nil
}
