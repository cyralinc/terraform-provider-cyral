package rule_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/deprecated/policy/rule"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
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
	Identities       *rule.Identity
}

var initialPolicyRuleConfig PolicyRuleConfig = PolicyRuleConfig{
	DeletedSeverity:  "medium",
	UpdatedSeverity:  "high",
	ReadSeverity:     "low",
	DeletedRateLimit: 1,
	UpdatedRateLimit: 2,
	ReadRateLimit:    3,
	Identities: &rule.Identity{
		DBRoles: []string{
			"db-role-1",
		},
	},
}

var updatedGroupsPolicyRuleConfig PolicyRuleConfig = PolicyRuleConfig{
	DeletedSeverity:  "low",
	UpdatedSeverity:  "medium",
	ReadSeverity:     "high",
	DeletedRateLimit: 2,
	UpdatedRateLimit: 3,
	ReadRateLimit:    4,
	Identities: &rule.Identity{
		Groups: []string{
			"groups-1",
		},
	},
}

var updatedServicesPolicyRuleConfig PolicyRuleConfig = PolicyRuleConfig{
	DeletedSeverity:  "high",
	UpdatedSeverity:  "low",
	ReadSeverity:     "medium",
	DeletedRateLimit: 5,
	UpdatedRateLimit: 6,
	ReadRateLimit:    7,
	Identities: &rule.Identity{
		Services: []string{
			"services-1",
		},
	},
}

var updatedUsersPolicyRuleConfig PolicyRuleConfig = PolicyRuleConfig{
	DeletedSeverity:  "high",
	UpdatedSeverity:  "low",
	ReadSeverity:     "medium",
	DeletedRateLimit: 5,
	UpdatedRateLimit: 6,
	ReadRateLimit:    7,
	Identities: &rule.Identity{
		Users: []string{
			"users-1",
		},
	},
}

var updatedNoIdentityPolicyRuleConfig PolicyRuleConfig = PolicyRuleConfig{
	DeletedSeverity:  "low",
	UpdatedSeverity:  "medium",
	ReadSeverity:     "high",
	DeletedRateLimit: 1,
	UpdatedRateLimit: 2,
	ReadRateLimit:    3,
}

func TestAccPolicyRuleResource(t *testing.T) {
	testConfig, testFunc := setupPolicyRuleTest(initialPolicyRuleConfig)
	testGroupsConfig, testGroupsFunc := setupPolicyRuleTest(updatedGroupsPolicyRuleConfig)
	testServicesConfig, testServicesFunc := setupPolicyRuleTest(updatedServicesPolicyRuleConfig)
	testUsersConfig, testUsersFunc := setupPolicyRuleTest(updatedUsersPolicyRuleConfig)
	testNoIdentityConfig, testNoIdentityFunc := setupPolicyRuleTest(updatedNoIdentityPolicyRuleConfig)

	importStateResName := "cyral_policy_rule.policy_rule_test"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testGroupsConfig,
				Check:  testGroupsFunc,
			},
			{
				Config: testServicesConfig,
				Check:  testServicesFunc,
			},
			{
				Config: testUsersConfig,
				Check:  testUsersFunc,
			},
			{
				Config: testNoIdentityConfig,
				Check:  testNoIdentityFunc,
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
	actualNewState, err := rule.UpgradePolicyRuleV0(context.Background(),
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
	config += utils.FormatBasicPolicyIntoConfig(
		utils.AccTestName(policyRuleResourceName, "policy"),
		[]string{testLabelName},
	)
	config += formatPolicyRuleConfigIntoConfig(
		policyRule,
		testLabelName,
	)

	resFullName := "cyral_policy_rule.policy_rule_test"
	testFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resFullName,
			"deletes.0.severity", policyRule.DeletedSeverity),
		resource.TestCheckResourceAttr(resFullName,
			"reads.0.severity", policyRule.ReadSeverity),
		resource.TestCheckResourceAttr(resFullName,
			"updates.0.severity", policyRule.UpdatedSeverity),
	}

	if policyRule.Identities != nil {
		identities := policyRule.Identities

		for i, dbRole := range identities.DBRoles {
			testFunctions = append(testFunctions,
				resource.TestCheckResourceAttr(resFullName,
					fmt.Sprintf("identities.0.db_roles.%d", i), dbRole),
			)
		}
		for i, group := range identities.Groups {
			testFunctions = append(testFunctions,
				resource.TestCheckResourceAttr(resFullName,
					fmt.Sprintf("identities.0.groups.%d", i), group),
			)
		}
		for i, service := range identities.Services {
			testFunctions = append(testFunctions,
				resource.TestCheckResourceAttr(resFullName,
					fmt.Sprintf("identities.0.services.%d", i), service),
			)
		}
		for i, user := range identities.Users {
			testFunctions = append(testFunctions,
				resource.TestCheckResourceAttr(resFullName,
					fmt.Sprintf("identities.0.users.%d", i), user),
			)
		}
	}

	return config, resource.ComposeTestCheckFunc(testFunctions...)
}

func formatPolicyRuleConfigIntoConfig(
	policyRule PolicyRuleConfig,
	dataLabelName string,
) string {
	var identitiesStr string
	if policyRule.Identities != nil {
		identities := policyRule.Identities
		identitiesStr += `
		identities {`

		if identities.DBRoles != nil {
			identitiesStr += fmt.Sprintf(`
			db_roles = %s`, utils.ListToStr(identities.DBRoles))
		}
		if identities.Groups != nil {
			identitiesStr += fmt.Sprintf(`
			groups = %s`, utils.ListToStr(identities.Groups))
		}
		if identities.Services != nil {
			identitiesStr += fmt.Sprintf(`
			services = %s`, utils.ListToStr(identities.Services))
		}
		if identities.Users != nil {
			identitiesStr += fmt.Sprintf(`
			users = %s`, utils.ListToStr(identities.Users))
		}

		identitiesStr += `
		}`
	}

	return fmt.Sprintf(`
	resource "cyral_policy_rule" "policy_rule_test" {
		policy_id = cyral_policy.test_policy.id
		hosts = ["192.0.2.22", "203.0.113.16/28"]
		%s
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
	}`, identitiesStr,
		dataLabelName, policyRule.DeletedSeverity, policyRule.DeletedRateLimit,
		dataLabelName, policyRule.ReadSeverity, policyRule.ReadRateLimit,
		dataLabelName, policyRule.UpdatedSeverity, policyRule.UpdatedRateLimit)
}
