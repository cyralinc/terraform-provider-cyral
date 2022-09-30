package cyral

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

const (
	policyRuleResourceName = "policy-rule"
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
				ResourceName:      importStateResName,
			},
		},
	})
}

func TestPolicyRuleResourceUpgradeV0(t *testing.T) {
	previousState := map[string]interface{}{
		"id":        "policy-rule-id",
		"policy_id": "policy-id",
	}
	actualNewState, err := upgradePolicyRuleV0(context.Background(),
		previousState, nil)
	require.NoError(t, err)
	expectedNewState := map[string]interface{}{
		"id":        "policy-id/policy-rule-id",
		"policy_id": "policy-id",
	}
	require.Equal(t, expectedNewState, actualNewState)
}

func setupPolicyRuleTest(policyRule PolicyRuleConfig) (string, resource.TestCheckFunc) {
	testLabelName := "TEST_CCN"
	var config string
	config += formatBasicPolicyIntoConfig(
		accTestName(policyRuleResourceName, "policy"),
		[]string{testLabelName},
	)
	config += formatPolicyRuleConfigIntoConfig(
		policyRule,
		testLabelName,
	)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_policy_rule.policy_rule_test", "deletes.0.severity", policyRule.DeletedSeverity),
		resource.TestCheckResourceAttr("cyral_policy_rule.policy_rule_test", "reads.0.severity", policyRule.ReadSeverity),
		resource.TestCheckResourceAttr("cyral_policy_rule.policy_rule_test", "updates.0.severity", policyRule.UpdatedSeverity),
	)

	return config, testFunction
}

func formatPolicyRuleConfigIntoConfig(
	policyRule PolicyRuleConfig,
	dataLabelName string,
) string {
	return fmt.Sprintf(`
	resource "cyral_policy_rule" "policy_rule_test" {
		policy_id = cyral_policy.test_policy.id
		hosts = ["192.0.2.22", "203.0.113.16/28"]
		identities {
			groups = ["analyst"]
		}
		deletes {
			data = ["%s"]
			rows = 1
			severity = "%s"
			rate_limit = %d
		}
		reads {
			data = ["%s"]
			rows = 1
			severity = "%s"
			rate_limit = %d
		}
		updates {
			data = ["%s"]
			rows = 1
			severity = "%s"
			rate_limit = %d
		}
	}`, dataLabelName, policyRule.DeletedSeverity, policyRule.DeletedRateLimit,
		dataLabelName, policyRule.ReadSeverity, policyRule.ReadRateLimit,
		dataLabelName, policyRule.UpdatedSeverity, policyRule.UpdatedRateLimit)
}
