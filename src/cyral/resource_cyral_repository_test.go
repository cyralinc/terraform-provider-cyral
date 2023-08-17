package cyral

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryResourceName = "repository"
)

var (
	initialRepoConfig = RepoInfo{
		Name:   utils.AccTestName(repositoryResourceName, "repo"),
		Type:   MongoDB,
		Labels: []string{"rds", "us-east-2"},
		RepoNodes: []*RepoNode{
			{
				Host: "mongo.local",
				Port: 3333,
			},
		},
		MongoDBSettings: &MongoDBSettings{
			ServerType: Standalone,
		},
	}

	updatedRepoConfig = RepoInfo{
		Name:   utils.AccTestName(repositoryResourceName, "repo-updated"),
		Type:   MongoDB,
		Labels: []string{"rds", "us-east-1"},
		RepoNodes: []*RepoNode{
			{
				Host: "mongo.local",
				Port: 3334,
			},
		},
		MongoDBSettings: &MongoDBSettings{
			ServerType: Standalone,
		},
	}

	emptyConnDrainingConfig = RepoInfo{
		Name: utils.AccTestName(repositoryResourceName, "repo-empty-conn-draining"),
		Type: MongoDB,
		ConnParams: &ConnParams{
			ConnDraining: &ConnDraining{},
		},
		RepoNodes: []*RepoNode{
			{
				Host: "mongo-cluster.local",
				Port: 27017,
			},
		},
		MongoDBSettings: &MongoDBSettings{
			ServerType: Standalone,
		},
	}

	connDrainingConfig = RepoInfo{
		Name: utils.AccTestName(repositoryResourceName, "repo-conn-draining"),
		Type: MongoDB,
		ConnParams: &ConnParams{
			ConnDraining: &ConnDraining{
				Auto:     true,
				WaitTime: 20,
			},
		},
		RepoNodes: []*RepoNode{
			{
				Host: "mongo-cluster.local",
				Port: 27017,
			},
		},
		MongoDBSettings: &MongoDBSettings{
			ServerType: Standalone,
		},
	}

	mixedMultipleNodesConfig = RepoInfo{
		Name: utils.AccTestName(repositoryResourceName, "repo-mixed-multi-node"),
		Type: MongoDB,
		ConnParams: &ConnParams{
			ConnDraining: &ConnDraining{
				Auto:     true,
				WaitTime: 20,
			},
		},
		RepoNodes: []*RepoNode{
			{
				Name: "node1",
				Host: "mongo-cluster.local.node1",
				Port: 27017,
			},
			{
				Name: "node2",
				Host: "mongo-cluster.local.node2",
				Port: 27017,
			},
			{
				Name: "node3",
				Host: "mongo-cluster.local.node3",
				Port: 27017,
			},
			{
				Dynamic: true,
			},
			{
				Name:    "node5",
				Dynamic: true,
			},
		},
		MongoDBSettings: &MongoDBSettings{
			ReplicaSetName: "some-replica-set",
			ServerType:     ReplicaSet,
		},
	}

	allRepoNodesAreDynamic = RepoInfo{
		Name: utils.AccTestName(repositoryResourceName, "repo-all-repo-nodes-are-dynamic"),
		Type: "mongodb",
		RepoNodes: []*RepoNode{
			{
				Dynamic: true,
			},
			{
				Dynamic: true,
			},
			{
				Dynamic: true,
			},
		},
		MongoDBSettings: &MongoDBSettings{
			ReplicaSetName: "myReplicaSet",
			ServerType:     "replicaset",
			SRVRecordName:  "mySRVRecord",
		},
	}
)

func TestAccRepositoryResource(t *testing.T) {
	initial := setupRepositoryTest(
		initialRepoConfig, "update_test")
	update := setupRepositoryTest(
		updatedRepoConfig, "update_test")
	connDrainingEmpty := setupRepositoryTest(
		emptyConnDrainingConfig, "conn_draining_empty_test")
	connDraining := setupRepositoryTest(
		connDrainingConfig, "conn_draining_test")
	allDynamic := setupRepositoryTest(
		allRepoNodesAreDynamic, "all_repo_nodes_are_dynamic")

	multiNode := setupRepositoryTest(
		mixedMultipleNodesConfig, "multi_node_test")

	// Should use name of the last resource created.
	importTest := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "cyral_repository.multi_node_test",
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			initial,
			update,
			connDrainingEmpty,
			connDraining,
			allDynamic,
			multiNode,
			importTest,
		},
	})
}

