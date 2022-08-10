package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialConfig RepoBindingData = RepoBindingData{
	Listener: Listener{
		Port: 1234,
	},
}

var updatedConfig RepoBindingData = RepoBindingData{
	Listener: Listener{
		Host: "host-updated.com",
		Port: 4321,
	},
	Enabled:                   true,
	SidecarAsIdPAccessGateway: false,
}

func TestAccRepositoryBindingResource(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryBindingConfig_DefaultValues(),
				Check:  testAccRepositoryBindingCheck_DefaultValues(),
			},
			{
				Config: testAccRepositoryBindingConfig_UpdatedIDs(),
				Check:  testAccRepositoryBindingCheck_UpdatedIDs(),
			},
			{
				Config: testAccRepositoryBindingConfig_AccessGatewayEnabled(),
				Check:  testAccRepositoryBindingCheck_AccessGatewayEnabled(),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_repository_binding.repo_binding",
			},
		},
	})
}

func testAccRepositoryBindingConfig_DefaultValues() string {
	var config string
	config += formatBasicSidecarIntoConfig(
		"sidecar_1",
		accTestName("repository-binding", "sidecar-1"),
		"cloudFormation",
	)
	config += formatBasicRepositoryIntoConfig(
		"repository_1",
		accTestName("repository-binding", "repository-1"),
		"mongodb",
		"mongodb.cyral.com",
		27017,
	)
	config += fmt.Sprintf(`
	resource "cyral_repository_binding" "repo_binding" {
		sidecar_id    = cyral_sidecar.sidecar_1.id
		repository_id = cyral_repository.repository_1.id
		listener_port = %d
	}`, initialConfig.Listener.Port)
	return config
}

func testAccRepositoryBindingCheck_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "repository_id",
			"cyral_repository.repository_1", "id"),
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "sidecar_id",
			"cyral_sidecar.sidecar_1", "id"),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "listener_port",
			fmt.Sprintf("%d", initialConfig.Listener.Port)),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding",
			"listener_host", "0.0.0.0"),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding",
			"enabled", "true"),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding",
			"sidecar_as_idp_access_gateway", "false"),
	)
}

func testAccRepositoryBindingConfig_UpdatedIDs() string {
	var config string
	config += formatBasicSidecarIntoConfig(
		"sidecar_2",
		accTestName("repository-binding", "sidecar-2"),
		"cloudFormation",
	)
	config += formatBasicRepositoryIntoConfig(
		"repository_2",
		accTestName("repository-binding", "repository-2"),
		"mongodb",
		"mongodb.cyral.com",
		27017,
	)
	config += fmt.Sprintf(`
	resource "cyral_repository_binding" "repo_binding" {
		sidecar_id    = cyral_sidecar.sidecar_2.id
		repository_id = cyral_repository.repository_2.id
		listener_port = %d
		listener_host = "%s"
		enabled = %t
		sidecar_as_idp_access_gateway = %t
	}`, updatedConfig.Listener.Port, updatedConfig.Listener.Host,
		updatedConfig.Enabled, updatedConfig.SidecarAsIdPAccessGateway)
	return config
}

func testAccRepositoryBindingCheck_UpdatedIDs() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "repository_id",
			"cyral_repository.repository_2", "id"),
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "sidecar_id",
			"cyral_sidecar.sidecar_2", "id"),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "listener_port",
			fmt.Sprintf("%d", updatedConfig.Listener.Port)),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding",
			"listener_host", updatedConfig.Listener.Host),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding",
			"enabled", "true"),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding",
			"sidecar_as_idp_access_gateway", "false"),
	)
}

func testAccRepositoryBindingConfig_AccessGatewayEnabled() string {
	var config string
	config += formatBasicSidecarIntoConfig(
		"sidecar_2",
		accTestName("repository-binding", "sidecar-2"),
		"cloudFormation",
	)
	config += formatBasicRepositoryIntoConfig(
		"repository_2",
		accTestName("repository-binding", "repository-2"),
		"mongodb",
		"mongodb.cyral.com",
		27017,
	)
	config += fmt.Sprintf(`
	resource "cyral_repository_binding" "repo_binding" {
		sidecar_id    = cyral_sidecar.sidecar_2.id
		repository_id = cyral_repository.repository_2.id
		listener_port = %d
		listener_host = "%s"
		enabled = %t
		sidecar_as_idp_access_gateway = true
	}`, updatedConfig.Listener.Port, updatedConfig.Listener.Host,
		updatedConfig.Enabled)
	return config
}

func testAccRepositoryBindingCheck_AccessGatewayEnabled() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "repository_id",
			"cyral_repository.repository_2", "id"),
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "sidecar_id",
			"cyral_sidecar.sidecar_2", "id"),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "listener_port",
			fmt.Sprintf("%d", updatedConfig.Listener.Port)),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding",
			"listener_host", updatedConfig.Listener.Host),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding",
			"enabled", "true"),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding",
			"sidecar_as_idp_access_gateway", "true"),
	)
}
