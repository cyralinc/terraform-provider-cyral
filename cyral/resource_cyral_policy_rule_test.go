package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	Identities       *Identity
}

var initialPolicyRuleConfig PolicyRuleConfig = PolicyRuleConfig{
	DeletedSeverity:  "medium",
	UpdatedSeverity:  "high",
	ReadSeverity:     "low",
	DeletedRateLimit: 1,
	UpdatedRateLimit: 2,
	ReadRateLimit:    3,
	Identities: &Identity{
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
	Identities: &Identity{
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
	Identities: &Identity{
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
	Identities: &Identity{
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
		ProviderFactories: providerFactories,
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
			db_roles = %s`, listToStr(identities.DBRoles))
		}
		if identities.Groups != nil {
			identitiesStr += fmt.Sprintf(`
			groups = %s`, listToStr(identities.Groups))
		}
		if identities.Services != nil {
			identitiesStr += fmt.Sprintf(`
			services = %s`, listToStr(identities.Services))
		}
		if identities.Users != nil {
			identitiesStr += fmt.Sprintf(`
			users = %s`, listToStr(identities.Users))
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
