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

var emptyPropertiesRepoConfig RepoData = RepoData{
	Name:       "repo-test-empty-properties",
	Host:       "mongo-cluster.local",
	Port:       27017,
	RepoType:   "mongodb",
	Properties: &RepositoryProperties{},
}

var replicaSetRepoConfig RepoData = RepoData{
	Name:                "repo-test-replica-set",
	Host:                "mongo-cluster.local",
	Port:                27017,
	RepoType:            "mongodb",
	MaxAllowedListeners: 2,
	Properties: &RepositoryProperties{
		MongoDBReplicaSetName: "replica-set-1",
		MongoDBServerType:     mongodbReplicaSetServerType,
	},
}

func TestAccRepositoryResource(t *testing.T) {
	testConfig, testFunc := setupRepositoryTest(initialRepoConfig)
	testUpdateConfig, testUpdateFunc := setupRepositoryTest(updatedRepoConfig)
	testEmptyPropertiesConfig, testEmptyPropertiesFunc := setupRepositoryTest(emptyPropertiesRepoConfig)
	testReplicaSetConfig, testReplicaSetFunc := setupRepositoryTest(replicaSetRepoConfig)

	// Should use name of the last resource created.
	importTestResourceName := repositoryConfigResourceFullName(replicaSetRepoConfig.Name)

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
				Config: testEmptyPropertiesConfig,
				Check:  testEmptyPropertiesFunc,
			},
			{
				Config: testReplicaSetConfig,
				Check:  testReplicaSetFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      importTestResourceName,
			},
		},
	})
}

func setupRepositoryTest(repoData RepoData) (string, resource.TestCheckFunc) {
	configuration := formatRepoDataIntoConfig(repoData)

	resourceFullName := repositoryConfigResourceFullName(repoData.Name)

	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceFullName, "type", repoData.RepoType),
		resource.TestCheckResourceAttr(resourceFullName, "host", repoData.Host),
		resource.TestCheckResourceAttr(resourceFullName, "port", fmt.Sprintf("%d", repoData.Port)),
		resource.TestCheckResourceAttr(resourceFullName, "name", repoData.Name),
		resource.TestCheckResourceAttr(resourceFullName, "labels.#", fmt.Sprintf("%d", len(repoData.Labels))),
	}

	if repoData.IsReplicaSet() {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceFullName,
				"properties.0.mongodb_replica_set.0.max_nodes", fmt.Sprintf("%d",
					repoData.MaxAllowedListeners)),

			resource.TestCheckResourceAttr(resourceFullName,
				"properties.0.mongodb_replica_set.0.replica_set_id",
				repoData.Properties.MongoDBReplicaSetName),
		}...)
	}

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return configuration, testFunction
}

func repositoryConfigResourceFullName(repoName string) string {
	return fmt.Sprintf("cyral_repository.%s", repositoryConfigResourceName(repoName))
}

func repositoryConfigResourceName(repoName string) string {
	return fmt.Sprintf("test_repository_%s", repoName)
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

	config := fmt.Sprintf(`
	resource "cyral_repository" "%s" {
		type  = "%s"
		host  = "%s"
		port  = %d
		name  = "%s"
		labels = [%s]
		%s
	}`, repositoryConfigResourceName(data.Name), data.RepoType, data.Host,
		data.Port, data.Name, formatAttributes(data.Labels), propertiesStr)

	return config
}
