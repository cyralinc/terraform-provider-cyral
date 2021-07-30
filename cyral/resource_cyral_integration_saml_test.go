package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialSAMLIntegrationConfig SAMLIntegration = SAMLIntegration{}

var updatedSAMLIntegrationConfig SAMLIntegration = SAMLIntegration{}

func TestAccSAMLIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupSAMLIntegrationTest(initialSAMLIntegrationConfig)
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
	})
}

func setupSAMLIntegrationTest(integrationData SAMLIntegration) (string, resource.TestCheckFunc) {
	configuration := formatSAMLIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc()

	return configuration, testFunction
}

func formatSAMLIntegrationDataIntoConfig(data SAMLIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_saml_integration" "saml_integration" {
	
	}`)
}
