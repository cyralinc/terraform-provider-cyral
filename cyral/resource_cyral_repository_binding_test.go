package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialConfig RepoBindingData = RepoBindingData{
	SidecarID:    "1",
	RepositoryID: "1",
	Listener: Listener{
		Port: 1234,
	},
}

var updatedConfig RepoBindingData = RepoBindingData{
	SidecarID:    "2",
	RepositoryID: "2",
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
	return fmt.Sprintf(`
	resource "cyral_sidecar" "test_repo_binding_sidecar_1" {
		name = "tf-provider-repo-binding-sidecar-1"
		deployment_method = "cloudFormation"
	}

	resource "cyral_repository" "test_repo_binding_repository_1" {
		name  = "tf-provider-repo-binding-repo-1"
		type  = "mongodb"
		host  = "mongodb.cyral.com"
		port  = 27017
	}

	resource "cyral_repository_binding" "repo_binding" {
		sidecar_id    = cyral_sidecar.test_repo_binding_sidecar_%s.id
		repository_id = cyral_repository.test_repo_binding_repository_%s.id
		listener_port = %d
	}`, initialConfig.SidecarID, initialConfig.RepositoryID, initialConfig.Listener.Port)
}

func testAccRepositoryBindingCheck_DefaultValues() resource.TestCheckFunc {
	sidecarResource := fmt.Sprintf("cyral_sidecar.test_repo_binding_sidecar_%s",
		initialConfig.SidecarID)
	repositoryResource := fmt.Sprintf("cyral_repository.test_repo_binding_repository_%s",
		initialConfig.RepositoryID)

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "repository_id",
			repositoryResource, "id"),
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "sidecar_id",
			sidecarResource, "id"),
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
	return fmt.Sprintf(`
	resource "cyral_sidecar" "test_repo_binding_sidecar_2" {
		name = "tf-provider-repo-binding-sidecar-2"
		deployment_method = "cloudFormation"
	}

	resource "cyral_repository" "test_repo_binding_repository_2" {
		name  = "tf-provider-repo-binding-repo-2"
		type  = "mongodb"
		host  = "mongodb.cyral.com"
		port  = 27017
	}

	resource "cyral_repository_binding" "repo_binding" {
		sidecar_id    = cyral_sidecar.test_repo_binding_sidecar_%s.id
		repository_id = cyral_repository.test_repo_binding_repository_%s.id
		listener_port = %d
		listener_host = "%s"
		enabled = %t
		sidecar_as_idp_access_gateway = %t
	}`, updatedConfig.SidecarID, updatedConfig.RepositoryID,
		updatedConfig.Listener.Port, updatedConfig.Listener.Host,
		updatedConfig.Enabled, updatedConfig.SidecarAsIdPAccessGateway)
}

func testAccRepositoryBindingCheck_UpdatedIDs() resource.TestCheckFunc {
	sidecarResource := fmt.Sprintf("cyral_sidecar.test_repo_binding_sidecar_%s",
		updatedConfig.SidecarID)
	repositoryResource := fmt.Sprintf("cyral_repository.test_repo_binding_repository_%s",
		updatedConfig.RepositoryID)

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "repository_id",
			repositoryResource, "id"),
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "sidecar_id",
			sidecarResource, "id"),
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
	return fmt.Sprintf(`
	resource "cyral_sidecar" "test_repo_binding_sidecar_2" {
		name = "tf-provider-repo-binding-sidecar-2"
		deployment_method = "cloudFormation"
	}

	resource "cyral_repository" "test_repo_binding_repository_2" {
		name  = "tf-provider-repo-binding-repo-2"
		type  = "mongodb"
		host  = "mongodb.cyral.com"
		port  = 27017
	}

	resource "cyral_repository_binding" "repo_binding" {
		sidecar_id    = cyral_sidecar.test_repo_binding_sidecar_%s.id
		repository_id = cyral_repository.test_repo_binding_repository_%s.id
		listener_port = %d
		listener_host = "%s"
		enabled = %t
		sidecar_as_idp_access_gateway = true
	}`, updatedConfig.SidecarID, updatedConfig.RepositoryID,
		updatedConfig.Listener.Port, updatedConfig.Listener.Host,
		updatedConfig.Enabled)
}

func testAccRepositoryBindingCheck_AccessGatewayEnabled() resource.TestCheckFunc {
	sidecarResource := fmt.Sprintf("cyral_sidecar.test_repo_binding_sidecar_%s",
		updatedConfig.SidecarID)
	repositoryResource := fmt.Sprintf("cyral_repository.test_repo_binding_repository_%s",
		updatedConfig.RepositoryID)

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "repository_id",
			repositoryResource, "id"),
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "sidecar_id",
			sidecarResource, "id"),
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
