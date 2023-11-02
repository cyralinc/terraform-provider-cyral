package slack_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/slack"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationSlackAlertsResourceName = "integration-slack-alerts"
)

var initialSlackAlertsConfig slack.SlackAlertsIntegration = slack.SlackAlertsIntegration{
	Name: utils.AccTestName(integrationSlackAlertsResourceName, "slack-alerts"),
	URL:  "https://slack.local",
}

var updatedSlackAlertsConfig slack.SlackAlertsIntegration = slack.SlackAlertsIntegration{
	Name: utils.AccTestName(integrationSlackAlertsResourceName, "slack-alerts-updated"),
	URL:  "https://slack-updated.local",
}

func TestAccSlackAlertsIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupSlackAlertTest(initialSlackAlertsConfig)
	testUpdateConfig, testUpdateFunc := setupSlackAlertTest(initialSlackAlertsConfig)

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
				ResourceName:      "cyral_integration_slack_alerts.test_slack_alerts",
			},
		},
	})
}

func setupSlackAlertTest(integrationData slack.SlackAlertsIntegration) (string, resource.TestCheckFunc) {
	configuration := formatSlackAlertsIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_slack_alerts.test_slack_alerts",
			"name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_slack_alerts.test_slack_alerts",
			"url", integrationData.URL),
	)

	return configuration, testFunction
}

func formatSlackAlertsIntegrationDataIntoConfig(data slack.SlackAlertsIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_slack_alerts" "test_slack_alerts" {
		name = "%s"
		url  = "%s"
	}`, data.Name, data.URL)
}
