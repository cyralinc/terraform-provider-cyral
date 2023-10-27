package internal_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationSlackAlertsResourceName = "integration-slack-alerts"
)

var initialSlackAlertsConfig internal.SlackAlertsIntegration = internal.SlackAlertsIntegration{
	Name: utils.AccTestName(integrationSlackAlertsResourceName, "slack-alerts"),
	URL:  "https://slack.local",
}

var updatedSlackAlertsConfig internal.SlackAlertsIntegration = internal.SlackAlertsIntegration{
	Name: utils.AccTestName(integrationSlackAlertsResourceName, "slack-alerts-updated"),
	URL:  "https://slack-updated.local",
}

func TestAccSlackAlertsIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupSlackAlertTest(initialSlackAlertsConfig)
	testUpdateConfig, testUpdateFunc := setupSlackAlertTest(initialSlackAlertsConfig)

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
				ResourceName:      "cyral_integration_slack_alerts.test_slack_alerts",
			},
		},
	})
}

func setupSlackAlertTest(integrationData internal.SlackAlertsIntegration) (string, resource.TestCheckFunc) {
	configuration := formatSlackAlertsIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_slack_alerts.test_slack_alerts",
			"name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_slack_alerts.test_slack_alerts",
			"url", integrationData.URL),
	)

	return configuration, testFunction
}

func formatSlackAlertsIntegrationDataIntoConfig(data internal.SlackAlertsIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_slack_alerts" "test_slack_alerts" {
		name = "%s"
		url  = "%s"
	}`, data.Name, data.URL)
}
