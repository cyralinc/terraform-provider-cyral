package mfaduo_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationDuoMFAResourceName = "integration-duo-mfa"
)

func initialDuoMFAIntegrationConfig() *confextension.IntegrationConfExtension {
	integration := confextension.NewIntegrationConfExtension(confextension.DuoMFATemplateType)
	integration.Name = utils.AccTestName(integrationDuoMFAResourceName, "integration")
	integration.Parameters = `{"integrationKey": "integration-key-1", "secretKey": "secret-key-1", "apiHostname": "api-hostname-1"}`
	return integration
}

func updatedDuoMFAIntegrationConfig() *confextension.IntegrationConfExtension {
	integration := confextension.NewIntegrationConfExtension(confextension.DuoMFATemplateType)
	integration.Name = utils.AccTestName(integrationDuoMFAResourceName, "integration-updated")
	integration.Parameters = `{"integrationKey": "integration-key-2", "secretKey": "secret-key-2", "apiHostname": "api-hostname-2"}`
	return integration
}

func TestAccDuoMFAIntegrationResource(t *testing.T) {
	tfResourceName := "duo_integration"

	testConfig, testFunc := setupDuoMFAIntegrationTest(t,
		initialDuoMFAIntegrationConfig(), tfResourceName)
	testUpdateConfig, testUpdateFunc := setupDuoMFAIntegrationTest(t,
		updatedDuoMFAIntegrationConfig(), tfResourceName)

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
				ResourceName:      fmt.Sprintf("cyral_integration_mfa_duo.%s", tfResourceName),
			},
		},
	})
}

func setupDuoMFAIntegrationTest(
	t *testing.T,
	integrationData *confextension.IntegrationConfExtension,
	resourceName string,
) (string, resource.TestCheckFunc) {

	var parameters confextension.IntegrationConfExtensionParameters
	err := json.Unmarshal([]byte(integrationData.Parameters), &parameters)
	if err != nil {
		t.Fatalf("Error unmarshalling parameters: %v", err)
	}

	configuration := formatDuoMFAIntegrationIntoConfig(
		integrationData.Name, parameters.IntegrationKey,
		parameters.SecretKey, parameters.APIHostname)

	resourceFullName := fmt.Sprintf("cyral_integration_mfa_duo.%s", resourceName)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", integrationData.Name),
		resource.TestCheckResourceAttr(resourceFullName,
			"integration_key", parameters.IntegrationKey),
		resource.TestCheckResourceAttr(resourceFullName,
			"secret_key", parameters.SecretKey),
		resource.TestCheckResourceAttr(resourceFullName,
			"api_hostname", parameters.APIHostname),
	)

	return configuration, testFunction
}

func formatDuoMFAIntegrationIntoConfig(name, integrationKey, secretKey, apiHostname string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_mfa_duo" "duo_integration" {
		name            = "%s"
		integration_key = "%s"
		secret_key      = "%s"
		api_hostname    = "%s"
	}`, name, integrationKey, secretKey, apiHostname)
}