func setupRepositoryTest(repo RepoInfo, resName string) resource.TestStep {
	return resource.TestStep{
		Config: repoAsConfig(repo, resName),
		Check:  repoCheckFuctions(repo, resName),
	}
}

func repoCheckFuctions(repo RepoInfo, resName string) resource.TestCheckFunc {
	resourceFullName := fmt.Sprintf("cyral_repository.%s", resName)

	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceFullName,
			"type", repo.Type),
		resource.TestCheckResourceAttr(resourceFullName,
			"name", repo.Name),
		resource.TestCheckResourceAttr(resourceFullName,
			"labels.#", fmt.Sprintf("%d", len(repo.Labels))),
	}

	for i, label := range repo.Labels {
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr(resourceFullName,
				fmt.Sprintf("labels.%d", i), label))
	}

	for i, node := range repo.RepoNodes {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceFullName,
				fmt.Sprintf("repo_node.%d.name", i), node.Name),
			resource.TestCheckResourceAttr(resourceFullName,
				fmt.Sprintf("repo_node.%d.host", i), node.Host),
			resource.TestCheckResourceAttr(resourceFullName,
				fmt.Sprintf("repo_node.%d.port", i),
				strconv.Itoa(int(node.Port))),
			resource.TestCheckResourceAttr(resourceFullName,
				fmt.Sprintf("repo_node.%d.dynamic", i),
				strconv.FormatBool(node.Dynamic)),
		}...)
	}

	if repo.ConnParams != nil {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceFullName,
				"connection_draining.0.auto",
				strconv.FormatBool(repo.ConnParams.ConnDraining.Auto)),

			resource.TestCheckResourceAttr(resourceFullName,
				"connection_draining.0.wait_time",
				strconv.Itoa(int(repo.ConnParams.ConnDraining.WaitTime))),
		}...)
	}

	if repo.MongoDBSettings != nil {
		if repo.MongoDBSettings.ServerType == ReplicaSet {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceFullName,
					"mongodb_settings.0.replica_set_name",
					repo.MongoDBSettings.ReplicaSetName),
			}...)
		}
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceFullName,
				"mongodb_settings.0.server_type",
				repo.MongoDBSettings.ServerType),

			resource.TestCheckResourceAttr(resourceFullName,
				"mongodb_settings.0.srv_record_name",
				repo.MongoDBSettings.SRVRecordName),
		}...)
	}

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func repoAsConfig(repo RepoInfo, resName string) string {
	config := fmt.Sprintf(`
	resource "cyral_repository" "%s" {
		type  = "%s"
		name  = "%s"
		labels = %s`, resName, repo.Type, repo.Name, utils.ListToStr(repo.Labels))

	if repo.ConnParams != nil {
		config += fmt.Sprintf(`
		connection_draining {
			auto = %s
			wait_time = %d
		}`, strconv.FormatBool(repo.ConnParams.ConnDraining.Auto),
			repo.ConnParams.ConnDraining.WaitTime,
		)
	}

	if repo.MongoDBSettings != nil {
		replicaSet := "null"
		serverType := "null"
		srvRecordName := "null"
		if repo.MongoDBSettings.ReplicaSetName != "" {
			replicaSet = fmt.Sprintf(`"%s"`, repo.MongoDBSettings.ReplicaSetName)
		}
		if repo.MongoDBSettings.ServerType != "" {
			serverType = fmt.Sprintf(`"%s"`, repo.MongoDBSettings.ServerType)
		}
		if repo.MongoDBSettings.SRVRecordName != "" {
			srvRecordName = fmt.Sprintf(`"%s"`, repo.MongoDBSettings.SRVRecordName)
		}
		config += fmt.Sprintf(`
		mongodb_settings {
			replica_set_name = %s
			server_type = %s
			srv_record_name = %s
		}`,
			replicaSet,
			serverType,
			srvRecordName,
		)
	}

	for _, node := range repo.RepoNodes {
		name, host := "null", "null"
		if node.Name != "" {
			name = fmt.Sprintf(`"%s"`, node.Name)
		}
		if node.Host != "" {
			host = fmt.Sprintf(`"%s"`, node.Host)
		}
		config += fmt.Sprintf(`
		repo_node {
			name = %s
			host = %s
			port = %d
			dynamic = %s
		}`, name, host, node.Port, strconv.FormatBool(node.Dynamic))
	}

	config += `
	}`

	return config
}
