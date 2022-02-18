package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialRepoAccountConfigAwsIam RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	AwsIAM: &AwsIAMResource{
		DatabaseName: "tf_test_db_name",
		RepoAccount:  "tf_test_repo_account",
		RoleArn:      "tf_test_role_arn",
	},
}

var updateRepoAccountConfigAwsIam RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	AwsIAM: &AwsIAMResource{
		DatabaseName: "tf_test_update_db_name",
		RepoAccount:  "tf_test_update_repo_account",
		RoleArn:      "tf_test_update_role_arn",
	},
}

var initialRepoAccountConfigAwsSecret RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	AwsSecretsManager: &AwsSecretsResource{
		DatabaseName: "tf_test_db_name",
		RepoAccount:  "tf_test_repo_account",
		SecretArn:    "tf_test_secret_arn",
	},
}

var updateRepoAccountConfigAwsSecret RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	AwsSecretsManager: &AwsSecretsResource{
		DatabaseName: "tf_test_update_db_name",
		RepoAccount:  "tf_test_update_repo_account",
		SecretArn:    "tf_test_update_secret_arn",
	},
}

var initialRepoAccountConfigCyralStorage RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	CyralStorage: &CyralStorageResource{
		DatabaseName: "tf_test_db_name",
		RepoAccount:  "tf_test_repo_account",
		Password:     "tf_test_pasword",
	},
}

var updateRepoAccountConfigCyralStorage RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	CyralStorage: &CyralStorageResource{
		DatabaseName: "tf_test_update_db_name",
		RepoAccount:  "tf_test_update_repo_account",
		Password:     "tf_test_update_pasword",
	},
}

var initialRepoAccountConfigHashicorpVault RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	HashicorpVault: &HashicorpVaultResource{
		DatabaseName: "tf_test_db_name",
		RepoAccount:  "tf_test_repo_account",
		Path:         "tf_test_path",
	},
}

var updateRepoAccountConfigHashicorpVault RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	HashicorpVault: &HashicorpVaultResource{
		DatabaseName: "tf_test_update_db_name",
		RepoAccount:  "tf_test_update_repo_account",
		Path:         "tf_test_update_path",
	},
}

var initialRepoAccountConfigEnviromentVariable RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	EnviromentVariable: &EnviromentVariableResource{
		DatabaseName: "tf_test_db_name",
		RepoAccount:  "tf_test_repo_account",
		VariableName: "CYRAL_DBSECRETS_TF_TEST_VARIABLE_NAME",
	},
}

var updateRepoAccountConfigEnviromentVariable RepositoryLocalAccountResource = RepositoryLocalAccountResource{
	EnviromentVariable: &EnviromentVariableResource{
		DatabaseName: "tf_test_update_db_name",
		RepoAccount:  "tf_test_update_repo_account",
		VariableName: "CYRAL_DBSECRETS_TF_TEST_UPDATE_VARIABLE_NAME",
	},
}

func TestAccRepositoryAccountEnviromentVariable(t *testing.T) {
	testConfigEnviromentVariable, testFuncEnviromentVariable := setupRepositoryAccountTest(initialRepoAccountConfigEnviromentVariable)
	testUpdateConfigEnviromentVariable, testUpdateFuncEnviromentVariable := setupRepositoryAccountTest(updateRepoAccountConfigEnviromentVariable)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigEnviromentVariable,
				Check:  testFuncEnviromentVariable,
			},
			{
				Config: testUpdateConfigEnviromentVariable,
				Check:  testUpdateFuncEnviromentVariable,
			},
		},
	})
}

func TestAccRepositoryAccountAwsHashicorpVault(t *testing.T) {
	testConfigHashicorpVault, testFuncHashicorpVault := setupRepositoryAccountTest(initialRepoAccountConfigHashicorpVault)
	testUpdateConfigHashicorpVault, testUpdateFuncHashicorpVault := setupRepositoryAccountTest(updateRepoAccountConfigHashicorpVault)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHashicorpVault,
				Check:  testFuncHashicorpVault,
			},
			{
				Config: testUpdateConfigHashicorpVault,
				Check:  testUpdateFuncHashicorpVault,
			},
		},
	})
}

func TestAccRepositoryAccountCyralStorage(t *testing.T) {
	testConfigCyralStorage, testFuncCyralStorage := setupRepositoryAccountTest(initialRepoAccountConfigCyralStorage)
	testUpdateConfigCyralStorage, testUpdateFuncCyralStorage := setupRepositoryAccountTest(updateRepoAccountConfigCyralStorage)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCyralStorage,
				Check:  testFuncCyralStorage,
			},
			{
				Config: testUpdateConfigCyralStorage,
				Check:  testUpdateFuncCyralStorage,
			},
		},
	})
}

func TestAccRepositoryAccountAwsSecretResource(t *testing.T) {
	testConfigAwsSecret, testFuncAwsSecret := setupRepositoryAccountTest(initialRepoAccountConfigAwsSecret)
	testUpdateConfigAwsSecret, testUpdateFuncAwsSecret := setupRepositoryAccountTest(updateRepoAccountConfigAwsSecret)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAwsSecret,
				Check:  testFuncAwsSecret,
			},
			{
				Config: testUpdateConfigAwsSecret,
				Check:  testUpdateFuncAwsSecret,
			},
		},
	})
}

