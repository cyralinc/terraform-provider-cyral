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

func TestAccRepositoryAccountAwsHashicorpVault(t *testing.T) {
	testConfigHashicorpVault, testFuncHashicorpVault := setupRepositoryAccountTest(initialRepoAccountConfigHashicorpVault)
	testUpdateConfigHashicorpVault, testUpdateFuncHashicorpVault := setupRepositoryAccountTest(updateRepoAccountConfigHashicorpVault)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
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
		PreCheck:          func() { testAccPreCheck(t) },
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

func TestAccRepositoryAccountAwsIamResource(t *testing.T) {
	testConfigAwsIam, testFuncAwsIam := setupRepositoryAccountTest(initialRepoAccountConfigAwsIam)
	testUpdateConfigAwsIam, testUpdateFuncAwsIam := setupRepositoryAccountTest(updateRepoAccountConfigAwsIam)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
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

func TestAccRepositoryAccountAwsSecretResource(t *testing.T) {
	testConfigAwsSecret, testFuncAwsSecret := setupRepositoryAccountTest(initialRepoAccountConfigAwsSecret)
	testUpdateConfigAwsSecret, testUpdateFuncAwsSecret := setupRepositoryAccountTest(updateRepoAccountConfigAwsSecret)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
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

func setupRepositoryAccountTest(integrationData RepositoryLocalAccountResource) (string, resource.TestCheckFunc) {
	configuration := formatRepoAccountIntoConfig(integrationData)

	testFunction := getTestFunctionForRepositoryLocalAccountResource(integrationData)

	return configuration, testFunction
}

func formatRepoAccountIntoConfig(data RepositoryLocalAccountResource) string {
	return fmt.Sprintf(`
	resource "cyral_repository" "tf_test_repository" {
		type = "mysql"
		host = "http://mysql.local/"
		port = 3306
		name = "tf-test-mysql-2"
	  }
	  %s 
	  `, formatRepoAccountAuthConfig(data))
}

func formatRepoAccountAuthConfig(data RepositoryLocalAccountResource) string {
	const RepositoryAccountTemplate = `resource "cyral_repository_local_account" "tf_test_repository_account" {
		repository_id = cyral_repository.tf_test_repository.id
		%s
	}`

	config := ""

	if data.AwsIAM != nil {
		config = fmt.Sprintf(`aws_iam {
			database_name = "%s"
			local_account = "%s"
			role_arn      = "%s"
		  }`, data.AwsIAM.DatabaseName, data.AwsIAM.RepoAccount, data.AwsIAM.RoleArn)
	} else if data.AwsSecretsManager != nil {
		config = fmt.Sprintf(`aws_secrets_manager {
			database_name = "%s"
			local_account = "%s"
			secret_arn    = "%s"
		  }`, data.AwsSecretsManager.DatabaseName, data.AwsSecretsManager.RepoAccount, data.AwsSecretsManager.SecretArn)
	} else if data.CyralStorage != nil {
		config = fmt.Sprintf(`cyral_storage {
			database_name = "%s"
			local_account = "%s"
			password      = "%s"
		  }`, data.CyralStorage.DatabaseName, data.CyralStorage.RepoAccount, data.CyralStorage.Password)
	} else if data.HashicorpVault != nil {
		config = fmt.Sprintf(`hashicorp_vault {
			database_name = "%s"
			local_account = "%s"
			path          = "%s"
		  }`, data.HashicorpVault.DatabaseName, data.HashicorpVault.RepoAccount, data.HashicorpVault.Path)
	}

	return fmt.Sprintf(RepositoryAccountTemplate, config)
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
	} else if data.AwsIAM != nil {
		testFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "hashicorp_vault.0.database_name", data.HashicorpVault.DatabaseName),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "hashicorp_vault.0.local_account", data.HashicorpVault.RepoAccount),
			resource.TestCheckResourceAttr("cyral_repository_local_account.tf_test_repository_account", "hashicorp_vault.0.path", data.HashicorpVault.Path),
		)
	}

	return testFunc
}
