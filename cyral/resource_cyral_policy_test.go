package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type PolicyTestConfig struct {
	Data        []string
	Description string
	Enabled     bool
	Name        string
	Tags        []string
}

var initialPolicyConfig PolicyTestConfig = PolicyTestConfig{
	Data:        []string{"data"},
	Description: "description",
	Enabled:     false,
	Name:        "policy-test",
	Tags:        []string{"tag"},
}

var updatedPolicyConfig PolicyTestConfig = PolicyTestConfig{
	Data:        []string{"data-updated"},
	Description: "desctiption-updated",
	Enabled:     true,
	Name:        "policy-test",
	Tags:        []string{"tag-updated"},
}

func TestAccPolicyResource(t *testing.T) {
	testConfig, testFunc := setupPolicyTest(initialPolicyConfig)
	testUpdateConfig, testUpdateFunc := setupPolicyTest(updatedPolicyConfig)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_policy.policy_test",
			},
		},
	})
}

func setupPolicyTest(integrationData PolicyTestConfig) (string, resource.TestCheckFunc) {
	configuration := formatPolicyTestConfigIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_policy.policy_test", "data.#", fmt.Sprintf("%d", len(integrationData.Data))),
		resource.TestCheckResourceAttr("cyral_policy.policy_test", "data.0", integrationData.Data[0]),
		resource.TestCheckResourceAttr("cyral_policy.policy_test", "description", integrationData.Description),
		resource.TestCheckResourceAttr("cyral_policy.policy_test", "enabled", fmt.Sprintf("%t", integrationData.Enabled)),
		resource.TestCheckResourceAttr("cyral_policy.policy_test", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_policy.policy_test", "tags.#", fmt.Sprintf("%d", len(integrationData.Tags))),
		resource.TestCheckResourceAttr("cyral_policy.policy_test", "tags.0", integrationData.Tags[0]),
	)

	return configuration, testFunction
}

func formatPolicyTestConfigIntoConfig(data PolicyTestConfig) string {
	return fmt.Sprintf(`
	resource "cyral_policy" "policy_test" {
		data = [%s]
		description = "%s"
		enabled = %t
		name = "%s"
		tags = [%s]
	  }`, formatAttributes(data.Data), data.Description, data.Enabled, data.Name, formatAttributes(data.Tags))
}
