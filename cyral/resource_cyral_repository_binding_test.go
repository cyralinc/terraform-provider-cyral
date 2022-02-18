package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialRepositoryBindingConfig RepoBindingData = RepoBindingData{
	SidecarID:    "1",
	RepositoryID: "1",
	Enabled:      false,
	Listener: Listener{
		Host: "host.com",
		Port: 3333,
	},
}

var updatedRepositoryBindingConfig RepoBindingData = RepoBindingData{
	SidecarID:    "2",
	RepositoryID: "2",
	Enabled:      true,
	Listener: Listener{
		Host: "host-updated.com",
		Port: 3334,
	},
}

func TestAccRepositoryBindingResource(t *testing.T) {
	testConfig, testFunc := setupRepositoryBindingTest(initialRepositoryBindingConfig)
	testUpdateConfig, testUpdateFunc := setupRepositoryBindingTest(updatedRepositoryBindingConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
		},
	})
}

func setupRepositoryBindingTest(integrationData RepoBindingData) (string, resource.TestCheckFunc) {
	configuration := formatRepoBindingDataIntoConfig(integrationData)

	sidecarResource := fmt.Sprintf("cyral_sidecar.test_repo_binding_sidecar_%s",
		integrationData.SidecarID)
	repositoryResource := fmt.Sprintf("cyral_repository.test_repo_binding_repository_%s",
		integrationData.RepositoryID)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "enabled",
			fmt.Sprintf("%t", integrationData.Enabled)),
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "repository_id",
			repositoryResource, "id"),
		resource.TestCheckResourceAttrPair("cyral_repository_binding.repo_binding", "sidecar_id",
			sidecarResource, "id"),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "listener_host",
			integrationData.Listener.Host),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "listener_port",
			fmt.Sprintf("%d", integrationData.Listener.Port)),
	)

	return configuration, testFunction
}

func formatRepoBindingDataIntoConfig(data RepoBindingData) string {
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
		enabled = %t
		sidecar_id    = cyral_sidecar.test_repo_binding_sidecar_%s.id
		repository_id = cyral_repository.test_repo_binding_repository_%s.id
		listener_host = "%s"
		listener_port = %d
	}`, data.Enabled, data.SidecarID, data.RepositoryID, data.Listener.Host, data.Listener.Port)
}
