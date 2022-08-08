package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	roleSSOGroupsTestRoleName = "tf-provider-role-sso-groups-role"
)

func TestAccRoleSSOGroupsResource(t *testing.T) {
	// This resource needs to exist when the last step executes, which is an
	// Import test.
	importTestResourceName := "cyral_role_sso_groups.test_role_sso_groups"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleSSOGroupsConfig_EmptyRoleId(),
				ExpectError: regexp.MustCompile(`The argument "role_id" is required`),
			},
			{
				Config:      testAccRoleSSOGroupsConfig_EmptySSOGroup(),
				ExpectError: regexp.MustCompile(`At least 1 "sso_group" blocks are required.`),
			},
			{
				Config:      testAccRoleSSOGroupsConfig_EmptyGroupName(),
				ExpectError: regexp.MustCompile(`The argument "group_name" is required`),
			},
			{
				Config:      testAccRoleSSOGroupsConfig_EmptyIdPID(),
				ExpectError: regexp.MustCompile(`The argument "idp_id" is required`),
			},
			{
				Config: testAccRoleSSOGroupsConfig_SingleSSOGroup(),
				Check:  testAccRoleSSOGroupsCheck_SingleSSOGroup(),
			},
			{
				Config: testAccRoleSSOGroupsConfig_MultipleSSOGroups(),
				Check:  testAccRoleSSOGroupsCheck_MultipleSSOGroups(),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      importTestResourceName,
			},
		},
	})
}

func testAccRoleSSOGroupsConfig_EmptyRoleId() string {
	return `
	resource "cyral_integration_idp_okta" "test_idp_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}

	resource "cyral_role_sso_groups" "test_role_sso_groups" {
		sso_group {
			group_name="Everyone"
			idp_id=cyral_integration_idp_okta.test_idp_integration.id
		}
	}
	`
}

func testAccRoleSSOGroupsConfig_EmptySSOGroup() string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
	}

	resource "cyral_role_sso_groups" "test_role_sso_groups" {
		role_id=cyral_role.test_role.id
	}
	`, roleSSOGroupsTestRoleName)
}

func testAccRoleSSOGroupsConfig_EmptyGroupName() string {
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
			idp_id=cyral_integration_idp_okta.test_idp_integration.id
		}
	}
	`, testSingleSignOnURL, roleSSOGroupsTestRoleName)
}

func testAccRoleSSOGroupsConfig_EmptyIdPID() string {
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
		}
	}
	`, testSingleSignOnURL, roleSSOGroupsTestRoleName)
}

func testAccRoleSSOGroupsConfig_SingleSSOGroup() string {
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
	`, testSingleSignOnURL, roleSSOGroupsTestRoleName)
}

func testAccRoleSSOGroupsCheck_SingleSSOGroup() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "role_id",
			"cyral_role.test_role", "id"),
		resource.TestCheckResourceAttr("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.#", "1"),
		resource.TestCheckResourceAttrSet("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.0.id"),
		resource.TestCheckResourceAttr("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.0.group_name", "Everyone"),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.0.idp_id",
			"cyral_integration_idp_okta.test_idp_integration", "id"),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.0.idp_name",
			"cyral_integration_idp_okta.test_idp_integration", "samlp.0.display_name"),
	)
}

func testAccRoleSSOGroupsConfig_MultipleSSOGroups() string {
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
			group_name="Admin"
			idp_id=cyral_integration_idp_okta.test_idp_integration.id
		}
		sso_group {
			group_name="Dev"
			idp_id=cyral_integration_idp_okta.test_idp_integration.id
		}
	}
	`, testSingleSignOnURL, roleSSOGroupsTestRoleName)
}

func testAccRoleSSOGroupsCheck_MultipleSSOGroups() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "role_id",
			"cyral_role.test_role", "id"),
		resource.TestCheckResourceAttr("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.#", "2"),
		resource.TestCheckResourceAttrSet("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.0.id"),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.*", map[string]string{"group_name": "Admin"}),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.0.idp_id",
			"cyral_integration_idp_okta.test_idp_integration", "id"),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.0.idp_name",
			"cyral_integration_idp_okta.test_idp_integration", "samlp.0.display_name"),
		resource.TestCheckResourceAttrSet("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.1.id"),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.*", map[string]string{"group_name": "Dev"}),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.1.idp_id",
			"cyral_integration_idp_okta.test_idp_integration", "id"),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.1.idp_name",
			"cyral_integration_idp_okta.test_idp_integration", "samlp.0.display_name"),
	)
}
