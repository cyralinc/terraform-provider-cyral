package confauth_test

import (
	"context"
	"fmt"
	"testing"

	auth "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/confauth"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
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
		"mongodb",
		"http://mongodb.local/",
		27017,
	)
}

func initialRepositoryConfAuthConfig() auth.RepositoryConfAuthData {
	return auth.RepositoryConfAuthData{
		AllowNativeAuth: false,
		ClientTLS:       string(auth.TLSDisable),
		RepoTLS:         string(auth.TLSEnable),
		AuthType:        "ACCESS_TOKEN",
	}
}

func update1RepositoryConfAuthConfig() auth.RepositoryConfAuthData {
	return auth.RepositoryConfAuthData{
		AllowNativeAuth: true,
		ClientTLS:       string(auth.TLSEnable),
		RepoTLS:         string(auth.TLSDisable),
		AuthType:        "AWS_IAM",
	}
}

func update2RepositoryConfAuthConfig() auth.RepositoryConfAuthData {
	return auth.RepositoryConfAuthData{
		AllowNativeAuth: false,
		ClientTLS:       string(auth.TLSEnable),
		RepoTLS:         string(auth.TLSDisable),
		AuthType:        "ACCESS_TOKEN",
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
			auth.RepositoryConfAuthData{
				ClientTLS: string(auth.TLSDisable),
				RepoTLS:   string(auth.TLSDisable),
				AuthType:  auth.DefaultAuthType,
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
		ProviderFactories: provider.ProviderFactories,
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
	actualNewState, err := auth.UpgradeRepositoryConfAuthV0(context.Background(),
		previousState, nil)
	require.NoError(t, err)
	expectedNewState := map[string]interface{}{
		"id":            "my-repository-id",
		"repository_id": "my-repository-id",
	}
	require.Equal(t, expectedNewState, actualNewState)
}

func setupRepositoryConfAuthTest(resName string, repositoryConf auth.RepositoryConfAuthData) resource.TestStep {
	return resource.TestStep{
		Config: setupRepositoryConfAuthConfig(resName, repositoryConf),
		Check:  setupRepositoryConfAuthCheck(resName, repositoryConf),
	}
}

func setupRepositoryConfAuthConfig(resName string, repositoryConf auth.RepositoryConfAuthData) string {
	var config string
	config += repositoryConfAuthDependencyConfig()
	config += formatRepositoryConfAuthDataIntoConfig(
		resName, repositoryConf, utils.BasicRepositoryID)

	return config
}

func setupRepositoryConfAuthCheck(resName string, repositoryConf auth.RepositoryConfAuthData) resource.TestCheckFunc {
	resourceFullName := fmt.Sprintf("cyral_repository_conf_auth.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"allow_native_auth", fmt.Sprintf("%t", repositoryConf.AllowNativeAuth),
		),
		resource.TestCheckResourceAttr(resourceFullName,
			"client_tls", repositoryConf.ClientTLS),
		resource.TestCheckResourceAttr(resourceFullName,
			"repo_tls", repositoryConf.RepoTLS),
		resource.TestCheckResourceAttr(resourceFullName,
			"auth_type", repositoryConf.AuthType,
		),
	)
}

func formatRepositoryConfAuthDataIntoConfig(
	resName string,
	data auth.RepositoryConfAuthData,
	repositoryID string,
) string {
	return fmt.Sprintf(`
	resource "cyral_repository_conf_auth" "%s" {
		repository_id = %s
		allow_native_auth = %t
		client_tls = "%s"
		identity_provider = "tf_test_conf_auth_okta"
		repo_tls = "%s"
		auth_type = "%s"
	}`, resName, repositoryID, data.AllowNativeAuth, data.ClientTLS,
		data.RepoTLS, data.AuthType)
}
