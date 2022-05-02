package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialRepoAccountConfigAwsIAM = RepositoryLocalAccountResource{
	AwsIAM: &AwsIAMResource{
		DatabaseName: "tf-test-db-name",
		RepoAccount:  "tf-test-repo-account",
		RoleArn:      "tf-test-role-arn",
	},
}

var updatedRepoAccountConfigAwsIAM = RepositoryLocalAccountResource{
	AwsIAM: &AwsIAMResource{
		DatabaseName: "tf-test-db-name-updated",
		RepoAccount:  "tf-test-repo-account-updated",
		RoleArn:      "tf-test-role-arn-updated",
	},
}

var initialRepoAccountConfigAwsSecretsManager = RepositoryLocalAccountResource{
	AwsSecretsManager: &AwsSecretsManagerResource{
		DatabaseName: "tf-test-db-name",
		RepoAccount:  "tf-test-repo-account",
		SecretArn:    "tf-test-secret-arn",
	},
}

var updatedRepoAccountConfigAwsSecretsManager = RepositoryLocalAccountResource{
	AwsSecretsManager: &AwsSecretsManagerResource{
		DatabaseName: "tf-test-db-name-updated",
		RepoAccount:  "tf-test-repo-account-updated",
		SecretArn:    "tf-test-secret-arn-updated",
	},
}

var initialRepoAccountConfigCyralStorage = RepositoryLocalAccountResource{
	CyralStorage: &CyralStorageResource{
		DatabaseName: "tf-test-db-name",
		RepoAccount:  "tf-test-repo-account",
		Password:     "tf-test-pasword",
	},
}

var updatedRepoAccountConfigCyralStorage = RepositoryLocalAccountResource{
	CyralStorage: &CyralStorageResource{
		DatabaseName: "tf-test-db-name-updated",
		RepoAccount:  "tf-test-repo-account-updated",
		Password:     "tf-test-pasword-updated",
	},
}

var initialRepoAccountConfigHashicorpVault = RepositoryLocalAccountResource{
	HashicorpVault: &HashicorpVaultResource{
		DatabaseName: "tf-test-db-name",
		RepoAccount:  "tf-test-repo-account",
		Path:         "tf-test-path",
	},
}

var updatedRepoAccountConfigHashicorpVault = RepositoryLocalAccountResource{
	HashicorpVault: &HashicorpVaultResource{
		DatabaseName: "tf-test-db-name-updated",
		RepoAccount:  "tf-test-repo-account-updated",
		Path:         "tf-test-path-updated",
	},
}

var initialRepoAccountConfigEnvironmentVariable = RepositoryLocalAccountResource{
	EnvironmentVariable: &EnvironmentVariableResource{
		DatabaseName: "tf-test-db-name",
		RepoAccount:  "tf-test-repo-account",
		VariableName: "CYRAL_DBSECRETS_TF_TEST_VARIABLE_NAME",
	},
}

var updatedRepoAccountConfigEnvironmentVariable = RepositoryLocalAccountResource{
	EnvironmentVariable: &EnvironmentVariableResource{
		DatabaseName: "tf-test-db-name-updated",
		RepoAccount:  "tf-test-repo-account-updated",
		VariableName: "CYRAL_DBSECRETS_TF_TEST_VARIABLE_NAME_UPDATED",
	},
}

var initialRepoAccountConfigKubernetesSecret = RepositoryLocalAccountResource{
	KubernetesSecret: &KubernetesSecretResource{
		DatabaseName: "tf-test-db-name",
		RepoAccount:  "tf-test-repo-account",
		SecretName:   "tf-test-db-secrets",
		SecretKey:    "tf-test-secret-key",
	},
}

var updatedRepoAccountConfigKubernetesSecret = RepositoryLocalAccountResource{
	KubernetesSecret: &KubernetesSecretResource{
		DatabaseName: "tf-test-db-name-updated",
		RepoAccount:  "tf-test-repo-account-updated",
		SecretName:   "db-secrets-updated",
		SecretKey:    "tf-test-secret-key-updated",
	},
}

func TestAccRepositoryLocalAccountResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccRepositoryLocalAccountConfig_MultipleSecretManagers(),
				ExpectError: regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			{
				Config:      testAccRepositoryLocalAccountConfig_MultipleSecretManagersOfSameType(),
				ExpectError: regexp.MustCompile("Error: Too many .* blocks"),
			},
		},
	})
}

func testAccRepositoryLocalAccountConfig_MultipleSecretManagers() string {
	return `
	resource "cyral_repository" "tf_test_repository" {
		type = "postgresql"
		host = "http://postgres.local/"
		port = 5432
		name = "tf-test-postgres-multiple-secret-managers"
	}

	resource "cyral_repository_local_account" "tf_test_repository_account" {
		repository_id = cyral_repository.tf_test_repository.id
		aws_iam {
			database_name = "some-db-name"
			local_account = "some-local-account"
			role_arn = "some-role-arn"
		}
		kubernetes_secret {
			database_name = "some-db-name"
			local_account = "some-local-account-2"
			secret_name = "some-secret-name"
			secret_key = "some-secret-key"
		}
	}
	`
}

func testAccRepositoryLocalAccountConfig_MultipleSecretManagersOfSameType() string {
	return `
	resource "cyral_repository" "tf_test_repository" {
		type = "postgresql"
		host = "http://postgres.local/"
		port = 5432
		name = "tf-test-postgres-multiple-secret-managers"
	}

	resource "cyral_repository_local_account" "tf_test_repository_account" {
		repository_id = cyral_repository.tf_test_repository.id
		kubernetes_secret {
			database_name = "some-db-name-1"
			local_account = "some-local-account-1"
			secret_name = "some-secret-name-1"
			secret_key = "some-secret-key-1"
		}
		kubernetes_secret {
			database_name = "some-db-name-2"
			local_account = "some-local-account-2"
			secret_name = "some-secret-name-2"
			secret_key = "some-secret-key-2"
		}
	}
	`
}

func TestAccRepositoryLocalAccountResource_KubernetesSecret(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigKubernetesSecret)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigKubernetesSecret)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialCheck,
			},
			{
				Config: testUpdatedConfig,
				Check:  testUpdatedCheck,
			},
		},
	})
}

func TestAccRepositoryLocalAccountResource_EnvironmentVariable(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigEnvironmentVariable)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigEnvironmentVariable)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialCheck,
			},
			{
				Config: testUpdatedConfig,
				Check:  testUpdatedCheck,
			},
			{ // Deprecated: Should be removed in the next MAJOR release
				Config:      testAccRepositoryLocalAccountConfig_UseEnvVarTogetherWithDeprecated(),
				ExpectError: regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			{ // Deprecated: Should be removed in the next MAJOR release
				Config: testAccRepositoryLocalAccountConfig_DeprecatedEnvironmentVariable(),
				Check:  testAccRepositoryLocalAccountCheck_DeprecatedEnvironmentVariable(),
			},
		},
	})
}

func testAccRepositoryLocalAccountConfig_UseEnvVarTogetherWithDeprecated() string {
	return `
	resource "cyral_repository" "tf_test_repository" {
		type = "postgresql"
		host = "http://postgres.local/"
		port = 5432
		name = "tf-test-postgres-multiple-secret-managers"
	}

	resource "cyral_repository_local_account" "tf_test_repository_account" {
		repository_id = cyral_repository.tf_test_repository.id
		environment_variable {
			database_name = "some-db-name-1"
			local_account = "some-local-account"
			variable_name = "CYRAL_DBSECRETS_SOME_VARIABLE_NAME"
		}
		enviroment_variable {
			database_name = "some-db-name-2"
			local_account = "some-local-account"
			variable_name = "CYRAL_DBSECRETS_SOME_VARIABLE_NAME"
		}
	}
	`
}

