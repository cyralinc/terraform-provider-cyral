package hcvault_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/hcvault"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	integrationHCVaultResourceName = "integration-hcvault"
)

var initialHCVaultIntegrationConfig hcvault.HCVaultIntegration = hcvault.HCVaultIntegration{
	AuthMethod: "unitTest-auth_method",
	ID:         "unitTest-id",
	AuthType:   "unitTest-auth_type",
	Name:       utils.AccTestName(integrationHCVaultResourceName, "hcvault"),
	Server:     "unitTest-server",
}

var updatedHCVaultIntegrationConfig hcvault.HCVaultIntegration = hcvault.HCVaultIntegration{
	AuthMethod: "unitTest-auth_method-updated",
	ID:         "unitTest-id-updated",
	AuthType:   "unitTest-auth_type-updated",
	Name:       utils.AccTestName(integrationHCVaultResourceName, "hcvault-updated"),
	Server:     "unitTest-server-updated",
}

func TestAccHCVaultIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupHCVaultIntegrationTest(initialHCVaultIntegrationConfig)
	testUpdateConfig, testUpdateFunc := setupHCVaultIntegrationTest(updatedHCVaultIntegrationConfig)

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
				ResourceName:      "cyral_integration_hc_vault.hc_vault_integration",
			},
		},
	})
}

func setupHCVaultIntegrationTest(integrationData hcvault.HCVaultIntegration) (string, resource.TestCheckFunc) {
	configuration := formatHCVaultIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(

		resource.TestCheckResourceAttr("cyral_integration_hc_vault.hc_vault_integration", "auth_method", integrationData.AuthMethod),
		resource.TestCheckResourceAttr("cyral_integration_hc_vault.hc_vault_integration", "auth_type", integrationData.AuthType),
		resource.TestCheckResourceAttr("cyral_integration_hc_vault.hc_vault_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_hc_vault.hc_vault_integration", "server", integrationData.Server),
	)

	return configuration, testFunction
}

func formatHCVaultIntegrationDataIntoConfig(data hcvault.HCVaultIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_hc_vault" "hc_vault_integration" {
		auth_method = "%s"
		auth_type = "%s"
		name = "%s"
		server = "%s"
	}`, data.AuthMethod, data.AuthType, data.Name, data.Server)
}
