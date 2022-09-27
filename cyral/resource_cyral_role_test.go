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
	"view_sidecars_and_repositories":   "false",
	"view_audit_logs":                  "false",
	"modify_integrations":              "false",
	"modify_roles":                     "false",
	"view_datamaps":                    "false",
}

var trueAndFalsePermissions = map[string]string{
	"modify_sidecars_and_repositories": "true",
	"modify_users":                     "true",
	"modify_policies":                  "true",
	"view_sidecars_and_repositories":   "true",
	"view_audit_logs":                  "false",
	"modify_integrations":              "false",
	"modify_roles":                     "false",
	"view_datamaps":                    "false",
}

var onlyTruePermissions = map[string]string{
	"modify_sidecars_and_repositories": "true",
	"modify_users":                     "true",
	"modify_policies":                  "true",
	"view_sidecars_and_repositories":   "true",
	"view_audit_logs":                  "true",
	"modify_integrations":              "true",
	"modify_roles":                     "true",
	"view_datamaps":                    "true",
}

func TestAccRoleResource(t *testing.T) {
	initialResName := "initial_role"
	updatedResName := "initial_role"

	importResourceName := fmt.Sprintf("cyral_role.%s", updatedResName)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleConfig_EmptyRoleName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccRoleConfig_MultiplePermissionsBlock(initialResName),
				ExpectError: regexp.MustCompile(`No more than 1 "permissions" blocks are allowed`),
			},
			{
				Config: testAccRoleConfig_DefaultValues(initialResName),
				Check:  testAccRoleCheck_DefaultValues(initialResName),
			},
			{
				Config: testAccRoleConfig_EmptyPermissions(updatedResName),
				Check:  testAccRoleCheck_EmptyPermissions(updatedResName),
			},
			{
				Config: testAccRoleConfig_OnlyFalsePermissions(updatedResName),
				Check:  testAccRoleCheck_OnlyFalsePermissions(updatedResName),
			},
			{
				Config: testAccRoleConfig_TrueAndFalsePermissions(updatedResName),
				Check:  testAccRoleCheck_TrueAndFalsePermissions(updatedResName),
			},
			{
				Config: testAccRoleConfig_OnlyTruePermissions(updatedResName),
				Check:  testAccRoleCheck_OnlyTruePermissions(updatedResName),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      importResourceName,
			},
		},
	})
}

func testRoleResourceFullName(resName string) string {
	return fmt.Sprintf("cyral_role.%s", resName)
}

func testAccRoleConfig_EmptyRoleName() string {
	return `
	resource "cyral_role" "test_role" {
	}
	`
}

func testAccRoleConfig_MultiplePermissionsBlock(resName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "%s" {
		name="%s"
		permissions{
		}
		permissions{
		}
	}
	`, resName, initialRoleName())
}

func testAccRoleConfig_DefaultValues(resName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "%s" {
		name="%s"
	}
	`, resName, initialRoleName())
}

func testAccRoleCheck_DefaultValues(resName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "name", initialRoleName()),
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "permissions.#", "0"),
	)
}

func testAccRoleConfig_EmptyPermissions(resName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "%s" {
		name="%s"
		permissions {
		}
	}
	`, resName, updatedRoleName())
}

func testAccRoleCheck_EmptyPermissions(resName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "name", updatedRoleName()),
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "permissions.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(testRoleResourceFullName(resName), "permissions.*",
			onlyFalsePermissions),
	)
}

func testAccRoleConfig_OnlyFalsePermissions(resName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "%s" {
		name="%s"
		permissions {
			modify_sidecars_and_repositories = false
			modify_users = false
			modify_policies = false
			view_sidecars_and_repositories = false
			view_audit_logs = false
			modify_integrations = false
			modify_roles = false
			view_datamaps = false
		}
	}
	`, resName, updatedRoleName())
}

func testAccRoleCheck_OnlyFalsePermissions(resName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "name", updatedRoleName()),
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "permissions.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(testRoleResourceFullName(resName), "permissions.*",
			onlyFalsePermissions),
	)
}

func testAccRoleConfig_TrueAndFalsePermissions(resName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "%s" {
		name="%s"
		permissions {
			modify_sidecars_and_repositories = true
			modify_users = true
			modify_policies = true
			view_sidecars_and_repositories = true
			view_audit_logs = false
			modify_integrations = false
			modify_roles = false
			view_datamaps = false
		}
	}
	`, resName, updatedRoleName())
}

func testAccRoleCheck_TrueAndFalsePermissions(resName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "name", updatedRoleName()),
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "permissions.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(testRoleResourceFullName(resName), "permissions.*",
			trueAndFalsePermissions),
	)
}

func testAccRoleConfig_OnlyTruePermissions(resName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "%s" {
		name="%s"
		permissions {
			modify_sidecars_and_repositories = true
			modify_users = true
			modify_policies = true
			view_sidecars_and_repositories = true
			view_audit_logs = true
			modify_integrations = true
			modify_roles = true
			view_datamaps = true
		}
	}
	`, resName, updatedRoleName())
}

func testAccRoleCheck_OnlyTruePermissions(resName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "name", updatedRoleName()),
		resource.TestCheckResourceAttr(testRoleResourceFullName(resName), "permissions.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(testRoleResourceFullName(resName), "permissions.*",
			onlyTruePermissions),
	)
}
