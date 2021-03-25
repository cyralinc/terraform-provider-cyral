package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	LogstashIntegrationEndpoint                   = "logstash.local/"
	LogstashIntegrationName                       = "logstash-test"
	LogstashIntegrationUseMutualAuthentication    = false
	LogstashIntegrationUsePrivateCertificateChain = false
	LogstashIntegrationUseTLS                     = false
)

// This is loosely based on this example:
// https://github.com/hashicorp/terraform-provider-vault/blob/master/vault/resource_azure_secret_backend_role_test.go
func TestAccLogstashIntegrationResource(t *testing.T) {
	name := LogstashIntegrationName
	endpoint := LogstashIntegrationEndpoint
	use_mutual_authentication := LogstashIntegrationUseMutualAuthentication
	use_private_certificate_chain := LogstashIntegrationUsePrivateCertificateChain
	use_tls := LogstashIntegrationUseTLS

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testLogstashIntegrationInitialConfig(name, endpoint, use_mutual_authentication, use_private_certificate_chain, use_tls),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "name", name),
					resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "endpoint", endpoint),
					resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "use_mutual_authentication", fmt.Sprintf("%t", use_mutual_authentication)),
					resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "use_private_certificate_chain", fmt.Sprintf("%t", use_private_certificate_chain)),
					resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "use_tls", fmt.Sprintf("%t", use_tls)),
				),
			},
		},
	})
}

func testLogstashIntegrationInitialConfig(name string, endpoint string, use_mutual_authentication bool, use_private_certificate_chain bool, use_tls bool) string {
	return fmt.Sprintf(`
	resource "cyral_integration_logstash" "logstash_integration" {
		name = "%s"
		endpoint = "%s"
		use_mutual_authentication = %t
		use_private_certificate_chain = %t
		use_tls = %t
	}`, name, endpoint, use_mutual_authentication, use_private_certificate_chain, use_tls)
}
