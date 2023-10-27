package internal_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

const (
	repositoryConfAuthResourceName = "repository-conf-auth"
)

func repositoryConfAuthDependencyConfig() string {
	return utils.FormatBasicRepositoryIntoConfig(
		utils.BasicRepositoryResName,
		utils.AccTestName(repositoryConfAuthResourceName, "repository"),
		"mysql",
		"http://mysql.local/",
		3306,
	)
}

func initialRepositoryConfAuthConfig() internal.RepositoryConfAuthData {
	return internal.RepositoryConfAuthData{
		AllowNativeAuth: false,
		ClientTLS:       "disable",
		RepoTLS:         "enable",
	}
}

func update1RepositoryConfAuthConfig() internal.RepositoryConfAuthData {
	return internal.RepositoryConfAuthData{
		AllowNativeAuth: true,
		ClientTLS:       "enable",
		RepoTLS:         "disable",
	}
}

func update2RepositoryConfAuthConfig() internal.RepositoryConfAuthData {
	return internal.RepositoryConfAuthData{
		AllowNativeAuth: false,
		ClientTLS:       "enable",
		RepoTLS:         "disable",
	}
}

// This tests an empty config to avoid regressions. In the past, we had a
// problem with infinite apply loops for this resource, when running for an
// empty config (using default values). See issue #286.
func repositoryConfAuthMinimalConfigTest(resName string) resource.TestStep {
	var config string
	config += repositoryConfAuthDependencyConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_conf_auth" "%s" {
		repository_id = %s
	}`, resName, utils.BasicRepositoryID)

	return resource.TestStep{
		Config: config,
		Check: setupRepositoryConfAuthCheck(
			resName,
			internal.RepositoryConfAuthData{
				ClientTLS: internal.DefaultClientTLS,
				RepoTLS:   internal.DefaultRepoTLS,
			},
		),
	}
}

func TestAccRepositoryConfAuthResource(t *testing.T) {
	testMinimal := repositoryConfAuthMinimalConfigTest("main_test")

	mainTest := setupRepositoryConfAuthTest("main_test", initialRepositoryConfAuthConfig())
	mainTestUpdate1 := setupRepositoryConfAuthTest("main_test", update1RepositoryConfAuthConfig())
	mainTestUpdate2 := setupRepositoryConfAuthTest("main_test", update2RepositoryConfAuthConfig())

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			testMinimal,

			mainTest,
			mainTestUpdate1,
			mainTestUpdate2,

			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_repository_conf_auth.main_test",
			},
		},
	})
}

func TestRepositoryConfAuthResourceUpgradeV0(t *testing.T) {
	previousState := map[string]interface{}{
		"id":            "repo-conf",
		"repository_id": "my-repository-id",
	}
	actualNewState, err := internal.UpgradeRepositoryConfAuthV0(context.Background(),
		previousState, nil)
	require.NoError(t, err)
	expectedNewState := map[string]interface{}{
		"id":            "my-repository-id",
		"repository_id": "my-repository-id",
	}
	require.Equal(t, expectedNewState, actualNewState)
}

func setupRepositoryConfAuthTest(resName string, repositoryConf internal.RepositoryConfAuthData) resource.TestStep {
	return resource.TestStep{
		Config: setupRepositoryConfAuthConfig(resName, repositoryConf),
		Check:  setupRepositoryConfAuthCheck(resName, repositoryConf),
	}
}

func setupRepositoryConfAuthConfig(resName string, repositoryConf internal.RepositoryConfAuthData) string {
	var config string
	config += repositoryConfAuthDependencyConfig()
	config += formatRepositoryConfAuthDataIntoConfig(
		resName, repositoryConf, utils.BasicRepositoryID)

	return config
}

func setupRepositoryConfAuthCheck(resName string, repositoryConf internal.RepositoryConfAuthData) resource.TestCheckFunc {
	resourceFullName := fmt.Sprintf("cyral_repository_conf_auth.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"allow_native_auth", fmt.Sprintf("%t", repositoryConf.AllowNativeAuth)),
		resource.TestCheckResourceAttr(resourceFullName,
			"client_tls", repositoryConf.ClientTLS),
		resource.TestCheckResourceAttr(resourceFullName,
			"repo_tls", repositoryConf.RepoTLS),
	)
}

func formatRepositoryConfAuthDataIntoConfig(
	resName string,
	data internal.RepositoryConfAuthData,
	repositoryID string,
) string {
	return fmt.Sprintf(`
	resource "cyral_repository_conf_auth" "%s" {
		repository_id = %s
		allow_native_auth = %t
		client_tls = "%s"
		identity_provider = "tf_test_conf_auth_okta"
		repo_tls = "%s"
	}`, resName, repositoryID, data.AllowNativeAuth, data.ClientTLS,
		data.RepoTLS)
}
