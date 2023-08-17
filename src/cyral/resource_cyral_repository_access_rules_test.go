package cyral

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryAccessRulesResourceName = "repository-access-rules"
)

var validFrom string = "2022-01-02T10:20:30Z"
var validUntil string = "3022-01-02T10:20:30Z"
var validFromUpdated string = "2023-11-12T10:20:30Z"
var validUntilUpdated string = "3023-11-12T10:20:30Z"

var initialAccessRulesConfig *AccessRulesResource = &AccessRulesResource{
	AccessRules: []*AccessRule{
		{
			Identity: &AccessRulesIdentity{
				Name: "identityEmail",
				Type: "email",
			},
			ValidFrom:  &validFrom,
			ValidUntil: &validUntil,
			Config: &AccessRulesConfig{
				AuthorizationPolicyInstanceIDs: []string{
					"policy1",
					"policy2",
				},
			},
		},
		{
			Identity: &AccessRulesIdentity{
				Name: "identityGroup",
				Type: "group",
			},
			ValidFrom:  &validFrom,
			ValidUntil: &validUntil,
			Config: &AccessRulesConfig{
				AuthorizationPolicyInstanceIDs: []string{
					"policy3",
					"policy4",
				},
			},
		},
	},
}

// Let's modify the identity names, the durations, and the policy IDs
var updatedAccessRulesConfig *AccessRulesResource = &AccessRulesResource{
	AccessRules: []*AccessRule{
		{
			Identity: &AccessRulesIdentity{
				Name: "identityEmailUpdated",
				Type: "email",
			},
			ValidFrom:  &validFromUpdated,
			ValidUntil: &validUntilUpdated,
			Config: &AccessRulesConfig{
				AuthorizationPolicyInstanceIDs: []string{
					"policy1Updated",
					"policy2Updated",
				},
			},
		},
		{
			Identity: &AccessRulesIdentity{
				Name: "identityGroupUpdated",
				Type: "group",
			},
			ValidFrom:  &validFromUpdated,
			ValidUntil: &validUntilUpdated,
			Config: &AccessRulesConfig{
				AuthorizationPolicyInstanceIDs: []string{
					"policy3Updated",
					"policy4Updated",
				},
			},
		},
	},
}

var barebonesAccessRulesConfig *AccessRulesResource = &AccessRulesResource{
	AccessRules: []*AccessRule{
		{
			Identity: &AccessRulesIdentity{
				Name: "identityUsername",
				Type: "username",
			},
		},
	},
}

func TestAccRepositoryAccessRulesResource(t *testing.T) {
	testConfig, testFunc := setupRepositoryAccessRulesTest(
		initialAccessRulesConfig,
		false,
	)

	updatedConfig, updatedFunc := setupRepositoryAccessRulesTest(
		updatedAccessRulesConfig,
		false,
	)

	barebonesConfig, barebonesFunc := setupRepositoryAccessRulesTest(
		barebonesAccessRulesConfig,
		true,
	)

	importStateResName := "cyral_repository_access_rules.acc_test_access_rules"

	resource.ParallelTest(
		t,
		resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: testConfig,
					Check:  testFunc,
				},
				{
					Config: updatedConfig,
					Check:  updatedFunc,
				},
				{
					Config: barebonesConfig,
					Check:  barebonesFunc,
				},
				{
					ImportState:       true,
					ImportStateVerify: true,
					ResourceName:      importStateResName,
				},
			},
		},
	)
}

