package internal_test

import (
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

func setupELKTest(d internal.ELKIntegration) (string, resource.TestCheckFunc) {
	configuration := utils.FormatELKIntegrationDataIntoConfig(d.Name, d.KibanaURL, d.ESURL)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "name", d.Name),
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "kibana_url", d.KibanaURL),
		resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "es_url", d.ESURL),
	)

	return configuration, testFunction
}
