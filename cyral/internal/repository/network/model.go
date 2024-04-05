package network

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NetworkAccessPolicyUpsertResp struct {
	Policy NetworkAccessPolicy `json:"policy"`
}

func (resp NetworkAccessPolicyUpsertResp) WriteToSchema(d *schema.ResourceData) error {
	return resp.Policy.WriteToSchema(d)
}

type NetworkAccessPolicy struct {
	Enabled            bool `json:"enabled"`
	NetworkAccessRules `json:"networkAccessRules,omitempty"`
}

type NetworkAccessRules struct {
	RulesBlockAccess bool                `json:"rulesBlockAccess"`
	Rules            []NetworkAccessRule `json:"rules"`
}

type NetworkAccessRule struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	DBAccounts  []string `json:"dbAccounts,omitempty"`
	SourceIPs   []string `json:"sourceIPs,omitempty"`
}

func (nap *NetworkAccessPolicy) ReadFromSchema(d *schema.ResourceData) error {
	nap.Enabled = d.Get("enabled").(bool)

	var networkAccessRulesIfaces []interface{}
	if set, ok := d.GetOk("network_access_rule"); ok {
		networkAccessRulesIfaces = set.(*schema.Set).List()
	} else {
		return nil
	}

	nap.NetworkAccessRules = NetworkAccessRules{
		RulesBlockAccess: d.Get("network_access_rules_block_access").(bool),
		Rules:            []NetworkAccessRule{},
	}
	for _, networkAccessRuleIface := range networkAccessRulesIfaces {
		networkAccessRuleMap := networkAccessRuleIface.(map[string]interface{})
		nap.NetworkAccessRules.Rules = append(nap.NetworkAccessRules.Rules,
			NetworkAccessRule{
				Name:        networkAccessRuleMap["name"].(string),
				Description: networkAccessRuleMap["description"].(string),
				DBAccounts:  utils.GetStrList(networkAccessRuleMap, "db_accounts"),
				SourceIPs:   utils.GetStrList(networkAccessRuleMap, "source_ips"),
			})
	}

	return nil
}

func (nap *NetworkAccessPolicy) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get("repository_id").(string))
	d.Set("enabled", nap.Enabled)
	d.Set("network_access_rules_block_access", nap.NetworkAccessRules.RulesBlockAccess)

	var networkAccessRules []interface{}
	for _, rule := range nap.NetworkAccessRules.Rules {
		rulesMap := map[string]interface{}{
			"name":        rule.Name,
			"description": rule.Description,
			"db_accounts": rule.DBAccounts,
			"source_ips":  rule.SourceIPs,
		}
		networkAccessRules = append(networkAccessRules, rulesMap)
	}

	return d.Set("network_access_rule", networkAccessRules)
}
