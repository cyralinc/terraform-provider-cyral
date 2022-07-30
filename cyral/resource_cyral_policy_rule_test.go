package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type PolicyRuleConfig struct {
	DeletedSeverity  string
	UpdatedSeverity  string
	ReadSeverity     string
	DeletedRateLimit int
	UpdatedRateLimit int
	ReadRateLimit    int
}

var initialPolicyRuleConfig PolicyRuleConfig = PolicyRuleConfig{
	DeletedSeverity:  "medium",
	UpdatedSeverity:  "high",
	ReadSeverity:     "low",
	DeletedRateLimit: 1,
	UpdatedRateLimit: 2,
	ReadRateLimit:    3,
}

var updated1PolicyRuleConfig PolicyRuleConfig = PolicyRuleConfig{
	DeletedSeverity:  "low",
	UpdatedSeverity:  "medium",
	ReadSeverity:     "high",
	DeletedRateLimit: 2,
	UpdatedRateLimit: 3,
	ReadRateLimit:    4,
}

var updated2PolicyRuleConfig PolicyRuleConfig = PolicyRuleConfig{
	DeletedSeverity:  "high",
	UpdatedSeverity:  "low",
	ReadSeverity:     "medium",
	DeletedRateLimit: 5,
	UpdatedRateLimit: 6,
	ReadRateLimit:    7,
}

func TestAccPolicyRuleResource(t *testing.T) {
	testConfig, testFunc := setupPolicyRuleTest(initialPolicyRuleConfig)
	testUpdate1Config, testUpdate1Func := setupPolicyRuleTest(updated1PolicyRuleConfig)
	testUpdate2Config, testUpdate2Func := setupPolicyRuleTest(updated2PolicyRuleConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdate1Config,
				Check:  testUpdate1Func,
			},
			{
				Config: testUpdate2Config,
				Check:  testUpdate2Func,
			},
			// TODO: when import functionality for cyral_policy_rule
			// is implemented, add Import test.
		},
	})
}

func setupPolicyRuleTest(integrationData PolicyRuleConfig) (string, resource.TestCheckFunc) {
	configuration := formatPolicyRuleConfigIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_policy_rule.policy_rule_test", "deletes.0.severity", integrationData.DeletedSeverity),
		resource.TestCheckResourceAttr("cyral_policy_rule.policy_rule_test", "reads.0.severity", integrationData.ReadSeverity),
		resource.TestCheckResourceAttr("cyral_policy_rule.policy_rule_test", "updates.0.severity", integrationData.UpdatedSeverity),
	)

	return configuration, testFunction
}

func formatPolicyRuleConfigIntoConfig(data PolicyRuleConfig) string {
	return fmt.Sprintf(`

		resource "cyral_repository" "tf_test_repository" {
			type = "mysql"
			host = "http://mysql.local/"
			port = 3306
			name = "tf-test-mysql"
	  }

	  resource "cyral_sidecar" "tf_test_sidecar" {
			name = "tf-test-sidecar"
			deployment_method = "cloudFormation"
	  }

	  resource "cyral_repository_binding" "repo_binding" {
			enabled       = true
			repository_id = cyral_repository.tf_test_repository.id
			listener_port = 3307
			sidecar_id    = cyral_sidecar.tf_test_sidecar.id
	  }

	  resource "cyral_datamap" "test_datamap" {
			mapping {
				label = "TEST_CCN"
				data_location {
				repo       = cyral_repository.tf_test_repository.name
				attributes = ["database.table.column"]
				}
			}
	  }


	resource "cyral_policy" "policy_rule_test_policy" {
		data = ["TEST_CCN"]
		description = "description"
		enabled = true
		name = "policy_rule_test_policy"
		tags = ["PCI"]
	}

	resource "cyral_policy_rule" "policy_rule_test" {
		policy_id = cyral_policy.policy_rule_test_policy.id
		hosts = ["192.0.2.22", "203.0.113.16/28"]
		identities {
			groups = ["analyst"]
		}
		deletes {
			data = ["TEST_CCN"]
			rows = 1
			severity = "%s"
			rate_limit = %d
		}
		reads {
			data = ["TEST_CCN"]
			rows = 1
			severity = "%s"
			rate_limit = %d
		}
		updates {
			data = ["TEST_CCN"]
			rows = 1
			severity = "%s"
			rate_limit = %d
		}
	}`, data.DeletedSeverity, data.DeletedRateLimit, data.ReadSeverity, data.ReadRateLimit, data.UpdatedSeverity, data.UpdatedRateLimit)
}
