package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialELKConfig ELKIntegration = ELKIntegration{
	Name:      "unitTest-ELK",
	KibanaURL: "kibana.local",
	ESURL:     "es.local",
}

var updatedELKConfig ELKIntegration = ELKIntegration{
	Name:      "unitTest-ELK-Updated",
	KibanaURL: "kibana-update.local",
	ESURL:     "es-update.local",
}

func TestAccELKIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupELKTest(initialELKConfig)
	testUpdateConfig, testUpdateFunc := setupELKTest(updatedELKConfig)

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

func setupELKTest(integrationData ELKIntegration) (string, resource.TestCheckFunc) {
	configuration := formatELKIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "kibana_url", integrationData.KibanaURL),
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "es_url", integrationData.ESURL),
	)

	return configuration, testFunction
}

func formatELKIntegrationDataIntoConfig(data ELKIntegration) string {
	return fmt.Sprintf(`
resource "cyral_integration_elk" "elk_integration" {
	name = "%s"
	kibana_url = "%s"
	es_url = "%s"
}`, data.Name, data.KibanaURL, data.ESURL)
}
