package cyral

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIntegrationIdP(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationIdPConfig_FilterTypeAAD(),
				Check:  testAccIntegrationIdPCheck_FilterTypeAAD(),
			},
		},
	})
}

func testAccIntegrationIdPConfig_FilterTypeAAD() string {
	return `
	resource "cyral_integration_idp_adfs" "idp_integration_adfs" {
		samlp {
			display_name = "testAcc_ADFS"
			config {
				single_sign_on_service_url = "https://testAcc_ADFS.com"
			}
		}
	}
	resource "cyral_integration_idp_okta" "idp_integration_okta" {
		samlp {
			display_name = "testAcc_OKTA"
			config {
				single_sign_on_service_url = "https://testAcc_OKTA.com"
			}
		}
	}
	data "cyral_integration_idp" "idp_integrations" {
		display_name = "testAcc_ADFS"
	}
	`
}

func testAccIntegrationIdPCheck_FilterTypeAAD() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.cyral_integration_idp.idp_integrations",
			"idp_list.#", "1",
		),
	)
}
*/
