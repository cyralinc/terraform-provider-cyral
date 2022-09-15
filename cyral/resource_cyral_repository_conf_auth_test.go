package cyral

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryConfAuthResourceName = "repository-conf-auth"
)

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

func repositoryConfAuthConfigNetworkShield() RepositoryConfAuthData {
	return RepositoryConfAuthData{
		AllowNativeAuth:            false,
		ClientTLS:                  "enable",
		RepoTLS:                    "disable",
		EnableNetworkAccessControl: true,
	}
}

func TestAccRepositoryConfAuthResource(t *testing.T) {
	testConfig, testFunc := setupRepositoryConfAuthTest(initialRepositoryConfAuthConfig())
	testUpdate1Config, testUpdate1Func := setupRepositoryConfAuthTest(update1RepositoryConfAuthConfig())
	testUpdate2Config, testUpdate2Func := setupRepositoryConfAuthTest(update2RepositoryConfAuthConfig())
	testNetworkShieldConfig, testNetworkShieldFunc := setupRepositoryConfAuthTest(repositoryConfAuthConfigNetworkShield())

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
			{
				Config: testNetworkShieldConfig,
				Check:  testNetworkShieldFunc,
			},
			// TODO: add import test -aholmquist 2022-08-05
		},
	})
}

func setupRepositoryConfAuthTest(repositoryConf RepositoryConfAuthData) (string, resource.TestCheckFunc) {
	var configuration string
	configuration += formatBasicRepositoryIntoConfig(
		basicRepositoryResName,
		accTestName(repositoryConfAuthResourceName, "repository"),
		"mysql",
		"http://mysql.local/",
		3306,
	)

	resourceName := "my-repository-conf-auth"
	resourceFullName := fmt.Sprintf("cyral_repository_conf_auth.%s", resourceName)

	configuration += formatRepositoryConfAuthDataIntoConfig(
		resourceName, repositoryConf, basicRepositoryID)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"allow_native_auth", fmt.Sprintf("%t", repositoryConf.AllowNativeAuth)),
		resource.TestCheckResourceAttr(resourceFullName,
			"client_tls", repositoryConf.ClientTLS),
		resource.TestCheckResourceAttr(resourceFullName,
			"repo_tls", repositoryConf.RepoTLS),
		resource.TestCheckResourceAttr(resourceFullName,
			"enable_network_access_control", strconv.FormatBool(repositoryConf.EnableNetworkAccessControl)),
	)

	return configuration, testFunction
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
		enable_network_access_control = %t
	}`, resName, repositoryID, data.AllowNativeAuth, data.ClientTLS,
		data.RepoTLS, data.EnableNetworkAccessControl)
}
