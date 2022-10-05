package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testRepoLocalAccImportName = "cyral_repository_local_account.test_repository_account"

	repositoryLocalAccountResourceName = "repository-local-account"
)

var initialRepoAccountConfigAwsIAM = RepositoryLocalAccountResource{
	AwsIAM: &AwsIAMResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account"),
		RoleArn:      "tf-test-role-arn",
	},
}

var updatedRepoAccountConfigAwsIAM = RepositoryLocalAccountResource{
	AwsIAM: &AwsIAMResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name-updated"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account-updated"),
		RoleArn:      "tf-test-role-arn-updated",
	},
}

var initialRepoAccountConfigAwsSecretsManager = RepositoryLocalAccountResource{
	AwsSecretsManager: &AwsSecretsManagerResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account"),
		SecretArn:    "tf-test-secret-arn",
	},
}

var updatedRepoAccountConfigAwsSecretsManager = RepositoryLocalAccountResource{
	AwsSecretsManager: &AwsSecretsManagerResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name-updated"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account-updated"),
		SecretArn:    "tf-test-secret-arn-updated",
	},
}

var initialRepoAccountConfigCyralStorage = RepositoryLocalAccountResource{
	CyralStorage: &CyralStorageResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account"),
		Password:     "tf-test-pasword",
	},
}

var updatedRepoAccountConfigCyralStorage = RepositoryLocalAccountResource{
	CyralStorage: &CyralStorageResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name-updated"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account-updated"),
		Password:     "tf-test-pasword-updated",
	},
}

var initialRepoAccountConfigHashicorpVault = RepositoryLocalAccountResource{
	HashicorpVault: &HashicorpVaultResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account"),
		Path:         "tf-test-path",
	},
}

var updatedRepoAccountConfigHashicorpVault = RepositoryLocalAccountResource{
	HashicorpVault: &HashicorpVaultResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name-updated"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account-updated"),
		Path:         "tf-test-path-updated",
	},
}

var initialRepoAccountConfigEnvironmentVariable = RepositoryLocalAccountResource{
	EnvironmentVariable: &EnvironmentVariableResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account"),
		VariableName: "CYRAL_DBSECRETS_TF_TEST_VARIABLE_NAME",
	},
}

var updatedRepoAccountConfigEnvironmentVariable = RepositoryLocalAccountResource{
	EnvironmentVariable: &EnvironmentVariableResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name-updated"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account-updated"),
		VariableName: "CYRAL_DBSECRETS_TF_TEST_VARIABLE_NAME_UPDATED",
	},
}

var initialRepoAccountConfigKubernetesSecret = RepositoryLocalAccountResource{
	KubernetesSecret: &KubernetesSecretResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account"),
		SecretName:   "tf-test-db-secrets",
		SecretKey:    "tf-test-secret-key",
	},
}

var updatedRepoAccountConfigKubernetesSecret = RepositoryLocalAccountResource{
	KubernetesSecret: &KubernetesSecretResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name-updated"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account-updated"),
		SecretName:   "db-secrets-updated",
		SecretKey:    "tf-test-secret-key-updated",
	},
}

var initialRepoAccountConfigGcpSecretManager = RepositoryLocalAccountResource{
	GcpSecretManager: &GcpSecretManagerResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account"),
		SecretName:   "projects/1234567890/secrets/my-secret/versions/1",
	},
}

var updatedRepoAccountConfigGcpSecretManager = RepositoryLocalAccountResource{
	GcpSecretManager: &GcpSecretManagerResource{
		DatabaseName: accTestName(repositoryLocalAccountResourceName, "db-name-updated"),
		RepoAccount:  accTestName(repositoryLocalAccountResourceName, "repo-account-updated"),
		SecretName:   "projects/1234567890/secrets/my-secret-updated/versions/2",
	},
}

func repositoryLocalAccountSampleRepositoryConfig() string {
	return formatBasicRepositoryIntoConfig(
		basicRepositoryResName,
		accTestName(repositoryLocalAccountResourceName, "repository"),
		"postgresql",
		"http://postgres.local/",
		5432,
	)
}