func testAccRepositoryLocalAccountConfig_DeprecatedEnvironmentVariable() string {
	return `
	resource "cyral_repository" "tf_test_repository" {
		type = "postgresql"
		host = "http://postgres.local/"
		port = 5432
		name = "tf-test-postgres-multiple-secret-managers"
	}

	resource "cyral_repository_local_account" "tf_test_repository_account" {
		repository_id = cyral_repository.tf_test_repository.id
		enviroment_variable {
			database_name = "some-db-name"
			local_account = "some-local-account"
			variable_name = "CYRAL_DBSECRETS_SOME_VARIABLE_NAME"
		}
	}
	`
}

func testAccRepositoryLocalAccountCheck_DeprecatedEnvironmentVariable() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		// Ensure that non-deprecated argument is not set.
		resource.TestCheckNoResourceAttr(
			"cyral_repository_local_account.tf_test_repository_account",
			"environment_variable.0.database_name",
		),
		resource.TestCheckNoResourceAttr(
			"cyral_repository_local_account.tf_test_repository_account",
			"environment_variable.0.local_account",
		),
		resource.TestCheckNoResourceAttr(
			"cyral_repository_local_account.tf_test_repository_account",
			"environment_variable.0.variable_name",
		),
		// Check values for deprecated argument.
		resource.TestCheckResourceAttr(
			"cyral_repository_local_account.tf_test_repository_account",
			"enviroment_variable.0.database_name", "some-db-name",
		),
		resource.TestCheckResourceAttr(
			"cyral_repository_local_account.tf_test_repository_account",
			"enviroment_variable.0.local_account", "some-local-account",
		),
		resource.TestCheckResourceAttr(
			"cyral_repository_local_account.tf_test_repository_account",
			"enviroment_variable.0.variable_name", "CYRAL_DBSECRETS_SOME_VARIABLE_NAME",
		),
	)
}

func TestAccRepositoryLocalAccountResource_HashicorpVault(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigHashicorpVault)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigHashicorpVault)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialCheck,
			},
			{
				Config: testUpdatedConfig,
				Check:  testUpdatedCheck,
			},
		},
	})
}

func TestAccRepositoryLocalAccountResource_CyralStorage(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigCyralStorage)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigCyralStorage)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialCheck,
			},
			{
				Config: testUpdatedConfig,
				Check:  testUpdatedCheck,
			},
		},
	})
}

func TestAccRepositoryLocalAccountResource_AwsSecretsManager(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigAwsSecretsManager)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigAwsSecretsManager)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialCheck,
			},
			{
				Config: testUpdatedConfig,
				Check:  testUpdatedCheck,
			},
		},
	})
}

func TestAccRepositoryLocalAccountResource_AwsIAM(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigAwsIAM)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigAwsIAM)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialCheck,
			},
			{
				Config: testUpdatedConfig,
				Check:  testUpdatedCheck,
			},
		},
	})
}

func setupRepositoryLocalAccountTest(
	data RepositoryLocalAccountResource,
) (string, resource.TestCheckFunc) {
	testConfig := formatRepositoryLocalAccountIntoConfig(data)
	testCheck := getTestCheckForRepositoryLocalAccountResource(data)

	return testConfig, testCheck
}

