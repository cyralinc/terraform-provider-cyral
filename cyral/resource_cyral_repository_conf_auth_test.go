package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialRepositoryConfAuthConfig RepositoryConfAuthData = RepositoryConfAuthData{
	AllowNativeAuth: false,
	ClientTLS:       "disable",
	RepoTLS:         "enable",
}

var update1RepositoryConfAuthConfig RepositoryConfAuthData = RepositoryConfAuthData{
	AllowNativeAuth: true,
	ClientTLS:       "enable",
	RepoTLS:         "disable",
}

var update2RepositoryConfAuthConfig RepositoryConfAuthData = RepositoryConfAuthData{
	AllowNativeAuth: false,
	ClientTLS:       "enable",
	RepoTLS:         "disable",
}

func TestAccRepositoryConfAuthResource(t *testing.T) {
	testConfig, testFunc := setupRepositoryConfAuthTest(initialRepositoryConfAuthConfig)
	testUpdate1Config, testUpdate1Func := setupRepositoryConfAuthTest(update1RepositoryConfAuthConfig)
	testUpdate2Config, testUpdate2Func := setupRepositoryConfAuthTest(update2RepositoryConfAuthConfig)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdate1Config,
				Check:  testUpdate1Func,
			},
			{
				Config: testUpdate2Config,
				Check:  testUpdate2Func,
			},
			// TODO: add import test -aholmquist 2022-08-05
		},
	})
}

func setupRepositoryConfAuthTest(repositoryConf RepositoryConfAuthData) (string, resource.TestCheckFunc) {
	var configuration string
	configuration += formatBasicRepositoryIntoConfig(
		basicRepositoryResName,
		accTestName("repository-conf-auth", "repository"),
		"mysql",
		"http://mysql.local/",
		3306,
	)
	configuration += formatRepositoryConfAuthDataIntoConfig(
		repositoryConf, basicRepositoryID)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_repository_conf_auth.my-repository-conf-auth",
			"allow_native_auth", fmt.Sprintf("%t", repositoryConf.AllowNativeAuth)),
		resource.TestCheckResourceAttr("cyral_repository_conf_auth.my-repository-conf-auth",
			"client_tls", repositoryConf.ClientTLS),
		resource.TestCheckResourceAttr("cyral_repository_conf_auth.my-repository-conf-auth",
			"repo_tls", repositoryConf.RepoTLS),
	)

	return configuration, testFunction
}

func formatRepositoryConfAuthDataIntoConfig(
	data RepositoryConfAuthData,
	repositoryID string,
) string {
	return fmt.Sprintf(`
	resource "cyral_repository_conf_auth" "my-repository-conf-auth" {
		repository_id = %s
		allow_native_auth = %t
		client_tls = "%s"
		identity_provider = "tf_test_conf_auth_okta"
		repo_tls = "%s"
	}`, repositoryID, data.AllowNativeAuth, data.ClientTLS, data.RepoTLS)
}