func testRespositoryLocalAccountImportState(resName string) resource.TestStep {
	return resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: importStateComposedIDFunc(
			resName,
			[]string{"repository_id", "id"},
			"/"),
		ResourceName: resName,
	}
}

func TestAccRepositoryLocalAccountResource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
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
			testAccRepositoryLocalAccount_AutoApprovalNoMax(),
			testAccRepositoryLocalAccount_AutoApprovalWithMax(),
		},
	})
}

func testAccRepositoryLocalAccountConfig_MultipleSecretManagers() string {
	var config string
	config += repositoryLocalAccountSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_local_account" "test_repository_account" {
		repository_id = %s
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
	}`, basicRepositoryID)
	return config
}

func testAccRepositoryLocalAccountConfig_MultipleSecretManagersOfSameType() string {
	var config string
	config += repositoryLocalAccountSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_local_account" "test_repository_account" {
		repository_id = %s
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
	}`, basicRepositoryID)
	return config
}

func testAccRepositoryLocalAccount_AutoApprovalNoMax() resource.TestStep {
	var config string
	config += repositoryLocalAccountSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_local_account" "test_repository_account" {
		repository_id = %s
		config {
			auto_approve_access = true
		}
		aws_iam {
			database_name = "some-db-name"
			local_account = "some-local-account"
			role_arn = "some-role-arn"
		}
	}`, basicRepositoryID)

	fullResName := "cyral_repository_local_account.test_repository_account"
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(fullResName,
			"config.0.auto_approve_access", "true",
		),
		resource.TestCheckResourceAttr(fullResName,
			"config.0.max_auto_approve_duration", "P0D",
		),
	}

	return resource.TestStep{
		Config: config,
		Check:  resource.ComposeTestCheckFunc(checkFuncs...),
	}
}

func testAccRepositoryLocalAccount_AutoApprovalWithMax() resource.TestStep {
	var config string
	config += repositoryLocalAccountSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_local_account" "test_repository_account" {
		repository_id = %s
		config {
			auto_approve_access = true
			max_auto_approve_duration = "PT4S"
		}
		aws_iam {
			database_name = "some-db-name"
			local_account = "some-local-account"
			role_arn = "some-role-arn"
		}
	}`, basicRepositoryID)

	fullResName := "cyral_repository_local_account.test_repository_account"
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(fullResName,
			"config.0.auto_approve_access", "true",
		),
		resource.TestCheckResourceAttr(fullResName,
			"config.0.max_auto_approve_duration", "PT4S",
		),
	}

	return resource.TestStep{
		Config: config,
		Check:  resource.ComposeTestCheckFunc(checkFuncs...),
	}
}

func TestAccRepositoryLocalAccountResource_GcpSecretManager(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigGcpSecretManager)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigGcpSecretManager)

	resource.ParallelTest(t, resource.TestCase{
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
			testRespositoryLocalAccountImportState(testRepoLocalAccImportName),
		},
	})
}

func TestAccRepositoryLocalAccountResource_KubernetesSecret(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigKubernetesSecret)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigKubernetesSecret)

	resource.ParallelTest(t, resource.TestCase{
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
			testRespositoryLocalAccountImportState(testRepoLocalAccImportName),
		},
	})
}

func TestAccRepositoryLocalAccountResource_EnvironmentVariable(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigEnvironmentVariable)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigEnvironmentVariable)

	resource.ParallelTest(t, resource.TestCase{
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
			testRespositoryLocalAccountImportState(testRepoLocalAccImportName),
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
	var config string
	config += repositoryLocalAccountSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_local_account" "test_repository_account" {
		repository_id = %s
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
	}`, basicRepositoryID)
	return config
}

func testAccRepositoryLocalAccountConfig_DeprecatedEnvironmentVariable() string {
	var config string
	config += repositoryLocalAccountSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_local_account" "test_repository_account" {
		repository_id = %s
		enviroment_variable {
			database_name = "some-db-name"
			local_account = "some-local-account"
			variable_name = "CYRAL_DBSECRETS_SOME_VARIABLE_NAME"
		}
	}`, basicRepositoryID)
	return config
}

