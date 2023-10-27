package internal_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	integrationDatadogResourceName = "integration-datadog"
)

var initialDatadogConfig internal.DatadogIntegration = internal.DatadogIntegration{
	Name:   utils.AccTestName(integrationDatadogResourceName, "datadog"),
	APIKey: "some-api-key",
}

var updatedDatadogConfig internal.DatadogIntegration = internal.DatadogIntegration{
	Name:   utils.AccTestName(integrationDatadogResourceName, "datadog-updated"),
	APIKey: "some-api-key-updated",
}

func TestAccDatadogIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupDatadogTest(initialDatadogConfig)
	testUpdateConfig, testUpdateFunc := setupDatadogTest(updatedDatadogConfig)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"cyral": func() (*schema.Provider, error) {
				return provider.Provider(), nil
			},
		},
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

func setupDatadogTest(integrationData internal.DatadogIntegration) (string, resource.TestCheckFunc) {
	configuration := formatDatadogIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration",
			"name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration",
			"api_key", integrationData.APIKey))

	return configuration, testFunction
}

func formatDatadogIntegrationDataIntoConfig(data internal.DatadogIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_datadog" "datadog_integration" {
		name = "%s"
		api_key = "%s"
	}`, data.Name, data.APIKey)
}
