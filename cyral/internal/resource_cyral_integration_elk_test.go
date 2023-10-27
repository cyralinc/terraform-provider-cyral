package internal_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	integrationELKResourceName = "integration-elk"
)

var initialELKConfig internal.ELKIntegration = internal.ELKIntegration{
	Name:      utils.AccTestName(integrationELKResourceName, "ELK"),
	KibanaURL: "kibana.local",
	ESURL:     "es.local",
}

var updatedELKConfig internal.ELKIntegration = internal.ELKIntegration{
	Name:      utils.AccTestName(integrationELKResourceName, "ELK-updated"),
	KibanaURL: "kibana-update.local",
	ESURL:     "es-update.local",
}

func TestAccELKIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupELKTest(initialELKConfig)
	testUpdateConfig, testUpdateFunc := setupELKTest(updatedELKConfig)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"cyral": func() (*schema.Provider, error) {
				return provider.Provider(), nil
			},
		},
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
				ResourceName:      "cyral_integration_elk.elk_integration",
			},
		},
	})
}

func setupELKTest(integrationData internal.ELKIntegration) (string, resource.TestCheckFunc) {
	configuration := formatELKIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "kibana_url", integrationData.KibanaURL),
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "es_url", integrationData.ESURL),
	)

	return configuration, testFunction
}

func formatELKIntegrationDataIntoConfig(data internal.ELKIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_elk" "elk_integration" {
		name = "%s"
		kibana_url = "%s"
		es_url = "%s"
	}`, data.Name, data.KibanaURL, data.ESURL)
}
