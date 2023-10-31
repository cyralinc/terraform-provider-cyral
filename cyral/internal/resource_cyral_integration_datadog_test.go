package internal_test

import (
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

func setupDatadogTest(d internal.DatadogIntegration) (string, resource.TestCheckFunc) {
	configuration := utils.FormatDatadogIntegrationDataIntoConfig(d.Name, d.APIKey)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration",
			"name", d.Name),
		resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration",
			"api_key", d.APIKey))

	return configuration, testFunction
}
