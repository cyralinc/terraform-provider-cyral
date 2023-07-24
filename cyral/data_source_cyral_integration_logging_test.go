package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationLogsDataSourceName = "integration-log-datasource"
)

func testIntegrationLoggingDataSource(t *testing.T, resName string, typeFilter string) (
	string, resource.TestCheckFunc,
) {
	return testIntegrationLoggingDataSourceConfig(resName, typeFilter),
		testIntegrationLoggingDataSourceChecks(t, resName, typeFilter)
}

func testIntegrationLoggingDataSourceConfig(resName, typeFilter string) string {
	var config string
	config += testIntegrationLoggingDataSourceConfigDependencies(resName)
	config += integrationLogsDataSourceConfig(
		resName,
		typeFilter,
		[]string{
			fmt.Sprintf("cyral_integration_logging.%s_1", resName),
			fmt.Sprintf("cyral_integration_logging.%s_2", resName),
		})
	return config
}

// Setup two integrations that are retrieved and checked by the data source.
func testIntegrationLoggingDataSourceConfigDependencies(resName string) string {
	resName1 := resName + "_1"
	resName2 := resName + "_2"

	resource1, _ := formatLogsIntegrationDataIntoConfig(LoggingIntegration{
		Name:             resName1,
		ReceiveAuditLogs: true,
		LoggingIntegrationConfig: LoggingIntegrationConfig{
			CloudWatch: &CloudWatchConfig{
				Region: "us-east-2",
				Group:  "group2",
				Stream: "abcd",
			},
		},
	}, resName1)

	resource2, _ := formatLogsIntegrationDataIntoConfig(LoggingIntegration{
		Name:             resName2,
		ReceiveAuditLogs: true,
		LoggingIntegrationConfig: LoggingIntegrationConfig{
			Datadog: &DataDogConfig{
				ApiKey: "API_KEY_A",
			},
		},
	}, resName2)

	var config string
	config += resource1
	config += resource2
	return config
}

func TestAccLoggingIntegrationDataSource(t *testing.T) {
	testConfig1, testFunc1 := testIntegrationLoggingDataSource(t,
		accTestName(integrationLogsDataSourceName, "test1"), "CLOUDWATCH")
	testConfig2, testFunc2 := testIntegrationLoggingDataSource(t,
		accTestName(integrationLogsDataSourceName, "test2"), "DATADOG")

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig1,
				Check:  testFunc1,
			},
			{
				Config: testConfig2,
				Check:  testFunc2,
			},
		},
	})
}

func testIntegrationLoggingDataSourceChecks(t *testing.T, resName, typeFilter string) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	pathResource := fmt.Sprintf("data.cyral_integration_logging.%s", resName)

	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(pathResource,
			"integrations.0.name"),
		resource.TestCheckResourceAttrSet(pathResource,
			"integrations.0.receive_audit_logs"),
	}...)

	if typeFilter == "DATADOG" {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttrSet(pathResource,
				"integrations.0.datadog.0.api_key"),
		}...)
	}

	if typeFilter == "CLOUDWATCH" {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttrSet(pathResource,
				"integrations.0.cloudwatch.0.region"),
			resource.TestCheckResourceAttrSet(pathResource,
				"integrations.0.cloudwatch.0.group"),
			resource.TestCheckResourceAttrSet(pathResource,
				"integrations.0.cloudwatch.0.stream"),
		}...)
	}

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return testFunction
}

func integrationLogsDataSourceConfig(resName string, typeFilter string, dependsOn []string) string {
	return fmt.Sprintf(`
	data "cyral_integration_logging" "%s" {
		depends_on = %s
		type = "%s"
	}`, resName, listToStr(dependsOn), typeFilter)
}