func setupRepositoryAccessRulesTest(
	accessRulesData *AccessRulesResource,
	bareBones bool,
) (string, resource.TestCheckFunc) {
	var configuration string
	configuration += utils.FormatBasicRepositoryIntoConfig(
		BasicRepositoryResName,
		utils.AccTestName(repositoryAccessRulesResourceName, "repository1"),
		"mongodb",
		"mongo.local",
		3333,
	) + "\n"

	configuration += userAccConfig(
		BasicRepositoryID,
		utils.AccTestName(repositoryAccessRulesResourceName, "user_acount"),
	) + "\n"

	configuration += accessRulesToConfig(
		accessRulesData,
		BasicRepositoryID,
		"cyral_repository_user_account.my_test_user_account.user_account_id",
		bareBones,
	) + "\n"

	var testFunction resource.TestCheckFunc
	if bareBones {
		testFunction = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.#",
				"1",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.identity.#",
				"1",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.identity.0.type",
				"username",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.identity.0.name",
				"identityUsername",
			),
		)
	} else {
		testFunction = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.#",
				"2",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.identity.#",
				"1",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.1.identity.#",
				"1",
			),

			// Checks for rule 0
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.identity.0.type",
				"email",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.identity.0.name",
				fmt.Sprintf("%s", accessRulesData.AccessRules[0].Identity.Name),
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.valid_from",
				fmt.Sprintf("%s", *accessRulesData.AccessRules[0].ValidFrom),
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.valid_until",
				fmt.Sprintf("%s", *accessRulesData.AccessRules[0].ValidUntil),
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.config.#",
				"1",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.config.0.policy_ids.#",
				"2",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.config.0.policy_ids.0",
				fmt.Sprintf(
					"%s",
					accessRulesData.AccessRules[0].Config.AuthorizationPolicyInstanceIDs[0],
				),
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.0.config.0.policy_ids.1",
				fmt.Sprintf(
					"%s",
					accessRulesData.AccessRules[0].Config.AuthorizationPolicyInstanceIDs[1],
				),
			),

			// Checks for rule 1
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.1.identity.0.type",
				"group",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.1.identity.0.name",
				fmt.Sprintf("%s", accessRulesData.AccessRules[1].Identity.Name),
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.1.valid_from",
				fmt.Sprintf("%s", *accessRulesData.AccessRules[0].ValidFrom),
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.1.valid_until",
				fmt.Sprintf("%s", *accessRulesData.AccessRules[0].ValidUntil),
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.1.config.#",
				"1",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.1.config.0.policy_ids.#",
				"2",
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.1.config.0.policy_ids.0",
				fmt.Sprintf(
					"%s",
					accessRulesData.AccessRules[1].Config.AuthorizationPolicyInstanceIDs[0],
				),
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_access_rules.acc_test_access_rules",
				"rule.1.config.0.policy_ids.1",
				fmt.Sprintf(
					"%s",
					accessRulesData.AccessRules[1].Config.AuthorizationPolicyInstanceIDs[1],
				),
			),
		)
	}

	return configuration, testFunction
}

func userAccConfig(
	repoID string,
	name string,
) string {
	return fmt.Sprintf(` resource "cyral_repository_user_account" "my_test_user_account" {
		repository_id = %s
		name = "%s"
		auth_scheme {
			environment_variable {
				variable_name = "FOOBAR"
			}
		}
	}`,
		repoID,
		name,
	)
}

func accessRulesToConfig(
	res *AccessRulesResource,
	repoID string,
	userAccountID string,
	bareBones bool,
) string {
	rules := res.AccessRules
	if !bareBones {
		return fmt.Sprintf(`
		resource "cyral_repository_access_rules" "acc_test_access_rules" {
			repository_id = %s
			user_account_id = %s

			rule {
				identity {
					type = "%s"
					name = "%s"
				}
				config {
					policy_ids = [
						"%s",
						"%s",
					]
				}
				valid_from = "%s"
				valid_until = "%s"
			}

			rule {
				identity {
					type = "%s"
					name = "%s"
				}
				config {
					policy_ids = [
						"%s",
						"%s",
					]
				}
				valid_from = "%s"
				valid_until = "%s"
			}
		}`,
			repoID,
			userAccountID,
			rules[0].Identity.Type,
			rules[0].Identity.Name,
			rules[0].Config.AuthorizationPolicyInstanceIDs[0],
			rules[0].Config.AuthorizationPolicyInstanceIDs[1],
			*rules[0].ValidFrom,
			*rules[0].ValidUntil,
			rules[1].Identity.Type,
			rules[1].Identity.Name,
			rules[1].Config.AuthorizationPolicyInstanceIDs[0],
			rules[1].Config.AuthorizationPolicyInstanceIDs[1],
			*rules[1].ValidFrom,
			*rules[1].ValidUntil,
		)
	}

	return fmt.Sprintf(`
		resource "cyral_repository_access_rules" "acc_test_access_rules" {
			repository_id = %s
			user_account_id = %s

			rule {
				identity {
					type = "%s"
					name = "%s"
				}
			}
		}`,
		repoID,
		userAccountID,
		rules[0].Identity.Type,
		rules[0].Identity.Name,
	)
}
