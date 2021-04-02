package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialSlackAlertsConfig SlackAlertsIntegration = SlackAlertsIntegration{
	Name: "tf-test-slack-alerts",
	URL:  "https://slack.local",
}

var updatedSlackAlertsConfig SlackAlertsIntegration = SlackAlertsIntegration{
	Name: "tf-test-update-slack-alerts",
	URL:  "https://slack-updated.local",
}

func TestAccSlackAlertsIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupSlackAlertTest(initialSlackAlertsConfig)
	testUpdateConfig, testUpdateFunc := setupSlackAlertTest(initialSlackAlertsConfig)

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

func setupSlackAlertTest(integrationData SlackAlertsIntegration) (string, resource.TestCheckFunc) {
	configuration := formatSlackAlertsIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_slack_alerts.test_slack_alerts",
			"name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_slack_alerts.test_slack_alerts",
			"url", integrationData.URL),
	)

	return configuration, testFunction
}

func formatSlackAlertsIntegrationDataIntoConfig(data SlackAlertsIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_slack_alerts" "test_slack_alerts" {
		name = "%s"
		url  = "%s"
	}`, data.Name, data.URL)
}
