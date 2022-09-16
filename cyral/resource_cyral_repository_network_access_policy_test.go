package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryNetworkAccessPolicyResourceName = "repository-network-access-policy"
)

func TestAccRepositoryNetworkAccessPolicy(t *testing.T) {
	// Recreate these resources each step
	emptyFields := &NetworkAccessPolicy{
		NetworkAccessRules: NetworkAccessRules{
			Rules: []NetworkAccessRule{},
		},
	}
	emptyFieldsExceptName := &NetworkAccessPolicy{
		NetworkAccessRules: NetworkAccessRules{
			Rules: []NetworkAccessRule{
				{
					Name: "name-1",
				},
			},
		},
	}
	emptyFieldsTest := setupRepositoryNetworkAccessPolicyTest(
		"initial", emptyFields, nil)
	emptyFieldsExceptNameTest := setupRepositoryNetworkAccessPolicyTest(
		"initial", emptyFieldsExceptName, nil)

	// Update these resources. These resources depend on
	// repository_local_account being in the config.
	fullRule := &NetworkAccessPolicy{
		NetworkAccessRules: NetworkAccessRules{
			Rules: []NetworkAccessRule{
				{
					Name:        "name-2",
					Description: "description-2",
					DBAccounts:  []string{"dbaccount-2"},
					SourceIPs:   []string{"1.2.3.4", "5.6.7.8"},
				},
			},
		},
	}
	fullRuleUpdated := &NetworkAccessPolicy{
		NetworkAccessRules: NetworkAccessRules{
			Rules: []NetworkAccessRule{
				{
					Name:        "name-3",
					Description: "description-3",
					DBAccounts:  []string{"dbaccount-3"},
					SourceIPs:   []string{"9.10.11.12", "13.14.15.16"},
				},
			},
		},
	}
	fullRuleUpdatedAdditionalRule := &NetworkAccessPolicy{
		NetworkAccessRules: NetworkAccessRules{
			Rules: []NetworkAccessRule{
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
			emptyFieldsTest,
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
	nap *NetworkAccessPolicy,
	dbAccounts []string,
) resource.TestStep {
	return resource.TestStep{
		Config: setupRepositoryNetworkAccessPolicyConfig(resName, nap, dbAccounts),
		Check:  setupRepositoryNetworkAccessPolicyCheck(resName, nap),
	}
}

func setupRepositoryNetworkAccessPolicyConfig(
	resName string,
	nap *NetworkAccessPolicy,
	dbAccounts []string,
) string {
	var config string

	// Repository
	repoResName := "test_repo"
	repoID := fmt.Sprintf("cyral_repository.%s.id", repoResName)
	config += formatBasicRepositoryIntoConfig(
		repoResName,
		accTestName(repositoryNetworkAccessPolicyResourceName, "repo-name"),
		"sqlserver",
		"my.host.com",
		1433,
	)

	// Local accounts
	config += sampleMultipleBasicRepositoryLocalAccountIntoConfig(
		repoID, dbAccounts)

	// Repo Conf Auth
	config += formatRepositoryConfAuthDataIntoConfig(
		resName,
		RepositoryConfAuthData{
			EnableNetworkAccessControl: true,
		},
		repoID,
	)

	// Network Access Policy
	config += formatNetworkAccessPolicyIntoConfig(resName, repoID, nap)

	return config
}

func setupRepositoryNetworkAccessPolicyCheck(
	resName string,
	nap *NetworkAccessPolicy,
) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_repository_network_access_policy.%s", resName)

	testFuncs := []resource.TestCheckFunc{}

	testFuncs = append(testFuncs, resource.TestCheckResourceAttr(resFullName,
		"network_access_rule.#",
		fmt.Sprintf("%d", len(nap.NetworkAccessRules.Rules)),
	))

	// for i, accessRule := range nap.NetworkAccessRules.Rules {
	// 	testFuncs = append(testFuncs, []resource.TestCheckFunc{
	// 		resource.TestCheckAttr(resFullName,
	// 			"network_access_rule.0")
	// 	}...)
	// }

	return resource.ComposeTestCheckFunc(testFuncs...)
}

func formatNetworkAccessPolicyIntoConfig(
	resName, repositoryID string, nap *NetworkAccessPolicy,
) string {
	var narStr string
	for _, nar := range nap.Rules {
		narStr += fmt.Sprintf(`
		network_access_rule {
			name = "%s"
			description = "%s"
			db_accounts = %s
			source_ips = %s
		}`, nar.Name, nar.Description, listToStr(nar.DBAccounts),
			listToStr(nar.SourceIPs))
	}

	config := fmt.Sprintf(`
	resource "cyral_repository_network_access_policy" "%s" {
		repository_id = %s
		%s
	}`, resName, repositoryID, narStr)

	return config
}
