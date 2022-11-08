package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	roleResourceName = "role"
)

func initialRoleName() string {
	return accTestName(roleResourceName, "role")
}

func updatedRoleName() string {
	return accTestName(roleResourceName, "role-updated")
}

var onlyFalsePermissions = map[string]string{
	"modify_sidecars_and_repositories": "false",
	"modify_users":                     "false",
	"modify_policies":                  "false",
	"view_audit_logs":                  "false",
	"modify_integrations":              "false",
	"modify_roles":                     "false",
	"view_datamaps":                    "false",
}

var trueAndFalsePermissions = map[string]string{
	"modify_sidecars_and_repositories": "true",
	"modify_users":                     "true",
	"modify_policies":                  "true",
	"view_audit_logs":                  "false",
	"modify_integrations":              "false",
	"modify_roles":                     "false",
	"view_datamaps":                    "false",
}

var onlyTruePermissions = map[string]string{
	"modify_sidecars_and_repositories": "true",
	"modify_users":                     "true",
	"modify_policies":                  "true",
	"view_audit_logs":                  "true",
	"modify_integrations":              "true",
	"modify_roles":                     "true",
	"view_datamaps":                    "true",
}

func TestAccRoleResource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleConfig_EmptyRoleName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccRoleConfig_MultiplePermissionsBlock(),
				ExpectError: regexp.MustCompile(`No more than 1 "permissions" blocks are allowed`),
			},
			{
				Config: testAccRoleConfig_DefaultValues(),
				Check:  testAccRoleCheck_DefaultValues(),
			},
			{
				Config: testAccRoleConfig_EmptyPermissions(),
				Check:  testAccRoleCheck_EmptyPermissions(),
			},
			{
				Config: testAccRoleConfig_OnlyFalsePermissions(),
				Check:  testAccRoleCheck_OnlyFalsePermissions(),
			},
			{
				Config: testAccRoleConfig_TrueAndFalsePermissions(),
				Check:  testAccRoleCheck_TrueAndFalsePermissions(),
			},
			{
				Config: testAccRoleConfig_OnlyTruePermissions(),
				Check:  testAccRoleCheck_OnlyTruePermissions(),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_role.test_role",
			},
		},
	})
}

func testAccRoleConfig_EmptyRoleName() string {
	return `
	resource "cyral_role" "test_role" {
	}
	`
}

func testAccRoleConfig_MultiplePermissionsBlock() string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
		permissions{
		}
		permissions{
		}
	}
	`, initialRoleName())
}

func testAccRoleConfig_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
	}
	`, initialRoleName())
}

func testAccRoleCheck_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_role.test_role", "name", initialRoleName()),
		resource.TestCheckResourceAttr("cyral_role.test_role", "permissions.#", "0"),
	)
}

func testAccRoleConfig_EmptyPermissions() string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
		permissions {
		}
	}
	`, updatedRoleName())
}

func testAccRoleCheck_EmptyPermissions() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_role.test_role", "name", updatedRoleName()),
		resource.TestCheckResourceAttr("cyral_role.test_role", "permissions.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_role.test_role", "permissions.*",
			onlyFalsePermissions),
	)
}

func testAccRoleConfig_OnlyFalsePermissions() string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
		permissions {
			modify_sidecars_and_repositories = false
			modify_users = false
			modify_policies = false
			view_audit_logs = false
			modify_integrations = false
			modify_roles = false
			view_datamaps = false
		}
	}
	`, updatedRoleName())
}

func testAccRoleCheck_OnlyFalsePermissions() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_role.test_role", "name", updatedRoleName()),
		resource.TestCheckResourceAttr("cyral_role.test_role", "permissions.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_role.test_role", "permissions.*",
			onlyFalsePermissions),
	)
}

func testAccRoleConfig_TrueAndFalsePermissions() string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
		permissions {
			modify_sidecars_and_repositories = true
			modify_users = true
			modify_policies = true
			view_audit_logs = false
			modify_integrations = false
			modify_roles = false
			view_datamaps = false
		}
	}
	`, updatedRoleName())
}

func testAccRoleCheck_TrueAndFalsePermissions() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_role.test_role", "name", updatedRoleName()),
		resource.TestCheckResourceAttr("cyral_role.test_role", "permissions.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_role.test_role", "permissions.*",
			trueAndFalsePermissions),
	)
}

func testAccRoleConfig_OnlyTruePermissions() string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
		permissions {
			modify_sidecars_and_repositories = true
			modify_users = true
			modify_policies = true
			view_audit_logs = true
			modify_integrations = true
			modify_roles = true
			view_datamaps = true
		}
	}
	`, updatedRoleName())
}

func testAccRoleCheck_OnlyTruePermissions() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_role.test_role", "name", updatedRoleName()),
		resource.TestCheckResourceAttr("cyral_role.test_role", "permissions.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_role.test_role", "permissions.*",
			onlyTruePermissions),
	)
}
