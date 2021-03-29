package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialDatadogConfig DatadogIntegrationData = DatadogIntegrationData{
	Name:   "unitTest-Datadog",
	APIKey: "some-api-key",
}

var updatedDatadogConfig DatadogIntegrationData = DatadogIntegrationData{
	Name:   "unitTest-Datadog-updated",
	APIKey: "some-api-key-updated",
}

// This is loosely based on this example:
// https://github.com/hashicorp/terraform-provider-vault/blob/master/vault/resource_azure_secret_backend_role_test.go
func TestAccDatadogIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupDatadogTest(initialDatadogConfig)
	testUpdateConfig, testUpdateFunc := setupDatadogTest(updatedDatadogConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
		},
	})
}

func setupDatadogTest(integrationData DatadogIntegrationData) (string, resource.TestCheckFunc) {
	configuration := formatDatadogIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration", "api_key", integrationData.APIKey))

	return configuration, testFunction
}

func formatDatadogIntegrationDataIntoConfig(data DatadogIntegrationData) string {
	return fmt.Sprintf(`
resource "cyral_integration_datadog" "datadog_integration" {
    name = "%s"
    api_key = "%s"
}`, data.Name, data.APIKey)
}
