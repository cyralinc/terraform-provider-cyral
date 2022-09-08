package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationDatadogResourceName = "integration-datadog"
)

var initialDatadogConfig DatadogIntegration = DatadogIntegration{
	Name:   accTestName(integrationDatadogResourceName, "datadog"),
	APIKey: "some-api-key",
}

var updatedDatadogConfig DatadogIntegration = DatadogIntegration{
	Name:   accTestName(integrationDatadogResourceName, "datadog-updated"),
	APIKey: "some-api-key-updated",
}

func TestAccDatadogIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupDatadogTest(initialDatadogConfig)
	testUpdateConfig, testUpdateFunc := setupDatadogTest(updatedDatadogConfig)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_integration_datadog.datadog_integration",
			},
		},
	})
}

func setupDatadogTest(integrationData DatadogIntegration) (string, resource.TestCheckFunc) {
	configuration := formatDatadogIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration",
			"name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration",
			"api_key", integrationData.APIKey))

	return configuration, testFunction
}

func formatDatadogIntegrationDataIntoConfig(data DatadogIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_datadog" "datadog_integration" {
		name = "%s"
		api_key = "%s"
	}`, data.Name, data.APIKey)
}
