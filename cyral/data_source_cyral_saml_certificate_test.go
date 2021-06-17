package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSAMLCertificateDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			testCertificateSetup(),
		},
	})
}

func testCertificateSetup() resource.TestStep {
	format := `
		data "cyral_saml_certificate" "certificate" {}

		output "certificate" {
			value = data.cyral_saml_certificate.certificate
		}
	`

	return resource.TestStep{
		Config: format,
		Check: testOutputAttrFunction("certificate", func(s interface{}) error {
			if s == nil {
				return fmt.Errorf("certificate response empty")
			}
			ma, ok := s.(map[string]interface{})
			if !ok {
				return fmt.Errorf("wrong data type. value = %s", s)
			}
			str, ok := ma["certificate"].(string)
			if !ok {
				return fmt.Errorf("wrong data type. value = %s", s)
			}
			if len(str) == 0 {
				return fmt.Errorf("certificate empty")
			}
			return nil
		}),
	}

}

func testDataSourceAttrFunction(resource, attr string, checkFunc func(string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		at, ok := rs.Primary.Attributes[attr]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		return checkFunc(at)
	}
}

func testOutputAttrFunction(resource string, checkFunc func(interface{}) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Outputs[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		at := rs.Value

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		return checkFunc(at)
	}
}

func testResourceAttrFunction(resource, attr string, checkFunc func(string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		at, ok := rs.Primary.Attributes[attr]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		return checkFunc(at)
	}
}
