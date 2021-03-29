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
		Host: "local-repo-binding.local",
		Port: 3333,
	},
}

var updatedRepositoryBindingConfig RepoBindingData = RepoBindingData{
	SidecarID:    "2",
	RepositoryID: "2",
	Enabled:      true,
	Listener: Listener{
		Host: "local-repo-binding-update.local",
		Port: 3334,
	},
}

func TestAccRepositoryBindingResource(t *testing.T) {
	testConfig, testFunc := setupRepositoryBindingTest(initialRepositoryBindingConfig)
	testUpdateConfig, testUpdateFunc := setupRepositoryBindingTest(updatedRepositoryBindingConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
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

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "enabled", fmt.Sprintf("%t", integrationData.Enabled)),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "repository_id", integrationData.RepositoryID),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "sidecar_id", integrationData.SidecarID),
		resource.TestCheckResourceAttr("cyral_repository_binding.repo_binding", "listener_port", fmt.Sprintf("%d", integrationData.Listener.Port)),
	)

	return configuration, testFunction
}

func formatRepoBindingDataIntoConfig(data RepoBindingData) string {
	return `
	resource "cyral_repository" "test_repo_binding_repository" {	
		type  = "mongodb"
		host  = "mongodb.cyral.com"
		port  = 27017
		name  = "test_repo_binding_sidecar-test"
	}
	
	resource "cyral_sidecar" "test_repo_binding_sidecar" {
		name = "test_repo_binding_sidecar-test"
		deployment_method = "cloudFormation"
	}

	resource "cyral_repository_binding" "repo_binding" {
		enabled = "false"
		repository_id = cyral_repository.test_repo_binding_repository.id
		listener_port = cyral_repository.test_repo_binding_repository.port
		sidecar_id    = cyral_sidecar.test_repo_binding_sidecar.id
		listener_host = "0.0.0.0"
	}`
}
