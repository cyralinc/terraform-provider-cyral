package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialSumoLogicConfig SumoLogicIntegrationData = SumoLogicIntegrationData{
	Name:    "sumo-logic-test",
	Address: "sumologic.local",
}

var updatedSumoLogicConfig SumoLogicIntegrationData = SumoLogicIntegrationData{
	Name:    "sumo-logic-test-update",
	Address: "sumologic-updated.local",
}

// This is loosely based on this example:
// https://github.com/hashicorp/terraform-provider-vault/blob/master/vault/resource_azure_secret_backend_role_test.go
func TestAccSumoLogicIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupSumoLogicTest(initialSumoLogicConfig)
	testUpdateConfig, testUpdateFunc := setupSumoLogicTest(updatedSumoLogicConfig)

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

func setupSumoLogicTest(integrationData SumoLogicIntegrationData) (string, resource.TestCheckFunc) {
	configuration := formatSumoLogicIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sumo_logic.sumo_logic_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_sumo_logic.sumo_logic_integration", "address", integrationData.Address),
	)

	return configuration, testFunction
}

func formatSumoLogicIntegrationDataIntoConfig(data SumoLogicIntegrationData) string {
	return fmt.Sprintf(`
	resource "cyral_integration_sumo_logic" "sumo_logic_integration" {
		name = "%s"
		address = "%s"
	}`, data.Name, data.Address)
}
