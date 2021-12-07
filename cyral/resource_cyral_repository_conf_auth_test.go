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

func TestAccRepositoryConfAuthDataResource(t *testing.T) {
	testConfig, testFunc := setupRepositoryConfAuthTest(initialRepositoryConfAuthConfig)
	testUpdate1Config, testUpdate1Func := setupRepositoryConfAuthTest(update1RepositoryConfAuthConfig)
	testUpdate2Config, testUpdate2Func := setupRepositoryConfAuthTest(update2RepositoryConfAuthConfig)

	resource.Test(t, resource.TestCase{
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
		},
	})
}

func setupRepositoryConfAuthTest(integrationData RepositoryConfAuthData) (string, resource.TestCheckFunc) {
	configuration := formatRepositoryConfAuthDataDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_repository_conf_auth.my-repository-conf-auth",
			"allow_native_auth", fmt.Sprintf("%t", integrationData.AllowNativeAuth)),
		resource.TestCheckResourceAttr("cyral_repository_conf_auth.my-repository-conf-auth",
			"client_tls", integrationData.ClientTLS),
		resource.TestCheckResourceAttr("cyral_repository_conf_auth.my-repository-conf-auth",
			"repo_tls", integrationData.RepoTLS),
	)

	return configuration, testFunction
}

func formatRepositoryConfAuthDataDataIntoConfig(data RepositoryConfAuthData) string {
	return fmt.Sprintf(`
	resource "cyral_repository" "tf_test_repository" {
		type = "mysql"
		host = "http://mysql.local/"
		port = 3306
		name = "tf-test-mysql"
	  }

	resource "cyral_repository_conf_auth" "my-repository-conf-auth" {
		repository_id = cyral_repository.tf_test_repository.id
		allow_native_auth = %t
		client_tls = "%s"
		identity_provider = "tf_test_conf_auth_okta"
		repo_tls = "%s"
	}`, data.AllowNativeAuth, data.ClientTLS, data.RepoTLS)
}
