package policyset_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
)

func initialPolicySetConfig() map[string]interface{} {
	wizardParameters := `{"failClosed":true,"denyByDefault":true}`

	return map[string]interface{}{
		"wizard_id":         "repo-lockdown",
		"name":              "Lockdown all repos",
		"description":       "Policy set created via wizard",
		"tags":              []string{"default", "block", "fail closed"},
		"scope":             map[string][]string{"repo_ids": {}}, // Empty scope means to target all repos
		"wizard_parameters": wizardParameters,
		"enabled":           true,
	}
}

func updatedPolicySetConfig() map[string]interface{} {
	wizardParameters := `{"failClosed":false,"denyByDefault":true}`

	return map[string]interface{}{
		"wizard_id":         "repo-lockdown",
		"name":              "Lockdown all repos",
		"description":       "Updated policy set created via wizard",
		"tags":              []string{"default", "block", "fail open"},
		"scope":             map[string][]string{"repo_ids": {}},
		"wizard_parameters": wizardParameters,
		"enabled":           true,
	}
}

func TestAccPolicyWizardV1Resource(t *testing.T) {
	testInitialConfig, testInitialFunc := setupPolicySetTest("main_test", initialPolicySetConfig())
	testUpdatedConfig, testUpdatedFunc := setupPolicySetTest("main_test", updatedPolicySetConfig())
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
			{
				// Test importing the resource
				ResourceName:      "cyral_policy_set.main_test",
				ImportState:       true,
				ImportStateIdFunc: testImportStateIdFunc("cyral_policy_set.main_test"),
				ImportStateVerify: true,
			},
		},
	})
}

func setupPolicySetTest(resName string, policySet map[string]interface{}) (string, resource.TestCheckFunc) {
	config := formatPolicySetIntoConfig(resName, policySet)
	resourceFullName := fmt.Sprintf("cyral_policy_set.%s", resName)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "wizard_id", policySet["wizard_id"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "name", policySet["name"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "description", policySet["description"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "wizard_parameters", policySet["wizard_parameters"].(string)),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", fmt.Sprintf("%v", policySet["enabled"])),
	)

	// Check tags
	if tags, ok := policySet["tags"].([]string); ok && len(tags) > 0 {
		for i, tag := range tags {
			key := fmt.Sprintf("tags.%d", i)
			testFunction = resource.ComposeTestCheckFunc(
				testFunction,
				resource.TestCheckResourceAttr(resourceFullName, key, tag),
			)
		}
	}

	// Check scope
	if scope, ok := policySet["scope"].(map[string][]string); ok {
		if repoIds, ok := scope["repo_ids"]; ok && len(repoIds) > 0 {
			for i, repoId := range repoIds {
				key := fmt.Sprintf("scope.0.repo_ids.%d", i)
				testFunction = resource.ComposeTestCheckFunc(
					testFunction,
					resource.TestCheckResourceAttr(resourceFullName, key, repoId),
				)
			}
		}
	}

	return config, testFunction
}

func formatPolicySetIntoConfig(resName string, policySet map[string]interface{}) string {
	config := fmt.Sprintf(`
resource "cyral_policy_set" "%s" {
  wizard_id         = "%s"
  name              = "%s"
`, resName, policySet["wizard_id"].(string), policySet["name"].(string))

	if description, ok := policySet["description"].(string); ok && description != "" {
		config += fmt.Sprintf(`  description       = "%s"
`, description)
	}

	// Handle tags
	if tags, ok := policySet["tags"].([]string); ok && len(tags) > 0 {
		config += "  tags = [\n"
		for _, tag := range tags {
			config += fmt.Sprintf(`    "%s",
`, tag)
		}
		config += "  ]\n"
	}

	// Handle scope
	if scope, ok := policySet["scope"].(map[string][]string); ok {
		if repoIds, ok := scope["repo_ids"]; ok && len(repoIds) > 0 {
			config += "  scope {\n"
			config += "    repo_ids = [\n"
			for _, repoId := range repoIds {
				config += fmt.Sprintf(`      "%s",
`, repoId)
			}
			config += "    ]\n"
			config += "  }\n"
		}
	}

	// Properly escape wizard_parameters
	config += fmt.Sprintf("  wizard_parameters = %q\n", policySet["wizard_parameters"].(string))

	if enabled, ok := policySet["enabled"].(bool); ok {
		config += fmt.Sprintf("  enabled = %v\n", enabled)
	}

	config += "}\n"

	return config
}

func testImportStateIdFunc(resourceName string) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceName)
		}
		// The ID is the policy set ID
		return rs.Primary.ID, nil
	}
}
