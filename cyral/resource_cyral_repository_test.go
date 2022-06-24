package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialRepoConfig RepoData = RepoData{
	Name:     "repo-test",
	Host:     "mongo.local",
	Port:     3333,
	RepoType: "mongodb",
	Labels:   []string{"rds", "us-east-2"},
}

var updatedRepoConfig RepoData = RepoData{
	Name:     "repo-test-updated",
	Host:     "mongo-updated.local",
	Port:     3334,
	RepoType: "mongodb",
	Labels:   []string{"rds", "us-east-1"},
}

func TestAccRepositoryResource(t *testing.T) {
	testConfig, testFunc := setupRepositoryTest(initialRepoConfig)
	testUpdateConfig, testUpdateFunc := setupRepositoryTest(updatedRepoConfig)

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

func setupRepositoryTest(integrationData RepoData) (string, resource.TestCheckFunc) {
	configuration := formatRepoDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository", "type", integrationData.RepoType),
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository", "host", integrationData.Host),
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository", "port", fmt.Sprintf("%d", integrationData.Port)),
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository", "labels.#", "2"),
	)

	return configuration, testFunction
}

func formatRepoDataIntoConfig(data RepoData) string {
	return fmt.Sprintf(`
	resource "cyral_repository" "test_repo_repository" {
		type  = "%s"
		host  = "%s"
		port  = %d
		name  = "%s"
		labels = [%s]
	}`, data.RepoType, data.Host, data.Port, data.Name, formatAttributes(data.Labels))
}