func testAccRepositoryLocalAccountCheck_DeprecatedEnvironmentVariable() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		// Ensure that non-deprecated argument is not set.
		resource.TestCheckNoResourceAttr(
			"cyral_repository_local_account.test_repository_account",
			"environment_variable.0.database_name",
		),
		resource.TestCheckNoResourceAttr(
			"cyral_repository_local_account.test_repository_account",
			"environment_variable.0.local_account",
		),
		resource.TestCheckNoResourceAttr(
			"cyral_repository_local_account.test_repository_account",
			"environment_variable.0.variable_name",
		),
		// Check values for deprecated argument.
		resource.TestCheckResourceAttr(
			"cyral_repository_local_account.test_repository_account",
			"enviroment_variable.0.database_name", "some-db-name",
		),
		resource.TestCheckResourceAttr(
			"cyral_repository_local_account.test_repository_account",
			"enviroment_variable.0.local_account", "some-local-account",
		),
		resource.TestCheckResourceAttr(
			"cyral_repository_local_account.test_repository_account",
			"enviroment_variable.0.variable_name", "CYRAL_DBSECRETS_SOME_VARIABLE_NAME",
		),
	)
}

func TestAccRepositoryLocalAccountResource_HashicorpVault(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigHashicorpVault)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigHashicorpVault)

	resource.ParallelTest(t, resource.TestCase{
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
			testRespositoryLocalAccountImportState(testRepoLocalAccImportName),
		},
	})
}

func TestAccRepositoryLocalAccountResource_CyralStorage(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigCyralStorage)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigCyralStorage)

	resource.ParallelTest(t, resource.TestCase{
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
			testRespositoryLocalAccountImportState(testRepoLocalAccImportName),
		},
	})
}

func TestAccRepositoryLocalAccountResource_AwsSecretsManager(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigAwsSecretsManager)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigAwsSecretsManager)

	resource.ParallelTest(t, resource.TestCase{
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
			testRespositoryLocalAccountImportState(testRepoLocalAccImportName),
		},
	})
}

