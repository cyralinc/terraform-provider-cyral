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

var replicaSetRepoConfig RepoData = RepoData{
	Name:                "repo-test-replica-set",
	Host:                "mongo-cluster.local",
	Port:                27017,
	RepoType:            "mongodb",
	Labels:              []string{"rds", "us-east-1"},
	MaxAllowedListeners: 2,
	Properties: &RepositoryProperties{
		MongoDBReplicaSetName: "replica-set-1",
		MongoDBServerType:     mongodbReplicaSetServerType,
	},
}

func TestAccRepositoryResource(t *testing.T) {
	testConfig, testFunc := setupRepositoryTest(initialRepoConfig)
	testUpdateConfig, testUpdateFunc := setupRepositoryTest(updatedRepoConfig)
	testReplicaSetConfig, testReplicaSetFunc := setupRepositoryTest(replicaSetRepoConfig)

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
			{
				Config: testReplicaSetConfig,
				Check:  testReplicaSetFunc,
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

	if repoData.IsReplicaSet() {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("cyral_repository.test_repo_repository",
				"properties.0.mongodb_replica_set.0.max_nodes", fmt.Sprintf("%d",
					repoData.MaxAllowedListeners)),

			resource.TestCheckResourceAttr("cyral_repository.test_repo_repository",
				"properties.0.mongodb_replica_set.0.replica_set_id",
				repoData.Properties.MongoDBReplicaSetName),
		}...)
	}

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return configuration, testFunction
}

func formatRepoDataIntoConfig(data RepoData) string {
	var propertiesStr string
	if data.Properties != nil {
		properties := data.Properties

		var rsetStr string
		if data.IsReplicaSet() {
			rsetStr = fmt.Sprintf(`
			mongodb_replica_set {
				max_nodes = %d
				replica_set_id = "%s"
			}`, data.MaxAllowedListeners, properties.MongoDBReplicaSetName)
		}

		propertiesStr = fmt.Sprintf(`
		properties {%s
		}`, rsetStr)
	}

	return fmt.Sprintf(`
	resource "cyral_repository" "test_repo_repository" {
		type  = "%s"
		host  = "%s"
		port  = %d
		name  = "%s"
		labels = [%s]
		%s
	}`, data.RepoType, data.Host, data.Port, data.Name,
		formatAttributes(data.Labels), propertiesStr)
}
