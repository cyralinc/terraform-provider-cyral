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

var initialRepoConfigReplicaSet RepoData = RepoData{
	Name:     "repo-test-replica-set",
	Host:     "host-1",
	Port:     27017,
	RepoType: "mongodb",
	Labels:   []string{"rds", "us-east-1"},
	Properties: &RepositoryProperties{
		MaxNodes:              "2",
		MongoDBReplicaSetName: "replica-set-1",
		MongoDBServerType:     "replicaset",
	},
}

var updatedRepoConfigReplicaSet RepoData = RepoData{
	Name:     "repo-test-replica-set",
	Host:     "host-1",
	Port:     27017,
	RepoType: "mongodb",
	Labels:   []string{"rds", "us-east-1"},
	Properties: &RepositoryProperties{
		MaxNodes:              "2",
		MongoDBReplicaSetName: "replica-set-2",
		MongoDBServerType:     "replicaset",
	},
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

func TestAccRepositoryResource_MongoDBReplicaSet(t *testing.T) {
	testConfig, testFunc := setupRepositoryTest(initialRepoConfigReplicaSet)
	testUpdateConfig, testUpdateFunc := setupRepositoryTest(updatedRepoConfigReplicaSet)

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

func setupRepositoryTest(repoData RepoData) (string, resource.TestCheckFunc) {
	configuration := formatRepoDataIntoConfig(repoData)

	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository",
			"type", repoData.RepoType),
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository",
			"host", repoData.Host),
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository",
			"port", fmt.Sprintf("%d", repoData.Port)),
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository",
			"name", repoData.Name),
		resource.TestCheckResourceAttr("cyral_repository.test_repo_repository",
			"labels.#", "2"),
	}

	if repoReplicaSetEnabled(repoData) {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("cyral_repository.test_repo_repository",
				"properties.0.replica_set.0.max_nodes", repoData.Properties.MaxNodes),
			resource.TestCheckResourceAttr("cyral_repository.test_repo_repository",
				"properties.0.replica_set.0.MongoDBReplicaSetName",
				repoData.Properties.MongoDBReplicaSetName),
		}...)
	}

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return configuration, testFunction
}

func formatRepoDataIntoConfig(data RepoData) string {
	base := fmt.Sprintf(`
	resource "cyral_repository" "test_repo_repository" {
		type  = "%s"
		host  = "%s"
		port  = %d
		name  = "%s"
		labels = [%s]`, data.RepoType, data.Host, data.Port, data.Name, formatAttibutes(data.Labels))
	propertiesStr := ""
	if data.Properties != nil {
		propertiesStr += `
		properties {`
		properties := data.Properties
		if repoReplicaSetEnabled(data) {
			propertiesStr += fmt.Sprintf(`
			replica_set {
				max_nodes = %s
				replica_set_id = %s
			}`, properties.MaxNodes, properties.MongoDBReplicaSetName)
		}
		propertiesStr += `
		}
`
	}
	completeConfig := base + propertiesStr + `
	}`
	return completeConfig
}

func repoReplicaSetEnabled(repoData RepoData) bool {
	properties := repoData.Properties
	return properties != nil && properties.MaxNodes != "" && properties.MongoDBReplicaSetName != ""
}
