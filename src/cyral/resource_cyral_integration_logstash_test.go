package cyral

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationLogstashResourceName = "integration-logstash"
)

var initialLogstashConfig LogstashIntegration = LogstashIntegration{
	Endpoint:                   "logstash.local/",
	Name:                       utils.AccTestName(integrationLogstashResourceName, "logstash-test"),
	UseMutualAuthentication:    false,
	UsePrivateCertificateChain: false,
	UseTLS:                     false,
}

var updated1LogstashConfig LogstashIntegration = LogstashIntegration{
	Endpoint:                   "logstash-updated.local/",
	Name:                       utils.AccTestName(integrationLogstashResourceName, "logstash-update-test"),
	UseMutualAuthentication:    true,
	UsePrivateCertificateChain: false,
	UseTLS:                     false,
}

var updated2LogstashConfig LogstashIntegration = LogstashIntegration{
	Endpoint:                   "logstash-updated.local/",
	Name:                       utils.AccTestName(integrationLogstashResourceName, "logstash-update-test"),
	UseMutualAuthentication:    false,
	UsePrivateCertificateChain: true,
	UseTLS:                     false,
}

var updated3LogstashConfig LogstashIntegration = LogstashIntegration{
	Endpoint:                   "logstash-updated.local/",
	Name:                       utils.AccTestName(integrationLogstashResourceName, "logstash-update-test"),
	UseMutualAuthentication:    false,
	UsePrivateCertificateChain: false,
	UseTLS:                     true,
}

func TestAccLogstashIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupLogstashTest(initialLogstashConfig)
	testUpdate1Config, testUpdate1Func := setupLogstashTest(updated1LogstashConfig)
	testUpdate2Config, testUpdate2Func := setupLogstashTest(updated2LogstashConfig)
	testUpdate3Config, testUpdate3Func := setupLogstashTest(updated3LogstashConfig)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdate1Config,
				Check:  testUpdate1Func,
			},
			{
				Config: testUpdate2Config,
				Check:  testUpdate2Func,
			},
			{
				Config: testUpdate3Config,
				Check:  testUpdate3Func,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_integration_logstash.logstash_integration",
			},
		},
	})
}

func setupLogstashTest(integrationData LogstashIntegration) (string, resource.TestCheckFunc) {
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

func formatLogstashIntegrationDataIntoConfig(config LogstashIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_logstash" "logstash_integration" {
		name = "%s"
		endpoint = "%s"
		use_mutual_authentication = %t
		use_private_certificate_chain = %t
		use_tls = %t
	}`, config.Name, config.Endpoint, config.UseMutualAuthentication, config.UsePrivateCertificateChain, config.UseTLS)
}
