package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialSplunkConfig SplunkIntegration = SplunkIntegration{
	Name:        "splunk-test",
	AccessToken: "access-token",
	Port:        3333,
	Host:        "splunk.local",
	Index:       "index",
	UseTLS:      false,
}

var updatedSplunkConfig SplunkIntegration = SplunkIntegration{
	Name:        "splunk-test-update",
	AccessToken: "access-token-update",
	Port:        6666,
	Host:        "splunk-update.local",
	Index:       "index-update",
	UseTLS:      true,
}

func TestAccSplunkIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupSplunkTest(initialSplunkConfig)
	testUpdateConfig, testUpdateFunc := setupSplunkTest(updatedSplunkConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
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

func setupSplunkTest(integrationData SplunkIntegration) (string, resource.TestCheckFunc) {
	configuration := formatSplunkIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_splunk.splunk_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_integration_splunk.splunk_integration", "access_token", integrationData.AccessToken),
		resource.TestCheckResourceAttr("cyral_integration_splunk.splunk_integration", "port", fmt.Sprintf("%d", integrationData.Port)),
		resource.TestCheckResourceAttr("cyral_integration_splunk.splunk_integration", "host", integrationData.Host),
		resource.TestCheckResourceAttr("cyral_integration_splunk.splunk_integration", "index", integrationData.Index),
		resource.TestCheckResourceAttr("cyral_integration_splunk.splunk_integration", "use_tls", fmt.Sprintf("%t", integrationData.UseTLS)),
	)

	return configuration, testFunction
}

func formatSplunkIntegrationDataIntoConfig(data SplunkIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_splunk" "splunk_integration" {
		name = "%s"
		access_token = "%s"
		port = %d
		host = "%s"
		index = "%s"
		use_tls = %t
	}`, data.Name, data.AccessToken, data.Port, data.Host, data.Index, data.UseTLS)
}
