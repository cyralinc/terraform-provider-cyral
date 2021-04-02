package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialLookerConfig LookerIntegrationData = LookerIntegrationData{
	ClientId:     "lookerClientID",
	ClientSecret: "lookerClientSecret",
	URL:          "looker.local/",
}

var updatedLookerConfig LookerIntegrationData = LookerIntegrationData{
	ClientId:     "lookerClientIDUpdated",
	ClientSecret: "lookerClientSecretUpdated",
	URL:          "looker-updated.local/",
}

func TestAccLookerIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupLookerTest(initialLookerConfig)
	testUpdateConfig, testUpdateFunc := setupLookerTest(updatedLookerConfig)

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

func setupLookerTest(integrationData LookerIntegrationData) (string, resource.TestCheckFunc) {
	configuration := formatLookerIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_looker.looker_integration", "client_id", integrationData.ClientId),
		resource.TestCheckResourceAttr("cyral_integration_looker.looker_integration", "client_secret", integrationData.ClientSecret),
		resource.TestCheckResourceAttr("cyral_integration_looker.looker_integration", "url", integrationData.URL),
	)

	return configuration, testFunction
}

func formatLookerIntegrationDataIntoConfig(data LookerIntegrationData) string {
	return fmt.Sprintf(`
resource "cyral_integration_looker" "looker_integration" {
	client_id = "%s"
	client_secret = "%s"
	url = "%s"
}`, data.ClientId, data.ClientSecret, data.URL)
}