func TestAccRepositoryLocalAccountResource_AwsIAM(t *testing.T) {
	testInitialConfig, testInitialCheck :=
		setupRepositoryLocalAccountTest(initialRepoAccountConfigAwsIAM)
	testUpdatedConfig, testUpdatedCheck :=
		setupRepositoryLocalAccountTest(updatedRepoAccountConfigAwsIAM)

	resource.ParallelTest(t, resource.TestCase{
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
			testRespositoryLocalAccountImportState(testRepoLocalAccImportName),
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
	repositoryName := ""
	localAccountConfig := ""

	if data.AwsIAM != nil {
		repositoryName = accTestName(repositoryLocalAccountResourceName,
			"mysql-aws-iam")
		localAccountConfig = fmt.Sprintf(`aws_iam {
			database_name = "%s"
			local_account = "%s"
			role_arn      = "%s"
		  }`,
			data.AwsIAM.DatabaseName,
			data.AwsIAM.RepoAccount,
			data.AwsIAM.RoleArn,
		)
	} else if data.AwsSecretsManager != nil {
		repositoryName = accTestName(repositoryLocalAccountResourceName,
			"mysql-aws-secrets-manager")
		localAccountConfig = fmt.Sprintf(`aws_secrets_manager {
			database_name = "%s"
			local_account = "%s"
			secret_arn    = "%s"
		  }`,
			data.AwsSecretsManager.DatabaseName,
			data.AwsSecretsManager.RepoAccount,
			data.AwsSecretsManager.SecretArn,
		)
	} else if data.CyralStorage != nil {
		repositoryName = accTestName(repositoryLocalAccountResourceName,
			"mysql-cyral-storage")
		localAccountConfig = fmt.Sprintf(`cyral_storage {
			database_name = "%s"
			local_account = "%s"
			password      = "%s"
			}`,
			data.CyralStorage.DatabaseName,
			data.CyralStorage.RepoAccount,
			data.CyralStorage.Password,
		)
	} else if data.HashicorpVault != nil {
		repositoryName = accTestName(repositoryLocalAccountResourceName,
			"mysql-hashicorp-vault")
		localAccountConfig = fmt.Sprintf(`hashicorp_vault {
			database_name = "%s"
			local_account = "%s"
			path          = "%s"
			}`,
			data.HashicorpVault.DatabaseName,
			data.HashicorpVault.RepoAccount,
			data.HashicorpVault.Path,
		)
	} else if data.EnvironmentVariable != nil {
		repositoryName = accTestName(repositoryLocalAccountResourceName,
			"mysql-environment-variable")
		localAccountConfig = fmt.Sprintf(`environment_variable {
			database_name = "%s"
			local_account = "%s"
			variable_name = "%s"
		  }`,
			data.EnvironmentVariable.DatabaseName,
			data.EnvironmentVariable.RepoAccount,
			data.EnvironmentVariable.VariableName,
		)
	} else if data.KubernetesSecret != nil {
		repositoryName = accTestName(repositoryLocalAccountResourceName,
			"mysql-kubernetes-secret")
		localAccountConfig = fmt.Sprintf(`kubernetes_secret {
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
	} else if data.GcpSecretManager != nil {
		repositoryName = accTestName(repositoryLocalAccountResourceName,
			"mysql-gcp-secret-manager")
		localAccountConfig = fmt.Sprintf(`gcp_secret_manager {
			database_name = "%s"
			local_account = "%s"
			secret_name = "%s"
		  }`,
			data.GcpSecretManager.DatabaseName,
			data.GcpSecretManager.RepoAccount,
			data.GcpSecretManager.SecretName,
		)
	}

	var config string
	config += formatBasicRepositoryIntoConfig(
		basicRepositoryResName,
		repositoryName,
		"mysql",
		"http://mysql.local/",
		3306,
	)
	config += fmt.Sprintf(`
	resource "cyral_repository_local_account" "test_repository_account" {
		repository_id = cyral_repository.test_repository.id
		%s
	}`, localAccountConfig)

	return config
}

func getTestCheckForRepositoryLocalAccountResource(
	data RepositoryLocalAccountResource,
) resource.TestCheckFunc {
	var testCheckFunc resource.TestCheckFunc

	if data.AwsIAM != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"aws_iam.0.database_name", data.AwsIAM.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"aws_iam.0.local_account", data.AwsIAM.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"aws_iam.0.role_arn", data.AwsIAM.RoleArn,
			),
		)
	} else if data.AwsSecretsManager != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"aws_secrets_manager.0.database_name", data.AwsSecretsManager.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"aws_secrets_manager.0.local_account", data.AwsSecretsManager.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"aws_secrets_manager.0.secret_arn", data.AwsSecretsManager.SecretArn,
			),
		)
	} else if data.CyralStorage != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"cyral_storage.0.database_name", data.CyralStorage.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"cyral_storage.0.local_account", data.CyralStorage.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"cyral_storage.0.password", data.CyralStorage.Password,
			),
		)
	} else if data.HashicorpVault != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"hashicorp_vault.0.database_name", data.HashicorpVault.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"hashicorp_vault.0.local_account", data.HashicorpVault.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"hashicorp_vault.0.path", data.HashicorpVault.Path,
			),
		)
	} else if data.EnvironmentVariable != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"environment_variable.0.database_name", data.EnvironmentVariable.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"environment_variable.0.local_account", data.EnvironmentVariable.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"environment_variable.0.variable_name", data.EnvironmentVariable.VariableName,
			),
		)
	} else if data.KubernetesSecret != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"kubernetes_secret.0.database_name", data.KubernetesSecret.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"kubernetes_secret.0.local_account", data.KubernetesSecret.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"kubernetes_secret.0.secret_name", data.KubernetesSecret.SecretName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"kubernetes_secret.0.secret_key", data.KubernetesSecret.SecretKey,
			),
		)
	} else if data.GcpSecretManager != nil {
		testCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"gcp_secret_manager.0.database_name", data.GcpSecretManager.DatabaseName,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"gcp_secret_manager.0.local_account", data.GcpSecretManager.RepoAccount,
			),
			resource.TestCheckResourceAttr(
				"cyral_repository_local_account.test_repository_account",
				"gcp_secret_manager.0.secret_name", data.GcpSecretManager.SecretName,
			),
		)
	}

	return testCheckFunc
}
