package cyral

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryUserAccountResourceName = "repository-user-account"
)

func repositoryUserAccountRepositoryConfig() string {
	return utils.FormatBasicRepositoryIntoConfig(
		BasicRepositoryResName,
		utils.AccTestName(repositoryUserAccountResourceName, "main-repo"),
		"mongodb",
		"mongodb.local",
		27017,
	)
}

func TestAccRepositoryUserAccountResource(t *testing.T) {
	// Upgrade test
	onlyRequiredFields := UserAccountResource{
		Name: "name-1",
		AuthScheme: &AuthScheme{
			AWSSecretsManager: &AuthSchemeAWSSecretsManager{
				SecretArn: "secret-arn-1",
			},
		},
	}
	anotherName := UserAccountResource{
		Name: "name-2",
		AuthScheme: &AuthScheme{
			AWSSecretsManager: &AuthSchemeAWSSecretsManager{
				SecretArn: "secret-arn-1",
			},
		},
	}
	anotherAuthScheme := UserAccountResource{
		Name: "name-2",
		AuthScheme: &AuthScheme{
			GCPSecretManager: &AuthSchemeGCPSecretManager{
				SecretName: "secret-name-1",
			},
		},
	}
	withOptionalFields := UserAccountResource{
		Name:             "name-2",
		AuthDatabaseName: "auth-database-name-1",
		Config: &UserAccountConfig{
			Approval: &ApprovalConfig{
				AutomaticGrant:            true,
				MaxAutomaticGrantDuration: "1234s",
			},
		},
		AuthScheme: &AuthScheme{
			GCPSecretManager: &AuthSchemeGCPSecretManager{
				SecretName: "secret-name-1",
			},
		},
	}
	withOptionalFieldsUpdated := UserAccountResource{
		Name:             "name-3",
		AuthDatabaseName: "auth-database-name-2",
		Config: &UserAccountConfig{
			Approval: &ApprovalConfig{
				AutomaticGrant:            true,
				MaxAutomaticGrantDuration: "4321s",
			},
		},
		AuthScheme: &AuthScheme{
			GCPSecretManager: &AuthSchemeGCPSecretManager{
				SecretName: "secret-name-2",
			},
		},
	}
	onlyRequiredFieldsTest := setupRepositoryUserAccountTest(
		"upgrade_test", onlyRequiredFields)
	anotherNameTest := setupRepositoryUserAccountTest(
		"upgrade_test", anotherName)
	anotherAuthSchemeTest := setupRepositoryUserAccountTest(
		"upgrade_test", anotherAuthScheme)
	withOptionalFieldsTest := setupRepositoryUserAccountTest(
		"upgrade_test", withOptionalFields)
	withOptionalFieldsUpdatedTest := setupRepositoryUserAccountTest(
		"upgrade_test", withOptionalFieldsUpdated)

	// Tests covering all auth scheme types
	awsIAM := UserAccountResource{
		Name: "aws-iam-useracc",
		AuthScheme: &AuthScheme{
			AWSIAM: &AuthSchemeAWSIAM{
				RoleARN: "role-arn-1",
			},
		},
	}
	awsSecretsManager := UserAccountResource{
		Name: "aws-sm-useracc",
		AuthScheme: &AuthScheme{
			AWSSecretsManager: &AuthSchemeAWSSecretsManager{
				SecretArn: "secret-arn-1",
			},
		},
	}
	cyralStorage := UserAccountResource{
		Name: "cyral-storage-useracc",
		AuthScheme: &AuthScheme{
			CyralStorage: &AuthSchemeCyralStorage{
				Password: "password-1",
			},
		},
	}
	hashicorpVault := UserAccountResource{
		Name: "hashicorp-vault-useracc",
		AuthScheme: &AuthScheme{
			HashicorpVault: &AuthSchemeHashicorpVault{
				Path:                 "path-1",
				IsDynamicUserAccount: true,
			},
		},
	}
	environmentVariable := UserAccountResource{
		Name: "env-var-useracc",
		AuthScheme: &AuthScheme{
			EnvironmentVariable: &AuthSchemeEnvironmentVariable{
				VariableName: "variable-name-1",
			},
		},
	}
	kubernetesSecret := UserAccountResource{
		Name: "kubesecrets-useracc",
		AuthScheme: &AuthScheme{
			KubernetesSecret: &AuthSchemeKubernetesSecret{
				SecretName: "secret-name-1",
				SecretKey:  "secret-key-1",
			},
		},
	}
	gcpSecretManager := UserAccountResource{
		Name: "gcp-useracc",
		AuthScheme: &AuthScheme{
			GCPSecretManager: &AuthSchemeGCPSecretManager{
				SecretName: "secret-name-1",
			},
		},
	}
	awsIAMTest := setupRepositoryUserAccountTest(
		"aws_iam_test", awsIAM)
	awsSecretsManagerTest := setupRepositoryUserAccountTest(
		"aws_secrets_manager_test", awsSecretsManager)
	cyralStorageTest := setupRepositoryUserAccountTest(
		"cyral_storage_test", cyralStorage)
	hashicorpVaultTest := setupRepositoryUserAccountTest(
		"hashicorp_vault_test", hashicorpVault)
	environmentVariableTest := setupRepositoryUserAccountTest(
		"environment_variable_test", environmentVariable)
	kubernetesSecretTest := setupRepositoryUserAccountTest(
		"kubernetes_secret_test", kubernetesSecret)
	gcpSecretManagerTest := setupRepositoryUserAccountTest(
		"gcp_secret_manager_test", gcpSecretManager)

	// Test with multiple user accounts
	userAccount1ResName := "multiple_accounts_test_1"
	userAccount2ResName := "multiple_accounts_test_2"
	userAccount1 := awsIAM
	userAccount2 := awsSecretsManager
	account1Config := setupRepositoryUserAccountConfig(
		userAccount1ResName, userAccount1)
	account2Config := setupRepositoryUserAccountConfig(
		userAccount2ResName, userAccount2)
	multipleAccountsConfig := repositoryUserAccountRepositoryConfig() +
		account1Config + account2Config
	multipleAccountsCheck1 := setupRepositoryUserAccountCheck(
		userAccount1ResName, userAccount1)
	multipleAccountsCheck2 := setupRepositoryUserAccountCheck(
		userAccount2ResName, userAccount2)
	multipleAccountsCheck := resource.ComposeTestCheckFunc(
		multipleAccountsCheck1, multipleAccountsCheck2)
	multipleAccountsTest := resource.TestStep{
		Config: multipleAccountsConfig,
		Check:  multipleAccountsCheck,
	}

	// Import test
	importResName := fmt.Sprintf("cyral_repository_user_account.%s",
		userAccount2ResName)
	importTest := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      importResName,
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Update tests
			onlyRequiredFieldsTest,
			anotherNameTest,
			anotherAuthSchemeTest,
			withOptionalFieldsTest,
			withOptionalFieldsUpdatedTest,

			// Tests covering all auth scheme types
			awsIAMTest,
			awsSecretsManagerTest,
			cyralStorageTest,
			hashicorpVaultTest,
			environmentVariableTest,
			kubernetesSecretTest,
			gcpSecretManagerTest,

			// Test with multiple user accounts
			multipleAccountsTest,

			// Import test
			importTest,
		},
	})
}

