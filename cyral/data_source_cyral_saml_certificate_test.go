package cyral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSAMLCertificateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
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
