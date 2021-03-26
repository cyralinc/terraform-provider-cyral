package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// This is loosely based on this example:
// https://github.com/hashicorp/terraform-provider-vault/blob/master/vault/resource_azure_secret_backend_role_test.go

var initialConfig LogstashIntegrationData = LogstashIntegrationData{
	Endpoint:                   "logstash.local/",
	Name:                       "logstash-test",
	UseMutualAuthentication:    false,
	UsePrivateCertificateChain: false,
	UseTLS:                     false,
}

var updatedConfig LogstashIntegrationData = LogstashIntegrationData{
	Endpoint:                   "logstash-updated.local/",
	Name:                       "logstash-update-test",
	UseMutualAuthentication:    true,
	UsePrivateCertificateChain: true,
	UseTLS:                     true,
}

func TestAccLogstashIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupTest(initialConfig)
	testUpdateConfig, testUpdateFunc := setupTest(initialConfig)

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

func setupTest(integrationData LogstashIntegrationData) (string, resource.TestCheckFunc) {
	configuration := formatLogstashIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "endpoint", integrationData.Endpoint),
		resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "use_mutual_authentication", fmt.Sprintf("%t", integrationData.UseMutualAuthentication)),
		resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "use_private_certificate_chain", fmt.Sprintf("%t", integrationData.UsePrivateCertificateChain)),
		resource.TestCheckResourceAttr("cyral_integration_logstash.logstash_integration", "use_tls", fmt.Sprintf("%t", integrationData.UseTLS)),
	)

	return configuration, testFunction
}

func formatLogstashIntegrationDataIntoConfig(config LogstashIntegrationData) string {
	return fmt.Sprintf(`
	resource "cyral_integration_logstash" "logstash_integration" {
		name = "%s"
		endpoint = "%s"
		use_mutual_authentication = %t
		use_private_certificate_chain = %t
		use_tls = %t
	}`, config.Name, config.Endpoint, config.UseMutualAuthentication, config.UsePrivateCertificateChain, config.UseTLS)
}
