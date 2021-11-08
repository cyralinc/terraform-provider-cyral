package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRoleSSOGroupsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleSSOGroupsConfig_DefaultValues(),
				Check:  testAccRoleSSOGroupsCheck_DefaultValues(),
			},
		},
	})
}

func testAccRoleSSOGroupsConfig_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "test_idp_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}

	resource "cyral_role" "test_role" {
		name="%s"
	}

	resource "cyral_role_sso_groups" "test_role_sso_groups" {
		role_id=cyral_role.test_role.id
		sso_group {
			group_name="Everyone"
			idp_id=cyral_integration_idp_okta.test_idp_integration.id
		}
	}
	`, testSingleSignOnURL, initialRoleName)
}

func testAccRoleSSOGroupsCheck_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair("cyral_role_sso_groups.test_role_sso_groups", "role_id",
			"cyral_role.test_role", "id"),
		resource.TestCheckResourceAttr("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.#", "1"),
		resource.TestCheckResourceAttrSet("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.0.id"),
		resource.TestCheckResourceAttr("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.0.group_name", "Everyone"),
		resource.TestCheckResourceAttrPair("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.0.idp_id", "cyral_integration_idp_okta.test_idp_integration", "id"),
	)
}
