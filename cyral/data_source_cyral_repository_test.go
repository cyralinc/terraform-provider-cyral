package cyral

import (
	"fmt"
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
			Host:   "localhost",
			Port:   1433,
			Type:   "sqlserver",
			Labels: []string{"rds", "us-east-2"},
		},
		{
			Name:                accTestName(repositoryDataSourceName, "mongodb-1"),
			Host:                "localhost",
			Port:                27017,
			Type:                "mongodb",
			Labels:              []string{"rds", "us-east-1"},
			MaxAllowedListeners: 2,
			Properties: &RepositoryProperties{
				MongoDBReplicaSetName: "replica-set-1",
				MongoDBServerType:     mongodbReplicaSetServerType,
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
		config += formatRepoDataIntoConfig(repoData, repoData.Name)
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
		repoData := filteredRepoDatas[0]
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceFullName,
				"repository_list.0.name", repoData.Name),
			resource.TestCheckResourceAttr(dataSourceFullName,
				"repository_list.0.type", repoData.Type),
			resource.TestCheckResourceAttr(dataSourceFullName,
				"repository_list.0.host", repoData.Host),
			resource.TestCheckResourceAttr(dataSourceFullName,
				"repository_list.0.port", fmt.Sprintf("%d", repoData.Port)),
			resource.TestCheckResourceAttr(dataSourceFullName,
				"repository_list.0.labels.#", fmt.Sprintf("%d", len(repoData.Labels)),
			),
		}...)

		if repoData.IsReplicaSet() {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(dataSourceFullName,
					"repository_list.0.properties.0.mongodb_replica_set.0.max_nodes",
					fmt.Sprintf("%d", repoData.MaxAllowedListeners)),
				resource.TestCheckResourceAttr(dataSourceFullName,
					"repository_list.0.properties.0.mongodb_replica_set.0.replica_set_id",
					repoData.Properties.MongoDBReplicaSetName),
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
