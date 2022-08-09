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

	importStateResName := "cyral_policy_rule.policy_rule_test"

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateComposedIDFunc(
					importStateResName,
					[]string{"policy_id", "id"},
					"/",
				),
				ResourceName: importStateResName,
			},
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

// TODO: finish decomposing this function -aholmquist 2022-08-08
func formatPolicyRuleConfigIntoConfig(data PolicyRuleConfig) string {
	testLabelName := "TEST_CCN"

	var config string
	config += formatBasicRepositoryIntoConfig(
		basicRepositoryResName,
		"tf-provider-policy-rule-repository",
		"mysql",
		"http://mysql.local/",
		3306,
	)
	config += formatBasicPolicyIntoConfig(
		"tf-provider-policy-rule-policy",
		[]string{testLabelName},
	)

	config += fmt.Sprintf(`
	resource "cyral_policy_rule" "policy_rule_test" {
		policy_id = cyral_policy.test_policy.id
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
	}`, data.DeletedSeverity, data.DeletedRateLimit, data.ReadSeverity,
		data.ReadRateLimit, data.UpdatedSeverity, data.UpdatedRateLimit)

	return config
}
