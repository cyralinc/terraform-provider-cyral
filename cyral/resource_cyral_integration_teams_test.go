package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialTeamsConfig TeamsIntegrationData = TeamsIntegrationData{
	Name: "tf-test-teams-alerts",
	Url:  "https://teams.local",
}

var updatedTeamsConfig TeamsIntegrationData = TeamsIntegrationData{
	Name: "tf-test-teams-alerts",
	Url:  "https://teams-updated.local",
}

func TestTeamsIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupTeamsTest(initialTeamsConfig)
	testUpdateConfig, testUpdateFunc := setupTeamsTest(initialTeamsConfig)

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

func setupTeamsTest(integrationData TeamsIntegrationData) (string, resource.TestCheckFunc) {
	configuration := formatTeamsIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_microsoft_teams.test_microsoft_teams", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_microsoft_teams.test_microsoft_teams", "url", integrationData.Url),
	)

	return configuration, testFunction
}

func formatTeamsIntegrationDataIntoConfig(data TeamsIntegrationData) string {
	return fmt.Sprintf(`
	resource "cyral_integration_microsoft_teams" "test_microsoft_teams" {
		name = "%s"
		url  = "%s"
	}`, data.Name, data.Url)
}
