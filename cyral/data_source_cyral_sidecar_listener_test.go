package cyral

const (
	listenerDataSourceName = "data-repository"
)

func listenerDataSourceTestRepos() []SidecarListener {
	return []SidecarListener{
		{
			RepoTypes: []string{"sqlserver"},
			NetworkAddress: &NetworkAddress{
				Port: 3333,
			},
		},
		{
			RepoTypes: []string{"mongodb"},
			NetworkAddress: &NetworkAddress{
				Port: 27017,
			},
		},
	}
}

// func TestAccSidecarListenerDataSource(t *testing.T) {
// 	testRepos := listenerDataSourceTestRepos()
// 	testConfigNameFilter, testFuncNameFilter := testListenerDataSource(
// 		testRepos, fmt.Sprintf("^%s$", testRepos[0].Name), "")
// 	testConfigTypeFilter, testFuncTypeFilter := testListenerDataSource(
// 		testRepos, "", "mongodb")
// 	testConfigNameTypeFilter, testFuncNameTypeFilter := testListenerDataSource(
// 		listenerDataSourceTestRepos(), fmt.Sprintf("^%s$", testRepos[1].Name), "mongodb")

// 	resource.ParallelTest(t, resource.TestCase{
// 		ProviderFactories: providerFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testConfigNameFilter,
// 				Check:  testFuncNameFilter,
// 			},
// 			{
// 				Config: testConfigTypeFilter,
// 				Check:  testFuncTypeFilter,
// 			},
// 			{
// 				Config: testConfigNameTypeFilter,
// 				Check:  testFuncNameTypeFilter,
// 			},
// 		},
// 	})
// }

// func listenerDataSourceConfig(nameFilter, typeFilter string, dependsOn []string) string {
// 	return fmt.Sprintf(`
// 	data "cyral_sidecar_listener" "test_listener" {
// 		depends_on = %s
// 		sidecar_id = "%s"
// 		repo_types = ["%s"]
// 	}`, listToStr(dependsOn), nameFilter, typeFilter)
// }
