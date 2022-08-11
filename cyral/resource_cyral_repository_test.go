package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryResourceName = "repository"
)

var initialRepoConfig RepoData = RepoData{
	Name:     accTestName(repositoryResourceName, "repo"),
	Host:     "mongo.local",
	Port:     3333,
	RepoType: "mongodb",
	Labels:   []string{"rds", "us-east-2"},
}

var updatedRepoConfig RepoData = RepoData{
	Name:     accTestName(repositoryResourceName, "repo-updated"),
	Host:     "mongo-updated.local",
	Port:     3334,
	RepoType: "mongodb",
	Labels:   []string{"rds", "us-east-1"},
}

var replicaSetRepoConfig RepoData = RepoData{
	Name:                accTestName(repositoryResourceName, "repo-replica-set"),
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
	testConfig, testFunc := setupRepositoryTest(
		initialRepoConfig, "update_test")
	testUpdateConfig, testUpdateFunc := setupRepositoryTest(
		updatedRepoConfig, "update_test")
	testReplicaSetConfig, testReplicaSetFunc := setupRepositoryTest(
		replicaSetRepoConfig, "replica_config_test")

	// Should use name of the last resource created.
	importTestResourceName := "cyral_repository.replica_config_test"

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      importTestResourceName,
			},
		},
	})
}

func setupRepositoryTest(repoData RepoData, resName string) (string, resource.TestCheckFunc) {
	configuration := formatRepoDataIntoConfig(repoData, resName)

	resourceFullName := fmt.Sprintf("cyral_repository.%s", resName)

	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceFullName, "type", repoData.RepoType),
		resource.TestCheckResourceAttr(resourceFullName, "host", repoData.Host),
		resource.TestCheckResourceAttr(resourceFullName, "port", fmt.Sprintf("%d", repoData.Port)),
		resource.TestCheckResourceAttr(resourceFullName, "name", repoData.Name),
		resource.TestCheckResourceAttr(resourceFullName, "labels.#", "2"),
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

func formatRepoDataIntoConfig(data RepoData, resName string) string {
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
	resource "cyral_repository" "%s" {
		type  = "%s"
		host  = "%s"
		port  = %d
		name  = "%s"
		labels = %s
		%s
	}`, resName, data.RepoType, data.Host,
		data.Port, data.Name, listToStr(data.Labels), propertiesStr)
}
