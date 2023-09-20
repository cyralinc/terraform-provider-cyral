package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServiceAccountResource(t *testing.T) {
	testSteps := []resource.TestStep{}
	serviceAccountName1 := accTestName("service-account", "1")
	testSteps = append(
		testSteps,
		[]resource.TestStep{
			accTestStepServiceAccountResource_RequiredArgumentDisplayName(serviceAccountName1),
			accTestStepServiceAccountResource_RequiredArgumentPermissions(serviceAccountName1),
			accTestStepServiceAccountResource_AllPermissionsFalse(serviceAccountName1),
			accTestStepServiceAccountResource_SinglePermissionTrue(serviceAccountName1),
			accTestStepServiceAccountResource_AllPermissionsTrue(serviceAccountName1),
			accTestStepServiceAccountResource_UpdatedFields(serviceAccountName1),
		}...,
	)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps:             testSteps,
	})
}

func accTestStepServiceAccountResource_RequiredArgumentDisplayName(resourceName string) resource.TestStep {
	config := fmt.Sprintf(`
		resource "cyral_service_account" "%s" {
		}
	`, resourceName)
	return resource.TestStep{
		Config: config,
		ExpectError: regexp.MustCompile(
			fmt.Sprintf(
				`The argument "%s" is required, but no definition was found.`,
				serviceAccountResourceDisplayNameKey,
			),
		),
	}
}

func accTestStepServiceAccountResource_RequiredArgumentPermissions(resourceName string) resource.TestStep {
	config := fmt.Sprintf(`
		resource "cyral_service_account" "%s" {
			display_name = "service-account-test"
		}
	`, resourceName,
	)
	return resource.TestStep{
		Config: config,
		ExpectError: regexp.MustCompile(
			fmt.Sprintf(
				`At least 1 "%s" blocks are required.`,
				serviceAccountResourcePermissionsKey,
			),
		),
	}
}

func accTestStepServiceAccountResource_AllPermissionsFalse(resourceName string) resource.TestStep {
	config := fmt.Sprintf(`
		resource "cyral_service_account" "%s" {
			display_name = "service-account-test"
			permissions {}
		}
	`,
		resourceName,
	)
	return resource.TestStep{
		Config:      config,
		ExpectError: regexp.MustCompile("at least one permission must be specified for the service account"),
	}
}

func accTestStepServiceAccountResource_SinglePermissionTrue(resourceName string) resource.TestStep {
	resourceFullName := fmt.Sprintf("cyral_service_account.%s", resourceName)
	config := fmt.Sprintf(`
		resource "cyral_service_account" "%s" {
			display_name = "service-account-test"
			permissions {
				modify_sidecars_and_repositories = true
			}
		}
	`, resourceName,
	)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			resourceFullName,
			serviceAccountResourceDisplayNameKey, "service-account-test",
		),
		resource.TestCheckResourceAttr(
			resourceFullName,
			fmt.Sprintf("%s.#", serviceAccountResourcePermissionsKey), "1",
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			resourceFullName,
			fmt.Sprintf("%s.*", serviceAccountResourcePermissionsKey), map[string]string{
				modifySidecarAndRepositoriesPermissionKey: "true",
				modifyPoliciesPermissionKey:               "false",
				modifyIntegrationsPermissionKey:           "false",
				modifyUsersPermissionKey:                  "false",
				modifyRolesPermissionKey:                  "false",
				viewUsersPermissionKey:                    "false",
				viewAuditLogsPermissionKey:                "false",
				repoCrawlerPermissionKey:                  "false",
				viewDatamapsPermissionKey:                 "false",
				viewRolesPermissionKey:                    "false",
				viewPoliciesPermissionKey:                 "false",
				approvalManagementPermissionKey:           "false",
				viewIntegrationsPermissionKey:             "false",
			},
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			serviceAccountResourceClientIDKey,
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			serviceAccountResourceClientSecretKey,
		),
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}

func accTestStepServiceAccountResource_AllPermissionsTrue(resourceName string) resource.TestStep {
	resourceFullName := fmt.Sprintf("cyral_service_account.%s", resourceName)
	config := fmt.Sprintf(`
		resource "cyral_service_account" "%s" {
			display_name = "service-account-test"
			permissions {
				modify_sidecars_and_repositories = true
				modify_policies = true
				modify_integrations = true
				modify_users = true
				modify_roles = true
				view_users = true
				view_audit_logs = true
				repo_crawler = true
				view_datamaps = true
				view_roles = true
				view_policies = true
				approval_management = true
				view_integrations = true
			}
		}
	`, resourceName,
	)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			resourceFullName,
			serviceAccountResourceDisplayNameKey, "service-account-test",
		),
		resource.TestCheckResourceAttr(
			resourceFullName,
			fmt.Sprintf("%s.#", serviceAccountResourcePermissionsKey), "1",
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			resourceFullName,
			fmt.Sprintf("%s.*", serviceAccountResourcePermissionsKey), map[string]string{
				modifySidecarAndRepositoriesPermissionKey: "true",
				modifyPoliciesPermissionKey:               "true",
				modifyIntegrationsPermissionKey:           "true",
				modifyUsersPermissionKey:                  "true",
				modifyRolesPermissionKey:                  "true",
				viewUsersPermissionKey:                    "true",
				viewAuditLogsPermissionKey:                "true",
				repoCrawlerPermissionKey:                  "true",
				viewDatamapsPermissionKey:                 "true",
				viewRolesPermissionKey:                    "true",
				viewPoliciesPermissionKey:                 "true",
				approvalManagementPermissionKey:           "true",
				viewIntegrationsPermissionKey:             "true",
			},
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			serviceAccountResourceClientIDKey,
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			serviceAccountResourceClientSecretKey,
		),
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}

func accTestStepServiceAccountResource_UpdatedFields(resourceName string) resource.TestStep {
	resourceFullName := fmt.Sprintf("cyral_service_account.%s", resourceName)
	config := fmt.Sprintf(`
		resource "cyral_service_account" "%s" {
			display_name = "service-account-test-updated"
			permissions {
				modify_sidecars_and_repositories = true
				modify_policies = false
				modify_integrations = true
				modify_users = false
				modify_roles = true
				view_users = false
				view_audit_logs = true
				repo_crawler = false
				view_datamaps = true
				view_roles = false
				view_policies = true
				approval_management = false
				view_integrations = true
			}
		}
	`, resourceName,
	)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			resourceFullName,
			serviceAccountResourceDisplayNameKey, "service-account-test-updated",
		),
		resource.TestCheckResourceAttr(
			resourceFullName,
			fmt.Sprintf("%s.#", serviceAccountResourcePermissionsKey), "1",
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			resourceFullName,
			fmt.Sprintf("%s.*", serviceAccountResourcePermissionsKey), map[string]string{
				modifySidecarAndRepositoriesPermissionKey: "true",
				modifyPoliciesPermissionKey:               "false",
				modifyIntegrationsPermissionKey:           "true",
				modifyUsersPermissionKey:                  "false",
				modifyRolesPermissionKey:                  "true",
				viewUsersPermissionKey:                    "false",
				viewAuditLogsPermissionKey:                "true",
				repoCrawlerPermissionKey:                  "false",
				viewDatamapsPermissionKey:                 "true",
				viewRolesPermissionKey:                    "false",
				viewPoliciesPermissionKey:                 "true",
				approvalManagementPermissionKey:           "false",
				viewIntegrationsPermissionKey:             "true",
			},
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			serviceAccountResourceClientIDKey,
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			serviceAccountResourceClientSecretKey,
		),
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}