func formatRepositoryLocalAccountIntoConfig(data RepositoryLocalAccountResource) string {
	const RepositoryAccountTemplate = `
	resource "cyral_repository" "tf_test_repository" {
		type = "mysql"
		host = "http://mysql.local/"
		port = 3306
		name = "%s"
	}

	resource "cyral_repository_local_account" "tf_test_repository_account" {
		repository_id = cyral_repository.tf_test_repository.id
		%s
	}
	`

	name := ""
	config := ""

	if data.AwsIAM != nil {
		name = "tf-test-mysql-aws-iam"
		config = fmt.Sprintf(`aws_iam {
			database_name = "%s"
			local_account = "%s"
			role_arn      = "%s"
		  }`,
			data.AwsIAM.DatabaseName,
			data.AwsIAM.RepoAccount,
			data.AwsIAM.RoleArn,
		)
	} else if data.AwsSecretsManager != nil {
		name = "tf-test-mysql-aws-secrets-manager"
		config = fmt.Sprintf(`aws_secrets_manager {
			database_name = "%s"
			local_account = "%s"
			secret_arn    = "%s"
		  }`,
			data.AwsSecretsManager.DatabaseName,
			data.AwsSecretsManager.RepoAccount,
			data.AwsSecretsManager.SecretArn,
		)
	} else if data.CyralStorage != nil {
		name = "tf-test-mysql-cyral-storage"
		config = fmt.Sprintf(`cyral_storage {
			database_name = "%s"
			local_account = "%s"
			password      = "%s"
			}`,
			data.CyralStorage.DatabaseName,
			data.CyralStorage.RepoAccount,
			data.CyralStorage.Password,
		)
	} else if data.HashicorpVault != nil {
		name = "tf-test-mysql-hashicorp-vault"
		config = fmt.Sprintf(`hashicorp_vault {
			database_name = "%s"
			local_account = "%s"
			path          = "%s"
			}`,
			data.HashicorpVault.DatabaseName,
			data.HashicorpVault.RepoAccount,
			data.HashicorpVault.Path,
		)
	} else if data.EnvironmentVariable != nil {
		name = "tf-test-mysql-environment-variable"
		config = fmt.Sprintf(`environment_variable {
			database_name = "%s"
			local_account = "%s"
			variable_name = "%s"
		  }`,
			data.EnvironmentVariable.DatabaseName,
			data.EnvironmentVariable.RepoAccount,
			data.EnvironmentVariable.VariableName,
		)
	} else if data.KubernetesSecret != nil {
		name = "tf-test-mysql-kubernetes-secret"
		config = fmt.Sprintf(`kubernetes_secret {
			database_name = "%s"
			local_account = "%s"
			secret_name = "%s"
			secret_key = "%s"
		  }`,
			data.KubernetesSecret.DatabaseName,
			data.KubernetesSecret.RepoAccount,
			data.KubernetesSecret.SecretName,
			data.KubernetesSecret.SecretKey,
		)
	}

	return fmt.Sprintf(RepositoryAccountTemplate, name, config)
}

func getTestCheckForRepositoryLocalAccountResource(
	data RepositoryLocalAccountResource,
) resource.TestCheckFunc {
	var testCheckFunc resource.TestCheckFunc

	if data.AwsIAM != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"aws_iam.0.database_name", data.AwsIAM.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"aws_iam.0.local_account", data.AwsIAM.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"aws_iam.0.role_arn", data.AwsIAM.RoleArn,
			),
		)
	} else if data.AwsSecretsManager != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"aws_secrets_manager.0.database_name", data.AwsSecretsManager.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"aws_secrets_manager.0.local_account", data.AwsSecretsManager.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"aws_secrets_manager.0.secret_arn", data.AwsSecretsManager.SecretArn,
			),
		)
	} else if data.CyralStorage != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"cyral_storage.0.database_name", data.CyralStorage.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"cyral_storage.0.local_account", data.CyralStorage.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"cyral_storage.0.password", data.CyralStorage.Password,
			),
		)
	} else if data.HashicorpVault != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"hashicorp_vault.0.database_name", data.HashicorpVault.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"hashicorp_vault.0.local_account", data.HashicorpVault.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"hashicorp_vault.0.path", data.HashicorpVault.Path,
			),
		)
	} else if data.EnvironmentVariable != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"environment_variable.0.database_name", data.EnvironmentVariable.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"environment_variable.0.local_account", data.EnvironmentVariable.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"environment_variable.0.variable_name", data.EnvironmentVariable.VariableName,
			),
		)
	} else if data.KubernetesSecret != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"kubernetes_secret.0.database_name", data.KubernetesSecret.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"kubernetes_secret.0.local_account", data.KubernetesSecret.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"kubernetes_secret.0.secret_name", data.KubernetesSecret.SecretName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.tf_test_repository_account",
				"kubernetes_secret.0.secret_key", data.KubernetesSecret.SecretKey,
			),
		)
	}

	return testCheckFunc
}
