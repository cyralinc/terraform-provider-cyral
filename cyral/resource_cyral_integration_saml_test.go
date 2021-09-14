package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialSAMLIntegrationConfig SAMLSetting = SAMLSetting{}

var updatedSAMLIntegrationConfig SAMLSetting = SAMLSetting{}

func TestAccSAMLIntegrationResource(t *testing.T) {
	/* testConfig, testFunc := setupSAMLIntegrationTest(initialSAMLIntegrationConfig)
	testUpdateConfig, testUpdateFunc := setupSAMLIntegrationTest(updatedSAMLIntegrationConfig)

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
	}) */
}

func setupSAMLIntegrationTest(integrationData SAMLSetting) (string, resource.TestCheckFunc) {
	configuration := formatSAMLIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc()

	return configuration, testFunction
}

func formatSAMLIntegrationDataIntoConfig(data SAMLSetting) string {
	return fmt.Sprintf(`
	resource "cyral_saml_integration" "saml_integration" {
	
	}`)
}
