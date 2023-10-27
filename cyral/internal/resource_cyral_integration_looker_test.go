package internal_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialLookerConfig internal.LookerIntegration = internal.LookerIntegration{
	ClientId:     "lookerClientID",
	ClientSecret: "lookerClientSecret",
	URL:          "looker.local/",
}

var updatedLookerConfig internal.LookerIntegration = internal.LookerIntegration{
	ClientId:     "lookerClientIDUpdated",
	ClientSecret: "lookerClientSecretUpdated",
	URL:          "looker-updated.local/",
}

func TestAccLookerIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupLookerTest(initialLookerConfig)
	testUpdateConfig, testUpdateFunc := setupLookerTest(updatedLookerConfig)

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
				ResourceName:      "cyral_integration_looker.looker_integration",
			},
		},
	})
}

func setupLookerTest(integrationData internal.LookerIntegration) (string, resource.TestCheckFunc) {
	configuration := formatLookerIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_looker.looker_integration", "client_id", integrationData.ClientId),
		resource.TestCheckResourceAttr("cyral_integration_looker.looker_integration", "client_secret", integrationData.ClientSecret),
		resource.TestCheckResourceAttr("cyral_integration_looker.looker_integration", "url", integrationData.URL),
	)

	return configuration, testFunction
}

func formatLookerIntegrationDataIntoConfig(data internal.LookerIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_looker" "looker_integration" {
		client_id = "%s"
		client_secret = "%s"
		url = "%s"
	}`, data.ClientId, data.ClientSecret, data.URL)
}
