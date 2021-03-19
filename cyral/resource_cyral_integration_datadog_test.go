package cyral

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	EnvVarDDIntegrationName   = "TEST_TPC_DD_INTEGRATION_NAME"
	EnvVarDDIntegrationApiKey = "TEST_TPC_DD_INTEGRATION_API_KEY"
)

// This is loosely based on this example:
// https://github.com/hashicorp/terraform-provider-vault/blob/master/vault/resource_azure_secret_backend_role_test.go
func TestDatadogIntegrationResource(t *testing.T) {
	name := os.Getenv(EnvVarDDIntegrationName)
	apiKey := os.Getenv(EnvVarDDIntegrationApiKey)
	if name == "" || apiKey == "" {
		t.Skipf("skipping because both %q and %q must be set", EnvVarDDIntegrationName, EnvVarDDIntegrationApiKey)
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testDDIntegrationInitialConfig(name, apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration", "name", name),
					resource.TestCheckResourceAttr("cyral_integration_datadog.datadog_integration", "api_key", apiKey),
				),
			},
		},
	})
}

func testDDIntegrationInitialConfig(name, apiKey string) string {
	return fmt.Sprintf(`
resource "cyral_integration_datadog" "datadog_integration" {
    name = "%s"
    api_key = "%s"
}`, name, apiKey)
}
