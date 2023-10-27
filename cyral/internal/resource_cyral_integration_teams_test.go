package internal_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationTeamsResourceName = "integrations-teams"
)

var initialTeamsConfig internal.MsTeamsIntegration = internal.MsTeamsIntegration{
	Name: utils.AccTestName(integrationTeamsResourceName, "msteams-alerts"),
	URL:  "https://msteams.local",
}

var updatedTeamsConfig internal.MsTeamsIntegration = internal.MsTeamsIntegration{
	Name: utils.AccTestName(integrationTeamsResourceName, "msteams-alerts"),
	URL:  "https://msteams-updated.local",
}

func TestAccMsTeamsIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupTeamsTest(initialTeamsConfig)
	testUpdateConfig, testUpdateFunc := setupTeamsTest(initialTeamsConfig)

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
				ResourceName:      "cyral_integration_microsoft_teams.test_microsoft_teams",
			},
		},
	})
}

func setupTeamsTest(integrationData internal.MsTeamsIntegration) (string, resource.TestCheckFunc) {
	configuration := formatMsTeamsIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_microsoft_teams.test_microsoft_teams",
			"name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_microsoft_teams.test_microsoft_teams",
			"url", integrationData.URL),
	)

	return configuration, testFunction
}

func formatMsTeamsIntegrationDataIntoConfig(data internal.MsTeamsIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_microsoft_teams" "test_microsoft_teams" {
		name = "%s"
		url  = "%s"
	}`, data.Name, data.URL)
}
