package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialHCVaultIntegrationConfig HCVaultIntegration = HCVaultIntegration{
	AuthMethod: "unitTest-auth_method",
	ID:         "unitTest-id",
	AuthType:   "unitTest-auth_type",
	Name:       "unitTest-name",
	Server:     "unitTest-server",
}

var updatedHCVaultIntegrationConfig HCVaultIntegration = HCVaultIntegration{
	AuthMethod: "unitTest-auth_method-updated",
	ID:         "unitTest-id-updated",
	AuthType:   "unitTest-auth_type-updated",
	Name:       "unitTest-name-updated",
	Server:     "unitTest-server-updated",
}

func TestAccHCVaultIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupHCVaultIntegrationTest(initialHCVaultIntegrationConfig)
	testUpdateConfig, testUpdateFunc := setupHCVaultIntegrationTest(updatedHCVaultIntegrationConfig)

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

func setupHCVaultIntegrationTest(integrationData HCVaultIntegration) (string, resource.TestCheckFunc) {
	configuration := formatHCVaultIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(

		resource.TestCheckResourceAttr("cyral_integration_hc_vault.hc_vault_integration", "auth_method", integrationData.AuthMethod),
		resource.TestCheckResourceAttr("cyral_integration_hc_vault.hc_vault_integration", "auth_type", integrationData.AuthType),
		resource.TestCheckResourceAttr("cyral_integration_hc_vault.hc_vault_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_hc_vault.hc_vault_integration", "server", integrationData.Server),
	)

	return configuration, testFunction
}

func formatHCVaultIntegrationDataIntoConfig(data HCVaultIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_hc_vault" "hc_vault_integration" {
		auth_method = "%s" 
		auth_type = "%s" 
		name = "%s" 
		server = "%s" 
	}`, data.AuthMethod, data.AuthType, data.Name, data.Server)
}
