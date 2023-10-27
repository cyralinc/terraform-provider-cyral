package internal_test

import (
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSAMLCertificateDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"cyral": func() (*schema.Provider, error) {
				return provider.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSAMLCertificateConfig(),
				Check:  testAccSAMLCertificateCheck(),
			},
		},
	})
}

func testAccSAMLCertificateConfig() string {
	return `
	data "cyral_saml_certificate" "test_saml_certificate" {
	}
	`
}

func testAccSAMLCertificateCheck() resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet("data.cyral_saml_certificate.test_saml_certificate",
		"certificate")
}
