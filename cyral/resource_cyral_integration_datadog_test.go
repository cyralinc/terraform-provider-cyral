package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialDatadogConfig DatadogIntegration = DatadogIntegration{
	Name:   "unitTest-Datadog",
	APIKey: "some-api-key",
}

var updatedDatadogConfig DatadogIntegration = DatadogIntegration{
	Name:   "unitTest-Datadog-updated",
	APIKey: "some-api-key-updated",
}

func TestAccDatadogIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupDatadogTest(initialDatadogConfig)
	testUpdateConfig, testUpdateFunc := setupDatadogTest(updatedDatadogConfig)

	resource.Test(t, resource.TestCase{
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
