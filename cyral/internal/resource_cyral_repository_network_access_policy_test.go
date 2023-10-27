package internal_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryNetworkAccessPolicyResourceName = "repository-network-access-policy"
)

func TestAccRepositoryNetworkAccessPolicyResource(t *testing.T) {
	// In these tests, the resources should be recreated every time.
	emptyFieldsNotEnabled := &internal.NetworkAccessPolicy{
		Enabled: false,
		NetworkAccessRules: internal.NetworkAccessRules{
			RulesBlockAccess: false,
			Rules:            []internal.NetworkAccessRule{},
		},
	}
	emptyFieldsExceptName := &internal.NetworkAccessPolicy{
		NetworkAccessRules: internal.NetworkAccessRules{
			Rules: []internal.NetworkAccessRule{
				{
					Name: "name-1",
				},
			},
		},
	}
	emptyFieldsNotEnabledTest := setupRepositoryNetworkAccessPolicyTest(
		"empty_fields_not_enabled", emptyFieldsNotEnabled, nil)
	emptyFieldsExceptNameTest := setupRepositoryNetworkAccessPolicyTest(
		"empty_fields_except_name", emptyFieldsExceptName, nil)

	// In these tests, the resources should be updated in-place
	fullRule := &internal.NetworkAccessPolicy{
		NetworkAccessRules: internal.NetworkAccessRules{
			Rules: []internal.NetworkAccessRule{
				{
					Name:        "name-2",
					Description: "description-2",
					DBAccounts:  []string{"dbaccount-2"},
					SourceIPs:   []string{"1.2.3.4", "5.6.7.8"},
				},
			},
		},
	}
	fullRuleUpdated := &internal.NetworkAccessPolicy{
		NetworkAccessRules: internal.NetworkAccessRules{
			Rules: []internal.NetworkAccessRule{
				{
					Name:        "name-3",
					Description: "description-3",
					DBAccounts:  []string{"dbaccount-3"},
					SourceIPs:   []string{"9.10.11.12", "13.14.15.16"},
				},
			},
		},
	}
	fullRuleUpdatedAdditionalRule := &internal.NetworkAccessPolicy{
		NetworkAccessRules: internal.NetworkAccessRules{
			Rules: []internal.NetworkAccessRule{
				{
					Name:        "name-3",
					Description: "description-3",
					DBAccounts:  []string{"dbaccount-3"},
					SourceIPs:   []string{"9.10.11.12", "13.14.15.16"},
				},
				{
					Name:        "name-4",
					Description: "description-4",
					DBAccounts:  []string{"dbaccount-4"},
					SourceIPs:   []string{"17.18.19.20", "21.22.23.24"},
				},
			},
		},
	}
	fullRuleTest := setupRepositoryNetworkAccessPolicyTest(
		"full_rule", fullRule, []string{"dbaccount-2"})
	fullRuleUpdatedTest := setupRepositoryNetworkAccessPolicyTest(
		"full_rule", fullRuleUpdated, []string{"dbaccount-3"})
	fullRuleUpdatedAdditionalRuleTest := setupRepositoryNetworkAccessPolicyTest(
		"full_rule", fullRuleUpdatedAdditionalRule, []string{"dbaccount-3", "dbaccount-4"})

	importResourceName := "cyral_repository_network_access_policy.full_rule"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Recreation tests
			emptyFieldsNotEnabledTest,
			emptyFieldsExceptNameTest,

			// Update tests
			fullRuleTest,
			fullRuleUpdatedTest,
			fullRuleUpdatedAdditionalRuleTest,

			// Import tests
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      importResourceName,
			},
		},
	})
}

func setupRepositoryNetworkAccessPolicyTest(
	resName string,
	nap *internal.NetworkAccessPolicy,
	dbAccounts []string,
) resource.TestStep {
	return resource.TestStep{
		Config: setupRepositoryNetworkAccessPolicyConfig(resName, nap, dbAccounts),
		Check:  setupRepositoryNetworkAccessPolicyCheck(resName, nap),
	}
}

func setupRepositoryNetworkAccessPolicyConfig(
	resName string,
	nap *internal.NetworkAccessPolicy,
	dbAccounts []string,
) string {
	var config string

	// Repository
	repoResName := "test_repo"
	repoID := fmt.Sprintf("cyral_repository.%s.id", repoResName)
	config += utils.FormatBasicRepositoryIntoConfig(
		repoResName,
		utils.AccTestName(repositoryNetworkAccessPolicyResourceName, resName),
		"sqlserver",
		"my.host.com",
		1433,
	)

	// Network Access Policy
	config += formatNetworkAccessPolicyIntoConfig(resName, repoID, nap, nil)

	return config
}

func setupRepositoryNetworkAccessPolicyCheck(
	resName string,
	nap *internal.NetworkAccessPolicy,
) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_repository_network_access_policy.%s", resName)

	testFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resFullName,
			"network_access_rule.#",
			fmt.Sprintf("%d", len(nap.NetworkAccessRules.Rules)),
		),
		resource.TestCheckResourceAttr(resFullName,
			"enabled", strconv.FormatBool(nap.Enabled),
		),
		resource.TestCheckResourceAttr(resFullName,
			"network_access_rules_block_access", strconv.FormatBool(nap.NetworkAccessRules.RulesBlockAccess),
		),
	}

	for i, accessRule := range nap.NetworkAccessRules.Rules {
		testFuncs = append(testFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resFullName,
				fmt.Sprintf("network_access_rule.%d.name", i),
				accessRule.Name),
			resource.TestCheckResourceAttr(resFullName,
				fmt.Sprintf("network_access_rule.%d.description", i),
				accessRule.Description),
		}...)

		for j, dbAccount := range accessRule.DBAccounts {
			testFuncs = append(testFuncs, resource.TestCheckResourceAttr(resFullName,
				fmt.Sprintf("network_access_rule.%d.db_accounts.%d", i, j),
				dbAccount),
			)
		}
		for j, sourceIP := range accessRule.SourceIPs {
			testFuncs = append(testFuncs, resource.TestCheckResourceAttr(resFullName,
				fmt.Sprintf("network_access_rule.%d.source_ips.%d", i, j),
				sourceIP),
			)
		}
	}

	return resource.ComposeTestCheckFunc(testFuncs...)
}

func formatNetworkAccessPolicyIntoConfig(
	resName, repositoryID string, nap *internal.NetworkAccessPolicy, dependsOn []string,
) string {
	var narStr string
	for _, nar := range nap.Rules {
		narStr += fmt.Sprintf(`
		network_access_rule {
			name = "%s"
			description = "%s"
			db_accounts = %s
			source_ips = %s
		}`, nar.Name, nar.Description, utils.ListToStr(nar.DBAccounts),
			utils.ListToStr(nar.SourceIPs))
	}

	config := fmt.Sprintf(`
	resource "cyral_repository_network_access_policy" "%s" {
		repository_id = %s
		enabled = %t
		network_access_rules_block_access = %t
		%s
		depends_on = %s
	}`, resName, repositoryID, nap.Enabled, nap.NetworkAccessRules.RulesBlockAccess,
		narStr, utils.ListToStrNoQuotes(dependsOn))

	return config
}
