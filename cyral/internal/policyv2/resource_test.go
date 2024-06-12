package policyv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

func initialPolicyConfig() map[string]interface{} {
	document := `{"governedData":{"locations":["gym_db.users"]},"readRules":[{"conditions":[{"attribute":"identity.userGroups","operator":"contains","value":"USERS"}],"constraints":{"datasetRewrite":"SELECT * FROM ${dataset} WHERE email = '${identity.endUserEmail}'"}},{"conditions":[],"constraints":{}}]}`

	return map[string]interface{}{
		"name":        "policy1",
		"description": "Local policies for users table access",
		"enabled":     true,
		"tags":        []string{"tag1", "tag2"},
		"scope":       map[string][]string{"repoIds": {"repo1", "repo2"}},
		"valid_from":  "2023-01-01T00:00:00Z",
		"valid_until": "2034-01-01T00:00:00Z",
		"document":    document,
		"enforced":    true,
		"type":        "POLICY_TYPE_LOCAL",
	}
}

func updatedPolicyConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":        "policy2",
		"description": "Updated local policies for users table access",
		"enabled":     false,
		"tags":        []string{"tag3"},
		"scope":       map[string][]string{"repoIds": {"repo3"}},
		"valid_from":  "2023-06-01T00:00:00Z",
		"valid_until": "2035-06-01T00:00:00Z",
		"document":    `{"governedData":{"locations":["gym_db.users"]},"readRules":[{"conditions":[{"attribute":"identity.userGroups","operator":"contains","value":"ADMINS"}],"constraints":{"datasetRewrite":"SELECT * FROM ${dataset} WHERE email = '${identity.endUserEmail}'"}},{"conditions":[],"constraints":{}}]}`,
		"enforced":    false,
		"type":        "POLICY_TYPE_LOCAL",
	}
}
func minimalPolicyConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":     "policy2",
		"document": `{"governedData":{"locations":["gym_db.users"]},"readRules":[{"conditions":[{"attribute":"identity.userGroups","operator":"contains","value":"ADMINS"}],"constraints":{"datasetRewrite":"SELECT * FROM ${dataset} WHERE email = '${identity.endUserEmail}'"}},{"conditions":[],"constraints":{}}]}`,
		"type":     "local",
	}
}

func TestAccPolicyV2Resource(t *testing.T) {
	testInitialConfig, testInitialFunc := setupPolicyTest("main_test", initialPolicyConfig())
	testUpdatedConfig, testUpdatedFunc := setupPolicyTest("main_test", updatedPolicyConfig())
	resource.Test(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialFunc,
			},
			{
				Config: testUpdatedConfig,
				Check:  testUpdatedFunc,
			},
		},
	})
}

func TestMinimalPolicyV2Resource(t *testing.T) {
	testMinimalConfig, testMinimalFunc := setupMinimalPolicyTest("main_test", minimalPolicyConfig())
	resource.Test(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testMinimalConfig,
				Check:  testMinimalFunc,
			},
		},
	})
}

func setupPolicyTest(resName string, policy map[string]interface{}) (string, resource.TestCheckFunc) {
	config := utils.FormatPolicyIntoConfig(resName, policy)
	resourceFullName := fmt.Sprintf("cyral_policy_v2.%s", resName)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", policy["name"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "description", policy["description"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", fmt.Sprintf("%v", policy["enabled"])),
		resource.TestCheckResourceAttr(resourceFullName, "valid_from", policy["valid_from"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "valid_until", policy["valid_until"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "document", policy["document"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "type", policy["type"].(string)),
	)

	return config, testFunction
}
func setupMinimalPolicyTest(resName string, policy map[string]interface{}) (string, resource.TestCheckFunc) {
	config := utils.FormatPolicyIntoConfig(resName, policy)
	resourceFullName := fmt.Sprintf("cyral_policy_v2.%s", resName)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", policy["name"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "document", policy["document"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "type", policy["type"].(string)),
	)

	return config, testFunction
}
