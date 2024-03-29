package pagerduty_test

import (
	"fmt"
	"testing"

	ce "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationPagerDutyResourceName = "integration-pager-duty"
)

func initialPagerDutyIntegrationConfig() *ce.IntegrationConfExtension {
	integration := ce.NewIntegrationConfExtension(ce.PagerDutyTemplateType)
	integration.Name = utils.AccTestName(integrationPagerDutyResourceName, "pager-duty")
	integration.Parameters = "unitTest-parameters"
	return integration
}

func updatedPagerDutyIntegrationConfig() *ce.IntegrationConfExtension {
	integration := ce.NewIntegrationConfExtension(ce.PagerDutyTemplateType)
	integration.Name = utils.AccTestName(integrationPagerDutyResourceName, "pager-duty-updated")
	integration.Parameters = "unitTest-parameters-updated"
	return integration
}

func TestAccPagerDutyIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupPagerDutyIntegrationTest(initialPagerDutyIntegrationConfig())
	testUpdateConfig, testUpdateFunc := setupPagerDutyIntegrationTest(updatedPagerDutyIntegrationConfig())

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
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

func setupPagerDutyIntegrationTest(integrationData *ce.IntegrationConfExtension) (string, resource.TestCheckFunc) {
	configuration := formatPagerDutyIntegrationIntoConfig(
		integrationData.Name, integrationData.Parameters)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_pager_duty.pager_duty_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_pager_duty.pager_duty_integration", "api_token", integrationData.Parameters),
	)

	return configuration, testFunction
}

func formatPagerDutyIntegrationIntoConfig(name, apiToken string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_pager_duty" "pager_duty_integration" {
		name = "%s"
		api_token = "%s"
	}`, name, apiToken)
}
