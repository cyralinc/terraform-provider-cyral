package accessrules

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AccessRulesIdentity struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type AccessRulesConfig struct {
	AuthorizationPolicyInstanceIDs []string `json:"authorizationPolicyInstanceIDs"`
}

type AccessRule struct {
	Identity   *AccessRulesIdentity `json:"identity"`
	ValidFrom  *string              `json:"validFrom"`
	ValidUntil *string              `json:"validUntil"`
	Config     *AccessRulesConfig   `json:"config"`
}

type AccessRulesResource struct {
	AccessRules []*AccessRule `json:"accessRules"`
}

type AccessRulesResponse struct {
	AccessRules []*AccessRule `json:"accessRules"`
}

// WriteToSchema is used when reading a resource. It takes whatever the API
// read call returned and translates it into the Terraform schema.
func (arr *AccessRulesResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(
		utils.MarshalComposedID(
			[]string{
				d.Get("repository_id").(string),
				d.Get("user_account_id").(string),
			},
			"/",
		),
	)
	// We'll have to build the access rule set in the format expected by Terraform,
	// which boils down to doing a bunch of type casts
	rules := make([]interface{}, 0, len(arr.AccessRules))
	for _, rule := range arr.AccessRules {
		m := make(map[string]interface{})

		m["identity"] = []interface{}{
			map[string]interface{}{
				"type": rule.Identity.Type,
				"name": rule.Identity.Name,
			},
		}

		m["valid_from"] = rule.ValidFrom
		m["valid_until"] = rule.ValidUntil

		if rule.Config != nil && len(rule.Config.AuthorizationPolicyInstanceIDs) > 0 {
			m["config"] = []interface{}{
				map[string]interface{}{
					"policy_ids": rule.Config.AuthorizationPolicyInstanceIDs,
				},
			}
		}

		rules = append(rules, m)
	}
	return d.Set("rule", rules)
}

// ReadFromSchema is called when *creating* or *updating* a resource.
// Essentially, it translates the stuff from the .tf file into whatever the
// API expects. The `AccessRulesResource` will be marshalled verbatim, so
// make sure that it matches *exactly* what the API needs.
func (arr *AccessRulesResource) ReadFromSchema(d *schema.ResourceData) error {
	rules := d.Get("rule").([]interface{})
	var accessRules []*AccessRule

	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})

		accessRule := &AccessRule{}

		identity := ruleMap["identity"].(*schema.Set).List()[0].(map[string]interface{})
		accessRule.Identity = &AccessRulesIdentity{
			Type: identity["type"].(string),
			Name: identity["name"].(string),
		}

		validFrom := ruleMap["valid_from"].(string)
		if validFrom != "" {
			accessRule.ValidFrom = &validFrom
		}

		validUntil := ruleMap["valid_until"].(string)
		if validUntil != "" {
			accessRule.ValidUntil = &validUntil
		}

		conf, ok := ruleMap["config"]
		if ok {
			confList := conf.(*schema.Set).List()
			if len(confList) > 0 {
				config := confList[0].(map[string]interface{})
				policyIDs := config["policy_ids"].([]interface{})
				ids := make([]string, 0, len(policyIDs))
				for _, policyID := range policyIDs {
					ids = append(ids, policyID.(string))
				}
				accessRule.Config = &AccessRulesConfig{
					AuthorizationPolicyInstanceIDs: ids,
				}
			}
		}

		accessRules = append(accessRules, accessRule)
	}

	arr.AccessRules = accessRules
	return nil
}
