package repository_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	initialRepoConfig = repository.RepoInfo{
		Name:   utils.AccTestName(utils.RepositoryResourceName, "repo"),
		Type:   repository.MongoDB,
		Labels: repository.Labels{"rds", "us-east-2"},
		RepoNodes: repository.RepoNodes{
			{
				Host: "mongo.local",
				Port: 3333,
			},
		},
		MongoDBSettings: &repository.MongoDBSettings{
			ServerType: repository.Standalone,
		},
	}

	updatedRepoConfig = repository.RepoInfo{
		Name:   utils.AccTestName(utils.RepositoryResourceName, "repo-updated"),
		Type:   repository.MongoDB,
		Labels: repository.Labels{"rds", "us-east-1"},
		RepoNodes: repository.RepoNodes{
			{
				Host: "mongo.local",
				Port: 3334,
			},
		},
		MongoDBSettings: &repository.MongoDBSettings{
			ServerType: repository.Standalone,
		},
	}

	emptyConnDrainingConfig = repository.RepoInfo{
		Name: utils.AccTestName(utils.RepositoryResourceName, "repo-empty-conn-draining"),
		Type: repository.MongoDB,
		ConnParams: &repository.ConnParams{
			ConnDraining: &repository.ConnDraining{},
		},
		RepoNodes: repository.RepoNodes{
			{
				Host: "mongo-cluster.local",
				Port: 27017,
			},
		},
		MongoDBSettings: &repository.MongoDBSettings{
			ServerType: repository.Standalone,
		},
	}

	connDrainingConfig = repository.RepoInfo{
		Name: utils.AccTestName(utils.RepositoryResourceName, "repo-conn-draining"),
		Type: repository.MongoDB,
		ConnParams: &repository.ConnParams{
			ConnDraining: &repository.ConnDraining{
				Auto:     true,
				WaitTime: 20,
			},
		},
		RepoNodes: repository.RepoNodes{
			{
				Host: "mongo-cluster.local",
				Port: 27017,
			},
		},
		MongoDBSettings: &repository.MongoDBSettings{
			ServerType: repository.Standalone,
		},
	}

	mixedMultipleNodesConfig = repository.RepoInfo{
		Name: utils.AccTestName(utils.RepositoryResourceName, "repo-mixed-multi-node"),
		Type: repository.MongoDB,
		ConnParams: &repository.ConnParams{
			ConnDraining: &repository.ConnDraining{
				Auto:     true,
				WaitTime: 20,
			},
		},
		RepoNodes: repository.RepoNodes{
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
		MongoDBSettings: &repository.MongoDBSettings{
			ReplicaSetName: "some-replica-set",
			ServerType:     repository.ReplicaSet,
			Flavor:         "mongodb",
		},
	}

	allRepoNodesAreDynamic = repository.RepoInfo{
		Name: utils.AccTestName(utils.RepositoryResourceName, "repo-all-repo-nodes-are-dynamic"),
		Type: "mongodb",
		RepoNodes: repository.RepoNodes{
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
		MongoDBSettings: &repository.MongoDBSettings{
			ReplicaSetName: "myReplicaSet",
			ServerType:     "replicaset",
			SRVRecordName:  "mySRVRecord",
			Flavor:         "documentdb",
		},
	}

	withRedshiftSettings = repository.RepoInfo{
		Name: utils.AccTestName(utils.RepositoryResourceName, "repo-with-redshift-settings"),
		Type: "redshift",
		RepoNodes: repository.RepoNodes{
			{
				Host: "redshift.local",
				Port: 3333,
			},
		},
		RedshiftSettings: &repository.RedshiftSettings{
			ClusterIdentifier: "myCluster",
			AwsRegion:         "us-east-1",
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
	redshift := setupRepositoryTest(
		withRedshiftSettings, "with_redshift_settings")

	multiNode := setupRepositoryTest(
		mixedMultipleNodesConfig, "multi_node_test")

	// Must use name of the last resource created.
	importTest := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "cyral_repository.multi_node_test",
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			initial,
			update,
			connDrainingEmpty,
			connDraining,
			redshift,
			allDynamic,
			multiNode,
			importTest,
		},
	})
}

func setupRepositoryTest(repo repository.RepoInfo, resName string) resource.TestStep {
	return resource.TestStep{
		Config: repoAsConfig(repo, resName),
		Check:  repoCheckFuctions(repo, resName),
	}
}

func repoCheckFuctions(repo repository.RepoInfo, resName string) resource.TestCheckFunc {
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
		if repo.MongoDBSettings.ServerType == repository.ReplicaSet {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceFullName,
					"mongodb_settings.0.replica_set_name",
					repo.MongoDBSettings.ReplicaSetName),
			}...)
		}
		if repo.MongoDBSettings.Flavor != "" {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceFullName,
					"mongodb_settings.0.flavor",
					repo.MongoDBSettings.Flavor),
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

	if repo.RedshiftSettings != nil {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceFullName,
				"redshift_settings.0.cluster_identifier",
				repo.RedshiftSettings.ClusterIdentifier,
			),
			resource.TestCheckResourceAttr(resourceFullName,
				"redshift_settings.0.workgroup_name",
				repo.RedshiftSettings.WorkgroupName,
			),
			resource.TestCheckResourceAttr(resourceFullName,
				"redshift_settings.0.aws_region",
				repo.RedshiftSettings.AwsRegion,
			),
		}...)
	}

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func repoAsConfig(repo repository.RepoInfo, resName string) string {
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
		flavor := "null"
		if repo.MongoDBSettings.ReplicaSetName != "" {
			replicaSet = fmt.Sprintf(`"%s"`, repo.MongoDBSettings.ReplicaSetName)
		}
		if repo.MongoDBSettings.ServerType != "" {
			serverType = fmt.Sprintf(`"%s"`, repo.MongoDBSettings.ServerType)
		}
		if repo.MongoDBSettings.SRVRecordName != "" {
			srvRecordName = fmt.Sprintf(`"%s"`, repo.MongoDBSettings.SRVRecordName)
		}
		if repo.MongoDBSettings.Flavor != "" {
			flavor = fmt.Sprintf(`"%s"`, repo.MongoDBSettings.Flavor)
		}
		config += fmt.Sprintf(`
		mongodb_settings {
			replica_set_name = %s
			server_type = %s
			srv_record_name = %s
			flavor = %s
		}`,
			replicaSet,
			serverType,
			srvRecordName,
			flavor,
		)
	}

	if repo.RedshiftSettings != nil {
		clusterIdentifier := "null"
		workgroupName := "null"
		awsRegion := "null"

		if repo.RedshiftSettings.ClusterIdentifier != "" {
			clusterIdentifier = fmt.Sprintf(`"%s"`, repo.RedshiftSettings.ClusterIdentifier)
		}

		if repo.RedshiftSettings.WorkgroupName != "" {
			workgroupName = fmt.Sprintf(`"%s"`, repo.RedshiftSettings.WorkgroupName)
		}

		if repo.RedshiftSettings.AwsRegion != "" {
			awsRegion = fmt.Sprintf(`"%s"`, repo.RedshiftSettings.AwsRegion)
		}

		config += fmt.Sprintf(`
			redshift_settings {
				cluster_identifier = %s
				workgroup_name = %s
				aws_region = %s
			}`,
			clusterIdentifier,
			workgroupName,
			awsRegion,
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