func TestAccRepositoryAccountAwsIamResource(t *testing.T) {
	testConfigAwsIam, testFuncAwsIam := setupRepositoryAccountTest(initialRepoAccountConfigAwsIam)
	testUpdateConfigAwsIam, testUpdateFuncAwsIam := setupRepositoryAccountTest(updateRepoAccountConfigAwsIam)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAwsIam,
				Check:  testFuncAwsIam,
			},
			{
				Config: testUpdateConfigAwsIam,
				Check:  testUpdateFuncAwsIam,
			},
		},
	})
}

func setupRepositoryAccountTest(integrationData RepositoryLocalAccountResource) (string, resource.TestCheckFunc) {
	configuration := formatRepoAccountIntoConfig(integrationData)

	testFunction := getTestFunctionForRepositoryLocalAccountResource(integrationData)

	return configuration, testFunction
}

func formatRepoAccountIntoConfig(data RepositoryLocalAccountResource) string {
	return fmt.Sprintf(`
	  %s
	  `, formatRepoAccountAuthConfig(data))
}

func formatRepoAccountAuthConfig(data RepositoryLocalAccountResource) string {
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
	}`

	name := ""
	config := ""

	if data.AwsIAM != nil {
		name = "tf-test-mysql-aws_iam"
		config = fmt.Sprintf(`aws_iam {
			database_name = "%s"
			local_account = "%s"
			role_arn      = "%s"
		  }`, data.AwsIAM.DatabaseName, data.AwsIAM.RepoAccount, data.AwsIAM.RoleArn)
	} else if data.AwsSecretsManager != nil {
		name = "tf-test-mysql-aws_secrets_manager"

		config = fmt.Sprintf(`aws_secrets_manager {
			database_name = "%s"
			local_account = "%s"
			secret_arn    = "%s"
		  }`, data.AwsSecretsManager.DatabaseName, data.AwsSecretsManager.RepoAccount, data.AwsSecretsManager.SecretArn)
	} else if data.CyralStorage != nil {
		name = "tf-test-mysql-cyral_storage"

		config = fmt.Sprintf(`cyral_storage {
			database_name = "%s"
			local_account = "%s"
			password      = "%s"
		  }`, data.CyralStorage.DatabaseName, data.CyralStorage.RepoAccount, data.CyralStorage.Password)
	} else if data.HashicorpVault != nil {
		name = "tf-test-mysql-hashicorp_vault"

		config = fmt.Sprintf(`hashicorp_vault {
			database_name = "%s"
			local_account = "%s"
			path          = "%s"
		  }`, data.HashicorpVault.DatabaseName, data.HashicorpVault.RepoAccount, data.HashicorpVault.Path)
	} else if data.EnviromentVariable != nil {
		name = "tf-test-mysql-enviroment_variable"

		config = fmt.Sprintf(`enviroment_variable {
			database_name = "%s"
			local_account = "%s"
			variable_name = "%s"
		  }`, data.EnviromentVariable.DatabaseName, data.EnviromentVariable.RepoAccount, data.EnviromentVariable.VariableName)
	}

	return fmt.Sprintf(RepositoryAccountTemplate, name, config)
}

func getTestFunctionForRepositoryLocalAccountResource(data RepositoryLocalAccountResource) resource.TestCheckFunc {
	var testFunc resource.TestCheckFunc

	if data.AwsIAM != nil {
		testFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "aws_iam.0.database_name", data.AwsIAM.DatabaseName),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "aws_iam.0.local_account", data.AwsIAM.RepoAccount),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "aws_iam.0.role_arn", data.AwsIAM.RoleArn),
		)
	} else if data.AwsSecretsManager != nil {
		testFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "aws_secrets_manager.0.database_name", data.AwsSecretsManager.DatabaseName),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "aws_secrets_manager.0.local_account", data.AwsSecretsManager.RepoAccount),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "aws_secrets_manager.0.secret_arn", data.AwsSecretsManager.SecretArn),
		)
	} else if data.CyralStorage != nil {
		testFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "cyral_storage.0.database_name", data.CyralStorage.DatabaseName),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "cyral_storage.0.local_account", data.CyralStorage.RepoAccount),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "cyral_storage.0.password", data.CyralStorage.Password),
		)
	} else if data.HashicorpVault != nil {
		testFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "hashicorp_vault.0.database_name", data.HashicorpVault.DatabaseName),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "hashicorp_vault.0.local_account", data.HashicorpVault.RepoAccount),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "hashicorp_vault.0.path", data.HashicorpVault.Path),
		)
	} else if data.EnviromentVariable != nil {
		testFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "enviroment_variable.0.database_name", data.EnviromentVariable.DatabaseName),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "enviroment_variable.0.local_account", data.EnviromentVariable.RepoAccount),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "enviroment_variable.0.variable_name", data.EnviromentVariable.VariableName),
		)
	}

	return testFunc
}
