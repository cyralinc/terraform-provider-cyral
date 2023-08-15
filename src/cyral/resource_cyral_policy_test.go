package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	policyResourceName = "policy"
)

var initialPolicyConfig = Policy{
	Meta: &PolicyMetadata{
		Name:        accTestName(policyResourceName, "test"),
		Description: "description",
		Enabled:     false,
		Tags:        []string{"tag"},
	},
	Data: []string{"data"},
	Tags: []string{"DATA_TAG_TEST"},
}

var updatedPolicyConfig = Policy{
	Meta: &PolicyMetadata{
		Name:        accTestName(policyResourceName, "test-updated"),
		Description: "desctiption-updated",
		Enabled:     true,
		Tags:        []string{"tag-updated"},
	},
	Data: []string{"data-updated"},
	Tags: []string{"DATA_TAG_TEST_UPDATED"},
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

func setupPolicyTest(integrationData Policy) (string, resource.TestCheckFunc) {
	configuration := formatPolicyTestConfigIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			"cyral_policy.policy_test", "name",
			integrationData.Meta.Name,
		),
		resource.TestCheckResourceAttr(
			"cyral_policy.policy_test", "description",
			integrationData.Meta.Description,
		),
		resource.TestCheckResourceAttr(
			"cyral_policy.policy_test", "enabled",
			fmt.Sprintf("%t", integrationData.Meta.Enabled),
		),
		resource.TestCheckResourceAttr(
			"cyral_policy.policy_test", "data.#",
			fmt.Sprintf("%d", len(integrationData.Data)),
		),
		resource.TestCheckResourceAttr(
			"cyral_policy.policy_test", "data.0",
			integrationData.Data[0],
		),
		resource.TestCheckResourceAttr(
			"cyral_policy.policy_test", "data_label_tags.#",
			fmt.Sprintf("%d", len(integrationData.Tags)),
		),
		resource.TestCheckResourceAttr(
			"cyral_policy.policy_test", "data_label_tags.0",
			integrationData.Tags[0],
		),
		resource.TestCheckResourceAttr(
			"cyral_policy.policy_test", "metadata_tags.#",
			fmt.Sprintf("%d", len(integrationData.Meta.Tags)),
		),
		resource.TestCheckResourceAttr(
			"cyral_policy.policy_test", "metadata_tags.0",
			integrationData.Meta.Tags[0],
		),
	)

	return configuration, testFunction
}

func formatPolicyTestConfigIntoConfig(data Policy) string {
	return fmt.Sprintf(`
	resource "cyral_policy" "policy_test" {
		name = "%s"
		description = "%s"
		enabled = %t
		data = %s
		data_label_tags = %s
		metadata_tags = %s
	}`,
		data.Meta.Name,
		data.Meta.Description,
		data.Meta.Enabled,
		utils.ListToStr(data.Data),
		utils.ListToStr(data.Tags),
		utils.ListToStr(data.Meta.Tags),
	)
}
