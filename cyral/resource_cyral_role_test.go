package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const InitialRoleName = "role-test"
const UpdatedRoleName = "updated-role-test"

var onlyFalsePermissions = map[string]string{
	"view_sidecars_and_repositories":   "false",
	"modify_sidecars_and_repositories": "false",
	"modify_policies":                  "false",
	"modify_users":                     "false",
	"modify_roles":                     "false",
	"view_audit_logs":                  "false",
	"modify_integrations":              "false",
}

var trueAndFalsePermissions = map[string]string{
	"view_sidecars_and_repositories":   "true",
	"modify_sidecars_and_repositories": "true",
	"modify_policies":                  "false",
	"modify_users":                     "false",
	"modify_roles":                     "false",
	"view_audit_logs":                  "false",
	"modify_integrations":              "false",
}

var onlyTruePermissions = map[string]string{
	"view_sidecars_and_repositories":   "true",
	"modify_sidecars_and_repositories": "true",
	"modify_policies":                  "true",
	"modify_users":                     "true",
	"modify_roles":                     "true",
	"view_audit_logs":                  "true",
	"modify_integrations":              "true",
}

func TestAccRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: initialRoleConfigWithoutPermissions(InitialRoleName),
				Check:  initialRoleCheck(InitialRoleName),
			},
			{
				Config: updatedRoleConfigEmptyPermissions(UpdatedRoleName),
				Check:  updatedRoleCheck(UpdatedRoleName, onlyFalsePermissions),
			},
			{
				Config: updatedRoleConfigOnlyFalsePermissions(UpdatedRoleName),
				Check:  updatedRoleCheck(UpdatedRoleName, onlyFalsePermissions),
			},
			{
				Config: updatedRoleConfigTrueAndFalsePermissions(UpdatedRoleName),
				Check:  updatedRoleCheck(UpdatedRoleName, trueAndFalsePermissions),
			},
			{
				Config: updatedRoleConfigOnlyTruePermissions(UpdatedRoleName),
				Check:  updatedRoleCheck(UpdatedRoleName, onlyTruePermissions),
			},
		},
	})
}

func initialRoleConfigWithoutPermissions(roleName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
	}
	`, roleName)
}

func updatedRoleConfigEmptyPermissions(roleName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
		permissions {
		}
	}
	`, roleName)
}

func updatedRoleConfigOnlyFalsePermissions(roleName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
		permissions {
			view_sidecars_and_repositories = false
			modify_sidecars_and_repositories = false
			modify_policies = false
			modify_users = false
			modify_roles = false
			view_audit_logs = false   
			modify_integrations = false 
		}
	}
	`, roleName)
}

func updatedRoleConfigTrueAndFalsePermissions(roleName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
		permissions {
			view_sidecars_and_repositories = true
			modify_sidecars_and_repositories = true
			modify_policies = false
			modify_users = false
			modify_roles = false
			view_audit_logs = false   
			modify_integrations = false 
		}
	}
	`, roleName)
}

func updatedRoleConfigOnlyTruePermissions(roleName string) string {
	return fmt.Sprintf(`
	resource "cyral_role" "test_role" {
		name="%s"
		permissions {
			view_sidecars_and_repositories = true
			modify_sidecars_and_repositories = true
			modify_policies = true
			modify_users = true
			modify_roles = true
			view_audit_logs = true   
			modify_integrations = true 
		}
	}
	`, roleName)
}

func initialRoleCheck(roleName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_role.test_role", "name", roleName),
		resource.TestCheckResourceAttr("cyral_role.test_role", "permissions.#", "0"),
	)
}

func updatedRoleCheck(roleName string, permissions map[string]string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_role.test_role", "name", roleName),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_role.test_role", "permissions.*",
			permissions),
	)
}
