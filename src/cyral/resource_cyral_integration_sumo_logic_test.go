package cyral

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationSumoLogicResourceName = "integration-sumo-logic"
)

var initialSumoLogicConfig SumoLogicIntegration = SumoLogicIntegration{
	Name:    utils.AccTestName(integrationSumoLogicResourceName, "sumo-logic"),
	Address: "https://sumologic.local/initial",
}

var updatedSumoLogicConfig SumoLogicIntegration = SumoLogicIntegration{
	Name:    utils.AccTestName(integrationSumoLogicResourceName, "sumo-logic-updated"),
	Address: "https://sumologic.local/updated",
}

func TestAccSumoLogicIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupSumoLogicTest(initialSumoLogicConfig)
	testUpdateConfig, testUpdateFunc := setupSumoLogicTest(updatedSumoLogicConfig)

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
				ResourceName:      "cyral_integration_sumo_logic.sumo_logic_integration",
			},
		},
	})
}

func setupSumoLogicTest(integrationData SumoLogicIntegration) (string, resource.TestCheckFunc) {
	configuration := formatSumoLogicIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sumo_logic.sumo_logic_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_sumo_logic.sumo_logic_integration", "address", integrationData.Address),
	)

	return configuration, testFunction
}

func formatSumoLogicIntegrationDataIntoConfig(data SumoLogicIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_sumo_logic" "sumo_logic_integration" {
		name = "%s"
		address = "%s"
	}`, data.Name, data.Address)
}