func setupRepositoryUserAccountTest(resName string, userAccount UserAccountResource) resource.TestStep {
	return resource.TestStep{
		Config: repositoryUserAccountRepositoryConfig() +
			setupRepositoryUserAccountConfig(resName, userAccount),
		Check: setupRepositoryUserAccountCheck(resName, userAccount),
	}
}

func setupRepositoryUserAccountCheck(resName string, userAccount UserAccountResource) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_repository_user_account.%s", resName)

	var checkFuncs []resource.TestCheckFunc

	// Required attributes
	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			resFullName, "repository_id",
			fmt.Sprintf("cyral_repository.%s", BasicRepositoryResName), "id"),
		resource.TestCheckResourceAttr(resFullName,
			"name", userAccount.Name),
		resource.TestCheckResourceAttr(resFullName,
			"auth_scheme.#", "1"),
	}...)

	// Optional attributes
	checkFuncs = append(checkFuncs,
		resource.TestCheckResourceAttr(resFullName,
			"auth_database_name", userAccount.AuthDatabaseName),
	)
	if userAccount.Config != nil && userAccount.Config.Approval != nil {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resFullName,
				"approval_config.#", "1"),
			resource.TestCheckResourceAttr(resFullName,
				"approval_config.0.automatic_grant",
				strconv.FormatBool(userAccount.Config.Approval.AutomaticGrant)),
			resource.TestCheckResourceAttr(resFullName,
				"approval_config.0.max_auto_grant_duration",
				userAccount.Config.Approval.MaxAutomaticGrantDuration),
		}...)
	}

	// Auth scheme
	authSchemeScope := "auth_scheme.0."
	switch authScheme := userAccount.AuthScheme; {
	case authScheme.AWSIAM != nil:
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr(resFullName,
				authSchemeScope+"aws_iam.0.role_arn",
				authScheme.AWSIAM.RoleARN))
	case authScheme.AWSSecretsManager != nil:
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr(resFullName,
				authSchemeScope+"aws_secrets_manager.0.secret_arn",
				authScheme.AWSSecretsManager.SecretArn))
	case authScheme.CyralStorage != nil:
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr(resFullName,
				authSchemeScope+"cyral_storage.0.password",
				authScheme.CyralStorage.Password))
	case authScheme.EnvironmentVariable != nil:
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr(resFullName,
				authSchemeScope+"environment_variable.0.variable_name",
				authScheme.EnvironmentVariable.VariableName))
	case authScheme.KubernetesSecret != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resFullName,
				authSchemeScope+"kubernetes_secret.0.secret_name",
				authScheme.KubernetesSecret.SecretName,
			),
			resource.TestCheckResourceAttr(resFullName,
				authSchemeScope+"kubernetes_secret.0.secret_key",
				authScheme.KubernetesSecret.SecretKey,
			),
		}...)
	case authScheme.GCPSecretManager != nil:
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr(resFullName,
				authSchemeScope+"gcp_secrets_manager.0.secret_name",
				authScheme.GCPSecretManager.SecretName))
	case authScheme.HashicorpVault != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resFullName,
				authSchemeScope+"hashicorp_vault.0.path",
				authScheme.HashicorpVault.Path,
			),
			resource.TestCheckResourceAttr(resFullName,
				authSchemeScope+"hashicorp_vault.0.is_dynamic_user_account",
				strconv.FormatBool(authScheme.HashicorpVault.IsDynamicUserAccount),
			),
		}...)
	}

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func setupRepositoryUserAccountConfig(resName string, userAccount UserAccountResource) string {
	var config string

	var authSchemeStr string
	switch authScheme := userAccount.AuthScheme; {
	case authScheme.AWSIAM != nil:
		authSchemeStr = fmt.Sprintf(`
			aws_iam {
				role_arn = "%s"
			}`, authScheme.AWSIAM.RoleARN)
	case authScheme.AWSSecretsManager != nil:
		authSchemeStr = fmt.Sprintf(`
			aws_secrets_manager {
				secret_arn = "%s"
			}`, authScheme.AWSSecretsManager.SecretArn)
	case authScheme.CyralStorage != nil:
		authSchemeStr = fmt.Sprintf(`
			cyral_storage {
				password = "%s"
			}`, authScheme.CyralStorage.Password)
	case authScheme.EnvironmentVariable != nil:
		authSchemeStr = fmt.Sprintf(`
			environment_variable {
				variable_name = "%s"
			}`, authScheme.EnvironmentVariable.VariableName)
	case authScheme.KubernetesSecret != nil:
		authSchemeStr = fmt.Sprintf(`
			kubernetes_secret {
				secret_name = "%s"
				secret_key = "%s"
			}`, authScheme.KubernetesSecret.SecretName,
			authScheme.KubernetesSecret.SecretKey)
	case authScheme.GCPSecretManager != nil:
		authSchemeStr = fmt.Sprintf(`
			gcp_secrets_manager {
				secret_name = "%s"
			}`, authScheme.GCPSecretManager.SecretName)
	case authScheme.HashicorpVault != nil:
		authSchemeStr = fmt.Sprintf(`
			hashicorp_vault {
				path = "%s"
				is_dynamic_user_account = %t
			}`, authScheme.HashicorpVault.Path,
			authScheme.HashicorpVault.IsDynamicUserAccount)
	}

	var approvalConfigStr string
	if userAccount.Config != nil && userAccount.Config.Approval != nil {
		approvalConfigStr = fmt.Sprintf(`
		approval_config {
			automatic_grant = %t
			max_auto_grant_duration = "%s"
		}`, userAccount.Config.Approval.AutomaticGrant,
			userAccount.Config.Approval.MaxAutomaticGrantDuration)
	}

	config += fmt.Sprintf(`
	resource "cyral_repository_user_account" "%s" {
		repository_id = %s
		name = "%s"
		auth_database_name = "%s"
		auth_scheme {
			%s
		}
		%s
	}`, resName, BasicRepositoryID, userAccount.Name,
		userAccount.AuthDatabaseName, authSchemeStr, approvalConfigStr)

	return config
}
