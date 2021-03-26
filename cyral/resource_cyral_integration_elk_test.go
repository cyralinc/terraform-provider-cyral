package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	ELKIntegrationName      = "unitTest-ELK"
	ELKIntegrationKibanaURL = "kibana.local"
	ELKIntegrationESURL     = "es.local"
)

// This is loosely based on this example:
// https://github.com/hashicorp/terraform-provider-vault/blob/master/vault/resource_azure_secret_backend_role_test.go
func TestAccELKIntegrationResource(t *testing.T) {
	name := ELKIntegrationName
	kibanaURL := ELKIntegrationKibanaURL
	esURL := ELKIntegrationESURL

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testELKIntegrationInitialConfig(name, kibanaURL, esURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "name", name),
					resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "kibana_url", kibanaURL),
					resource.TestCheckResourceAttr("cyral_integration_elk.elk_integration", "es_url", esURL),
				),
			},
		},
	})
}

func testELKIntegrationInitialConfig(name, kibanaURL, esURL string) string {
	return fmt.Sprintf(`
resource "cyral_integration_elk" "elk_integration" {
	name = "%s"
	kibana_url = "%s"
	es_url = "%s"
}`, name, kibanaURL, esURL)
}
