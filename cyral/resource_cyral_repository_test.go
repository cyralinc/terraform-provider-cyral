package cyral

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
			// Remove replica set config again to see if the replica
			// set properties are actually removed from state. We
			// had a bug in the past where this test would not pass.
			{
				Config: testEmptyPropertiesConfig,
				Check:  testEmptyPropertiesFunc,
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
			"labels.#", fmt.Sprintf("%d", len(repoData.Labels))),
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
	} else if repoData.Properties != nil {
		// If we have an empty properties block: `properties {}`, we
		// need to check that the replica set attributes are not set. We
		// had a bug where this check would not pass in the past.
		checkFuncs = append(checkFuncs,
			func(s *terraform.State) error {
				resName := "cyral_repository.test_repo_repository"
				res, ok := s.RootModule().Resources[resName]
				if !ok {
					return fmt.Errorf("not found: %s", resName)
				}
				maxNodes := res.Primary.Attributes["properties.0.mongodb_replica_set.0.max_nodes"]
				replicaSetID := res.Primary.Attributes["properties.0.mongodb_replica_set.0.replica_set_id"]
				if maxNodes != "" || replicaSetID != "" {
					return fmt.Errorf("expected replica set attributes to " +
						"be unset for empty properties block")
				}
				return nil
			},
		)
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

	config := fmt.Sprintf(`
	resource "cyral_repository" "test_repo_repository" {
		type  = "%s"
		host  = "%s"
		port  = %d
		name  = "%s"
		labels = [%s]
		%s
	}`, data.RepoType, data.Host, data.Port, data.Name,
		formatAttributes(data.Labels), propertiesStr)

	log.Printf("[DEBUG] Config: %s\n", config)

	return config
}
