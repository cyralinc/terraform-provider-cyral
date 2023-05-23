package cyral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationLogsDataSourceName = "integration-log-datasource"
)

func TestAccLogsIntegrationDataSource(t *testing.T) {
	// Using vars from resource test
	testConfigElk, testFuncElk := setupLogsTest(initialLogsConfigElk)
	testUpdateConfigElk, testUpdateFuncElk := setupLogsTest(updatedLogsConfigElk)

	testConfigCloudWatch, testFuncCloudWatch := setupLogsTest(initialLogsConfigCloudWatch)
	testUpdateConfigCloudWatch, testUpdateFuncCloudWatch := setupLogsTest(updatedLogsConfigCloudWatch)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigElk,
				Check:  testFuncElk,
			},
			{
				Config: testUpdateConfigElk,
				Check:  testUpdateFuncElk,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_integration_logs.logs_integration",
			},
			{
				Config: testAccIntegrationLogsDataSourceConfigElk(),
				Check:  testAccIntegrationLogsDataSourceCheckElk(),
			},
			{
				Config: testConfigCloudWatch,
				Check:  testFuncCloudWatch,
			},
			{
				Config: testUpdateConfigCloudWatch,
				Check:  testUpdateFuncCloudWatch,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_integration_logs.logs_integration",
			},
			{
				Config: testAccIntegrationLogsDataSourceConfigCloudWatch(),
				Check:  testAccIntegrationLogsDataSourceCheckCloudWatch(),
			},
		},
	})
}

func testAccIntegrationLogsDataSourceConfigElk() string {
	return `
	data "cyral_integration_logs" "list_integrations" {
		type = "ELK"
	}
	`
}

func testAccIntegrationLogsDataSourceCheckElk() resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	pathResource := "data.cyral_integration_logs.list_integrations"

	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.name"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.receive_audit_logs"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.config_scheme.0.elk.0.es_url"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.config_scheme.0.elk.0.kibana_url"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.config_scheme.0.elk.0.es_credentials.0.username"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.config_scheme.0.elk.0.es_credentials.0.password"),
	}...)

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return testFunction
}

func testAccIntegrationLogsDataSourceConfigCloudWatch() string {
	return `
	data "cyral_integration_logs" "list_integrations2" {
		type = "CLOUDWATCH"
	}
	`
}

func testAccIntegrationLogsDataSourceCheckCloudWatch() resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	pathResource := "data.cyral_integration_logs.list_integrations2"

	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.name"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.receive_audit_logs"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.config_scheme.0.cloud_watch.0.region"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.config_scheme.0.cloud_watch.0.group"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integration_list.0.config_scheme.0.cloud_watch.0.stream"),
	}...)

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return testFunction
}
