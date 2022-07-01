package cyral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIntegrationIdP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationIdPConfig_EmptyFilters(),
				Check:  testAccIntegrationIdPCheck_EmptyFilters(),
			},
			{
				Config: testAccIntegrationIdPConfig_FilterTypeAAD(),
				Check:  testAccIntegrationIdPCheck_FilterTypeAAD(),
			},
		},
	})
}

func testAccIntegrationIdPConfig_EmptyFilters() string {
	return `
	data "cyral_integration_idp" "idp_integrations" {
	}
	` + testAccIdPIntegrationConfig_AAD_DefaultValues() +
		testAccIdPIntegrationConfig_ADFS_DefaultValues() +
		testAccIdPIntegrationConfig_Forgerock_DefaultValues() +
		testAccIdPIntegrationConfig_GSuite_DefaultValues() +
		testAccIdPIntegrationConfig_Okta_DefaultValues() +
		testAccIdPIntegrationConfig_PingOne_DefaultValues()
}

func testAccIntegrationIdPCheck_EmptyFilters() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.cyral_integration_idp.idp_integrations",
			"idp_set.#", "6",
		),
	)
}

func testAccIntegrationIdPConfig_FilterTypeAAD() string {
	return `
	data "cyral_integration_idp" "idp_integrations" {
		type = "aad"
	}
	` + testAccIdPIntegrationConfig_AAD_DefaultValues() +
		testAccIdPIntegrationConfig_ADFS_DefaultValues() +
		testAccIdPIntegrationConfig_Forgerock_DefaultValues() +
		testAccIdPIntegrationConfig_GSuite_DefaultValues() +
		testAccIdPIntegrationConfig_Okta_DefaultValues() +
		testAccIdPIntegrationConfig_PingOne_DefaultValues()
}

func testAccIntegrationIdPCheck_FilterTypeAAD() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.cyral_integration_idp.idp_integrations",
			"idp_list.#", "1",
		),
	)
}
