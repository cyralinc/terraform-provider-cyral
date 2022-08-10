package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationOktaResourceName = "integration-okta"
)

var initialOktaConfig ResourceIntegrationOktaPayload = ResourceIntegrationOktaPayload{
	Samlp: ResourceIntegrationOkta{
		Name:         accTestName(integrationOktaResourceName, "okta"),
		Certificate:  "certificate",
		EmailDomains: []string{"sigin.com"},
		SignInUrl:    "https://sigin.com/in",
		SignOutUrl:   "https://signout.com/out",
	},
}

var updatedOktaConfig ResourceIntegrationOktaPayload = ResourceIntegrationOktaPayload{
	Samlp: ResourceIntegrationOkta{
		Name:         accTestName(integrationOktaResourceName, "okta"),
		Certificate:  "certificate-updated",
		EmailDomains: []string{"siginupdated.com"},
		SignInUrl:    "https://siginupdated.com/in",
		SignOutUrl:   "https://signoutupdated.com/out",
	},
}

func TestAccOktaIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupOktaTest(initialOktaConfig)
	testUpdateConfig, testUpdateFunc := setupOktaTest(updatedOktaConfig)

	resource.ParallelTest(t, resource.TestCase{
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

func setupOktaTest(integrationData ResourceIntegrationOktaPayload) (string, resource.TestCheckFunc) {
	configuration := formatOktaIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_okta.tf_test_okta",
			"name", integrationData.Samlp.Name),
		resource.TestCheckResourceAttr("cyral_integration_okta.tf_test_okta",
			"certificate", integrationData.Samlp.Certificate),
		resource.TestCheckResourceAttr("cyral_integration_okta.tf_test_okta",
			"email_domains.#", fmt.Sprintf("%d", len(integrationData.Samlp.EmailDomains))),
		resource.TestCheckResourceAttr("cyral_integration_okta.tf_test_okta",
			"email_domains.0", integrationData.Samlp.EmailDomains[0]),
		resource.TestCheckResourceAttr("cyral_integration_okta.tf_test_okta",
			"signin_url", integrationData.Samlp.SignInUrl),
		resource.TestCheckResourceAttr("cyral_integration_okta.tf_test_okta",
			"signout_url", integrationData.Samlp.SignOutUrl),
	)

	return configuration, testFunction
}

func formatOktaIntegrationDataIntoConfig(data ResourceIntegrationOktaPayload) string {
	return fmt.Sprintf(`
	resource "cyral_integration_okta" "tf_test_okta" {
		name          = "%s"
		certificate   = "%s"
		email_domains = [%s]
		signin_url    = "%s"
		signout_url   = "%s"
	  }`, data.Samlp.Name, data.Samlp.Certificate, listToStr(data.Samlp.EmailDomains), data.Samlp.SignInUrl, data.Samlp.SignOutUrl)
}
