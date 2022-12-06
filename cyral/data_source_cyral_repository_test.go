package cyral

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryDataSourceName = "data-repository"
)

func repositoryDataSourceTestRepos() []RepoInfo {
	return []RepoInfo{
		{
			Name:   accTestName(repositoryDataSourceName, "sqlserver-1"),
			Type:   "sqlserver",
			Labels: []string{"rds", "us-east-2"},
			RepoNodes: []*RepoNode{
				{
					Host: "sql.local",
					Port: 3333,
				},
			},
		},
		{
			Name:   accTestName(repositoryDataSourceName, "mongodb-1"),
			Type:   "mongodb",
			Labels: []string{"rds", "us-east-1"},
			RepoNodes: []*RepoNode{
				{
					Host: "mongo.local",
					Port: 27017,
				},
			},
		},
	}
}

func TestAccRepositoryDataSource(t *testing.T) {
	testRepos := repositoryDataSourceTestRepos()
	testConfigNameFilter, testFuncNameFilter := testRepositoryDataSource(
		testRepos, fmt.Sprintf("^%s$", testRepos[0].Name), "")
	testConfigTypeFilter, testFuncTypeFilter := testRepositoryDataSource(
		testRepos, "", "mongodb")
	testConfigNameTypeFilter, testFuncNameTypeFilter := testRepositoryDataSource(
		repositoryDataSourceTestRepos(), fmt.Sprintf("^%s$", testRepos[1].Name), "mongodb")

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigNameFilter,
				Check:  testFuncNameFilter,
			},
			{
				Config: testConfigTypeFilter,
				Check:  testFuncTypeFilter,
			},
			{
				Config: testConfigNameTypeFilter,
				Check:  testFuncNameTypeFilter,
			},
		},
	})
}

func testRepositoryDataSource(repoDatas []RepoInfo, nameFilter, typeFilter string) (
	string, resource.TestCheckFunc,
) {
	return testRepositoryDataSourceConfig(repoDatas, nameFilter, typeFilter),
		testRepositoryDataSourceChecks(repoDatas, nameFilter, typeFilter)
}

func testRepositoryDataSourceConfig(repoDatas []RepoInfo, nameFilter, typeFilter string) string {
	var config string
	var dependsOn []string
	for _, repoData := range repoDatas {
		config += repoAsConfig(repoData, repoData.Name)
		dependsOn = append(dependsOn, fmt.Sprintf("cyral_repository.%s", repoData.Name))
	}
	config += repositoryDataSourceConfig(nameFilter, typeFilter, dependsOn)

	return config
}

func testRepositoryDataSourceChecks(repoDatas []RepoInfo, nameFilter, typeFilter string) resource.TestCheckFunc {
	dataSourceFullName := "data.cyral_repository.test_repository"

	if nameFilter == "" {
		return resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(dataSourceFullName,
				"repository_list.#",
				notZeroRegex(),
			),
			dsourceCheckTypeFilter(
				dataSourceFullName,
				"repository_list.%d.type",
				typeFilter,
			),
		)
	}

	var checkFuncs []resource.TestCheckFunc
	filteredRepoDatas := filterRepoDatas(repoDatas, nameFilter, typeFilter)
	if len(filteredRepoDatas) == 1 {
		repo := filteredRepoDatas[0]
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceFullName,
				"repository_list.0.name", repo.Name),
			resource.TestCheckResourceAttr(dataSourceFullName,
				"repository_list.0.type", repo.Type),
			resource.TestCheckResourceAttr(dataSourceFullName,
				"repository_list.0.labels.#", fmt.Sprintf("%d", len(repo.Labels)),
			),
		}...)

		for i, label := range repo.Labels {
			checkFuncs = append(checkFuncs,
				resource.TestCheckResourceAttr(dataSourceFullName,
					fmt.Sprintf("repository_list.0.labels.%d", i), label))
		}

		for i, node := range repo.RepoNodes {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(dataSourceFullName,
					fmt.Sprintf("repository_list.0.repo_node.%d.name", i), node.Name),
				resource.TestCheckResourceAttr(dataSourceFullName,
					fmt.Sprintf("repository_list.0.repo_node.%d.host", i), node.Host),
				resource.TestCheckResourceAttr(dataSourceFullName,
					fmt.Sprintf("repository_list.0.repo_node.%d.port", i),
					strconv.Itoa(int(node.Port))),
				resource.TestCheckResourceAttr(dataSourceFullName,
					fmt.Sprintf("repository_list.0.repo_node.%d.dynamic", i),
					strconv.FormatBool(node.Dynamic)),
			}...)
		}

		if repo.ConnParams != nil {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(dataSourceFullName,
					"repository_list.0.connection_draining.0.auto",
					strconv.FormatBool(repo.ConnParams.ConnDraining.Auto)),

				resource.TestCheckResourceAttr(dataSourceFullName,
					"repository_list.0.connection_draining.0.wait_time",
					strconv.Itoa(int(repo.ConnParams.ConnDraining.WaitTime))),
			}...)
		}

		if repo.MongoDBSettings != nil {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(dataSourceFullName,
					"repository_list.0.mongodb_settings.0.replica_set_name",
					repo.MongoDBSettings.ReplicaSetName),

				resource.TestCheckResourceAttr(dataSourceFullName,
					"repository_list.0.connection_draining.0.server_type",
					repo.MongoDBSettings.ServerType),
			}...)
		}
	}

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return testFunction
}

func filterRepoDatas(repoDatas []RepoInfo, nameFilter, typeFilter string) []RepoInfo {
	var filteredRepoDatas []RepoInfo
	for _, repoData := range repoDatas {
		if (nameFilter == "" || repoData.Name == nameFilter) &&
			(typeFilter == "" || repoData.Type == typeFilter) {
			filteredRepoDatas = append(filteredRepoDatas, repoData)
		}
	}
	return filteredRepoDatas
}

func repositoryDataSourceConfig(nameFilter, typeFilter string, dependsOn []string) string {
	return fmt.Sprintf(`
	data "cyral_repository" "test_repository" {
		depends_on = %s
		name = "%s"
		type = "%s"
	}`, listToStr(dependsOn), nameFilter, typeFilter)
}
