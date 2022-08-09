package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func repositoryDataSourceTestRepos() []RepoData {
	return []RepoData{
		{
			Name:     "tfprov-test-repository-dsource-sqlserver-1",
			Host:     "localhost",
			Port:     1433,
			RepoType: "sqlserver",
			Labels:   []string{"rds", "us-east-2"},
		},
		{
			Name:                "tfprov-test-repository-dsource-mongodb-1",
			Host:                "localhost",
			Port:                27017,
			RepoType:            "mongodb",
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
	testConfigNameFilter, testFuncNameFilter := testRepositoryDataSource(
		repositoryDataSourceTestRepos(), "^tfprov-test-repository-dsource-sqlserver-1$", "")
	testConfigTypeFilter, testFuncTypeFilter := testRepositoryDataSource(
		repositoryDataSourceTestRepos(), "", "mongodb")
	testConfigNameTypeFilter, testFuncNameTypeFilter := testRepositoryDataSource(
		repositoryDataSourceTestRepos(), "^tfprov-test-repository-dsource-mongodb-1$", "mongodb")

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

func testRepositoryDataSource(repoDatas []RepoData, nameFilter, typeFilter string) (
	string, resource.TestCheckFunc,
) {
	return testRepositoryDataSourceConfig(repoDatas, nameFilter, typeFilter),
		testRepositoryDataSourceChecks(repoDatas, nameFilter, typeFilter)
}

func testRepositoryDataSourceConfig(repoDatas []RepoData, nameFilter, typeFilter string) string {
	var config string
	var dependsOn []string
	for _, repoData := range repoDatas {
		config += formatRepoDataIntoConfig(repoData)
		dependsOn = append(dependsOn, repositoryConfigResourceFullName(repoData.Name))
	}
	config += repositoryDataSourceConfig(nameFilter, typeFilter, dependsOn)

	return config
}

func testRepositoryDataSourceChecks(repoDatas []RepoData, nameFilter, typeFilter string) resource.TestCheckFunc {
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
				"repository_list.0.type", repoData.RepoType),
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

func filterRepoDatas(repoDatas []RepoData, nameFilter, typeFilter string) []RepoData {
	var filteredRepoDatas []RepoData
	for _, repoData := range repoDatas {
		if (nameFilter == "" || repoData.Name == nameFilter) &&
			(typeFilter == "" || repoData.RepoType == typeFilter) {
			filteredRepoDatas = append(filteredRepoDatas, repoData)
		}
	}
	return filteredRepoDatas
}

func repositoryDataSourceConfig(nameFilter, typeFilter string, dependsOn []string) string {
	return fmt.Sprintf(`
	data "cyral_repository" "test_repository" {
		depends_on = [%s]
		name = "%s"
		type = "%s"
	}`, listToStr(dependsOn), nameFilter, typeFilter)
}
