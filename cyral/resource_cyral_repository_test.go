package cyral

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryResourceName = "repository"
)

// func repositoryTestSidecarConfig() string {
// 	return formatBasicSidecarIntoConfig(
// 		basicSidecarResName,
// 		accTestName(repositoryBindingResourceName, "repo-test-sidecar"),
// 		"docker",
// 	)
// }

var (
	initialRepoConfig = RepoInfo{
		Name:   accTestName(repositoryResourceName, "repo"),
		Type:   "mongodb",
		Labels: []string{"rds", "us-east-2"},
		RepoNodes: []*RepoNode{
			{
				Host: "mongo.local",
				Port: 3333,
			},
		},
	}

	updatedRepoConfig = RepoInfo{
		Name:   accTestName(repositoryResourceName, "repo-updated"),
		Type:   "mongodb",
		Labels: []string{"rds", "us-east-1"},
		RepoNodes: []*RepoNode{
			{
				Host: "mongo.local",
				Port: 3334,
			},
		},
	}

	emptyConnDrainingConfig = RepoInfo{
		Name: accTestName(repositoryResourceName, "repo-empty-conn-draining"),
		Type: "mongodb",
		ConnParams: &ConnParams{
			ConnDraining: &ConnDraining{},
		},
		RepoNodes: []*RepoNode{
			{
				Host: "mongo-cluster.local",
				Port: 27017,
			},
		},
	}

	connDrainingConfig = RepoInfo{
		Name: accTestName(repositoryResourceName, "repo-conn-draining"),
		Type: "mongodb",
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
	}

	accessGatewayEmptyConfig = RepoInfo{
		Name: accTestName(repositoryResourceName, "repo-with-access-gateway-empty"),
		Type: "mongodb",
		ConnParams: &ConnParams{
			ConnDraining: &ConnDraining{
				Auto:     true,
				WaitTime: 20,
			},
		},
		RepoNodes: []*RepoNode{
			{
				Name: "node1",
				Host: "mongo.local.node1",
				Port: 27017,
			},
			{
				Name: "node2",
				Host: "mongo.local.node2",
				Port: 27017,
			},
		},
		PreferredAccessGwBinding: &BindingKey{},
	}

	// accessGatewayConfig = RepoInfo{
	// 	Name: accTestName(repositoryResourceName, "repo-with-access-gateway"),
	// 	Type: "mongodb",
	// 	ConnParams: &ConnParams{
	// 		ConnDraining: &ConnDraining{
	// 			Auto:     true,
	// 			WaitTime: 20,
	// 		},
	// 	},
	// 	RepoNodes: []*RepoNode{
	// 		{
	// 			Name: "node1",
	// 			Host: "mongo.local.node1",
	// 			Port: 27017,
	// 		},
	// 		{
	// 			Name: "node2",
	// 			Host: "mongo.local.node2",
	// 			Port: 27017,
	// 		},
	// 	},
	// 	PreferredAccessGwBinding: &BindingKey{
	// 		BindingID: "some-binding-id",
	// 	},
	// }

	mixedMultipleNodesConfig = RepoInfo{
		Name: accTestName(repositoryResourceName, "repo-mixed-multi-node"),
		Type: "mongodb",
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
	accessGatewayEmpty := setupRepositoryTest(
		accessGatewayEmptyConfig, "access_gateway_empty_test")
	// accessGateway := setupRepositoryTest(
	// 	accessGatewayConfig, "access_gateway_test")

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
			accessGatewayEmpty,
			// accessGateway,
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
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceFullName,
				"mongodb_settings.0.replica_set_name",
				repo.MongoDBSettings.ReplicaSetName),

			resource.TestCheckResourceAttr(resourceFullName,
				"connection_draining.0.server_type",
				repo.MongoDBSettings.ServerType),
		}...)
	}

	// if repo.PreferredAccessGwBinding != nil {
	// 	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
	// 		resource.TestCheckResourceAttrPair(
	// 			resourceFullName, "preferred_access_gateway.0.sidecar_id",
	// 			fmt.Sprintf("cyral_sidecar.%s", basicSidecarResName), "id"),
	// 		resource.TestCheckResourceAttr(resourceFullName,
	// 			"preferred_access_gateway.0.binding_id",
	// 			repo.PreferredAccessGwBinding.BindingID),
	// 	}...)
	// }

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func repoAsConfig(repo RepoInfo, resName string) string {
	config := fmt.Sprintf(`
	resource "cyral_repository" "%s" {
		type  = "%s"
		name  = "%s"
		labels = %s`, resName, repo.Type, repo.Name, listToStr(repo.Labels))

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
		replicaSet, serverType := "null", "null"
		if repo.MongoDBSettings.ReplicaSetName != "" {
			replicaSet = fmt.Sprintf(
				`"%s"`, repo.MongoDBSettings.ReplicaSetName)
		}
		if repo.MongoDBSettings.ServerType != "" {
			serverType = fmt.Sprintf(
				`"%s"`, repo.MongoDBSettings.ServerType)
		}
		config += fmt.Sprintf(`
		mongodb_settings {
			replica_set_name = %s
			server_type = %s
		}`, replicaSet,
			serverType,
		)
	}

	// if repo.PreferredAccessGwBinding != nil {
	// 	bindingID := "null"
	// 	if repo.PreferredAccessGwBinding.BindingID != "" {
	// 		bindingID = fmt.Sprintf(
	// 			`"%s"`, repo.PreferredAccessGwBinding.BindingID)
	// 	}
	// 	config += fmt.Sprintf(`
	// 	preferred_access_gateway {
	// 		sidecar_id = %s
	// 		binding_id = %s
	// 	}`, basicSidecarID,
	// 		bindingID,
	// 	)
	// }

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
