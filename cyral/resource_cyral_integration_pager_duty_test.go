package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationPagerDutyResourceName = "integration-pager-duty"
)

var initialPagerDutyIntegrationConfig PagerDutyIntegration = PagerDutyIntegration{
	Name:       accTestName(integrationPagerDutyResourceName, "pager-duty"),
	Parameters: "unitTest-parameters",
}

var updatedPagerDutyIntegrationConfig PagerDutyIntegration = PagerDutyIntegration{
	Name:       accTestName(integrationPagerDutyResourceName, "pager-duty-updated"),
	Parameters: "unitTest-parameters-updated",
}

func TestAccPagerDutyIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupPagerDutyIntegrationTest(initialPagerDutyIntegrationConfig)
	testUpdateConfig, testUpdateFunc := setupPagerDutyIntegrationTest(updatedPagerDutyIntegrationConfig)

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
				ResourceName:      "cyral_integration_pager_duty.pager_duty_integration",
			},
		},
	})
}

func setupPagerDutyIntegrationTest(integrationData PagerDutyIntegration) (string, resource.TestCheckFunc) {
	configuration := formatPagerDutyIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(

		resource.TestCheckResourceAttr("cyral_integration_pager_duty.pager_duty_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_pager_duty.pager_duty_integration", "api_token", integrationData.Parameters),
	)

	return configuration, testFunction
}

func formatPagerDutyIntegrationDataIntoConfig(data PagerDutyIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_pager_duty" "pager_duty_integration" {
		name = "%s"
		api_token = "%s"
	}`, data.Name, data.Parameters)
}
