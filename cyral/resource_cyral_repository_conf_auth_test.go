package cyral

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

const (
	repositoryConfAuthResourceName = "repository-conf-auth"
)

func repositoryConfAuthDependencyConfig() string {
	return formatBasicRepositoryIntoConfig(
		basicRepositoryResName,
		accTestName(repositoryConfAuthResourceName, "repository"),
		"mysql",
		"http://mysql.local/",
		3306,
	)
}

func initialRepositoryConfAuthConfig() RepositoryConfAuthData {
	return RepositoryConfAuthData{
		AllowNativeAuth: false,
		ClientTLS:       "disable",
		RepoTLS:         "enable",
	}
}

func update1RepositoryConfAuthConfig() RepositoryConfAuthData {
	return RepositoryConfAuthData{
		AllowNativeAuth: true,
		ClientTLS:       "enable",
		RepoTLS:         "disable",
	}
}

func update2RepositoryConfAuthConfig() RepositoryConfAuthData {
	return RepositoryConfAuthData{
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
	}`, resName, basicRepositoryID)

	return resource.TestStep{
		Config: config,
		Check: setupRepositoryConfAuthCheck(
			resName,
			RepositoryConfAuthData{
				ClientTLS: defaultClientTLS,
				RepoTLS:   defaultRepoTLS,
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
	actualNewState, err := upgradeRepositoryConfAuthV0(context.Background(),
		previousState, nil)
	require.NoError(t, err)
	expectedNewState := map[string]interface{}{
		"id":            "my-repository-id",
		"repository_id": "my-repository-id",
	}
	require.Equal(t, expectedNewState, actualNewState)
}

func setupRepositoryConfAuthTest(resName string, repositoryConf RepositoryConfAuthData) resource.TestStep {
	return resource.TestStep{
		Config: setupRepositoryConfAuthConfig(resName, repositoryConf),
		Check:  setupRepositoryConfAuthCheck(resName, repositoryConf),
	}
}

func setupRepositoryConfAuthConfig(resName string, repositoryConf RepositoryConfAuthData) string {
	var config string
	config += repositoryConfAuthDependencyConfig()
	config += formatRepositoryConfAuthDataIntoConfig(
		resName, repositoryConf, basicRepositoryID)

	return config
}

func setupRepositoryConfAuthCheck(resName string, repositoryConf RepositoryConfAuthData) resource.TestCheckFunc {
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
	data RepositoryConfAuthData,
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
