package deprecated_test

import (
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/deprecated"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationELKResourceName = "integration-elk"
)

var initialELKConfig deprecated.ELKIntegration = deprecated.ELKIntegration{
	Name:      utils.AccTestName(integrationELKResourceName, "ELK"),
	KibanaURL: "kibana.local",
	ESURL:     "es.local",
}

var updatedELKConfig deprecated.ELKIntegration = deprecated.ELKIntegration{
	Name:      utils.AccTestName(integrationELKResourceName, "ELK-updated"),
	KibanaURL: "kibana-update.local",
	ESURL:     "es-update.local",
}

func TestAccELKIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupELKTest(initialELKConfig)
	testUpdateConfig, testUpdateFunc := setupELKTest(updatedELKConfig)

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
				ResourceName:      "cyral_integration_elk.elk_integration",
			},
		},
	})
}

func setupELKTest(d deprecated.ELKIntegration) (string, resource.TestCheckFunc) {
	configuration := utils.FormatELKIntegrationDataIntoConfig(d.Name, d.KibanaURL, d.ESURL)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "name", d.Name),
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "kibana_url", d.KibanaURL),
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "es_url", d.ESURL),
	)

	return configuration, testFunction
}
