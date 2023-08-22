package cyral

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"golang.org/x/exp/slices"
)

const (
	listenerDataSourceName = "data-repository"
)

func listenerDataSourceTestRepos() []SidecarListener {
	return []SidecarListener{
		{
			RepoTypes: []string{"mysql"},
			NetworkAddress: &NetworkAddress{
				Port: 3306,
			},
		},
		{
			RepoTypes: []string{"mongodb"},
			NetworkAddress: &NetworkAddress{
				Port: 27017,
			},
		},
		{
			RepoTypes: []string{"mongodb"},
			NetworkAddress: &NetworkAddress{
				Port: 27018,
			},
		},
	}
}

func TestAccSidecarListenerDataSource(t *testing.T) {
	testListeners := listenerDataSourceTestRepos()

	testConfigTypeFilter, testFuncTypeFilter := testListenerDataSource(
		testListeners, testListeners[0].RepoTypes[0], 0)

	testConfigPortFilter, testFuncPortFilter := testListenerDataSource(
		testListeners, "", testListeners[1].NetworkAddress.Port)

	testConfigTypePortFilter, testFuncTypePortFilter := testListenerDataSource(
		testListeners, testListeners[2].RepoTypes[0], testListeners[0].NetworkAddress.Port)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigTypeFilter,
				Check:  testFuncTypeFilter,
			},
			{
				Config: testConfigPortFilter,
				Check:  testFuncPortFilter,
			},
			{
				Config: testConfigTypePortFilter,
				Check:  testFuncTypePortFilter,
			},
		},
	})
}

func testListenerDataSource(listeners []SidecarListener, repoTypeFilter string, portFilter int) (
	string, resource.TestCheckFunc,
) {
	return testListenerDataSourceConfig(listeners, repoTypeFilter, portFilter),
		testListenerDataSourceChecks(listeners, repoTypeFilter, portFilter)
}

func testListenerDataSourceConfig(listeners []SidecarListener, repoTypeFilter string, portFilter int) string {
	var config string
	var dependsOn []string
	for _, listener := range listeners {
		resourceName := fmt.Sprintf("%s_%d", listener.RepoTypes[0], listener.NetworkAddress.Port)
		config += setupSidecarListenerConfig(resourceName, listener)
		dependsOn = append(dependsOn, fmt.Sprintf("cyral_sidecar_listener.%s", resourceName))
	}
	sidecarConfig := utils.FormatBasicSidecarIntoConfig(
		utils.BasicSidecarResName,
		utils.AccTestName("ds-sidecar-listener", "sidecar"),
		"docker", "",
	)
	config += sidecarConfig + listenerDataSourceConfig(repoTypeFilter, portFilter, dependsOn)

	return config
}

func testListenerDataSourceChecks(listeners []SidecarListener, repoTypeFilter string, portFilter int) resource.TestCheckFunc {
	dataSourceFullName := "data.cyral_sidecar_listener.test_listener"

	if repoTypeFilter == "" && portFilter == 0 {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(dataSourceFullName,
				"listener_list.#",
				fmt.Sprintf("%d", len(listeners)),
			),
		)
	}

	var checkFuncs []resource.TestCheckFunc
	filteredListeners := filterListenerData(listeners, repoTypeFilter, portFilter)
	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dataSourceFullName,
			"listener_list.#", fmt.Sprintf("%d", len(filteredListeners))),
	}...)

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func filterListenerData(listeners []SidecarListener, repoTypeFilter string, portFilter int) []SidecarListener {
	var filteredListeners []SidecarListener
	for _, l := range listeners {
		if (repoTypeFilter == "" || slices.Contains(l.RepoTypes, repoTypeFilter)) &&
			(portFilter == 0 || l.NetworkAddress.Port == portFilter) {
			filteredListeners = append(filteredListeners, l)
		}
	}
	return filteredListeners
}

func listenerDataSourceConfig(repoTypeFilter string, portFilter int, dependsOn []string) string {
	if repoTypeFilter != "" && portFilter > 0 {
		return fmt.Sprintf(`
			data "cyral_sidecar_listener" "test_listener" {
				depends_on = %s
				sidecar_id = %s
				repo_type = "%s"
				port = %d
			}`, utils.ListToStr(dependsOn), utils.BasicSidecarID, repoTypeFilter, portFilter)
	} else if repoTypeFilter != "" {
		return fmt.Sprintf(`
			data "cyral_sidecar_listener" "test_listener" {
				depends_on = %s
				sidecar_id = %s
				repo_type = "%s"
			}`, utils.ListToStr(dependsOn), utils.BasicSidecarID, repoTypeFilter)
	} else {
		return fmt.Sprintf(`
			data "cyral_sidecar_listener" "test_listener" {
				depends_on = %s
				sidecar_id = %s
				port = %d
			}`, utils.ListToStr(dependsOn), utils.BasicSidecarID, portFilter)
	}
}
