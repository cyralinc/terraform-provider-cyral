package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	LookerIntegrationClientID     = "lookerClientID"
	LookerIntegrationClientSecret = "lookerClientSecret"
	LookerIntegrationURL          = "looker.local/"
)

// This is loosely based on this example:
// https://github.com/hashicorp/terraform-provider-vault/blob/master/vault/resource_azure_secret_backend_role_test.go
func TestLookerIntegrationResource(t *testing.T) {
	clientID := LookerIntegrationClientID
	clientSecret := LookerIntegrationClientSecret
	url := LookerIntegrationURL

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testLookerIntegrationInitialConfig(clientID, clientSecret, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cyral_integration_looker.looker_integration", "client_id", clientID),
					resource.TestCheckResourceAttr("cyral_integration_looker.looker_integration", "client_secret", clientSecret),
					resource.TestCheckResourceAttr("cyral_integration_looker.looker_integration", "url", url),
				),
			},
		},
	})
}

func testLookerIntegrationInitialConfig(clientID, clientSecret, url string) string {
	return fmt.Sprintf(`
resource "cyral_integration_looker" "looker_integration" {
	client_id = "%s"
	client_secret = "%s"
	url = "%s"
}`, clientID, clientSecret, url)
}
