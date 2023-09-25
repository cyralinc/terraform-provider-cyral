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
			accTestStepServiceAccountResource_EmptyPermissions(serviceAccountName1),
			accTestStepServiceAccountResource_SinglePermission(serviceAccountName1),
			accTestStepServiceAccountResource_DuplicatedPermission(serviceAccountName1),
			accTestStepServiceAccountResource_AllPermissions(serviceAccountName1),
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
				`The argument "%s" is required, but no definition was found.`,
				serviceAccountResourcePermissionIDsKey,
			),
		),
	}
}

func accTestStepServiceAccountResource_EmptyPermissions(resourceName string) resource.TestStep {
	config := fmt.Sprintf(`
		resource "cyral_service_account" "%s" {
			display_name = "service-account-test"
			permission_ids = []
		}
	`,
		resourceName,
	)
	return resource.TestStep{
		Config:      config,
		ExpectError: regexp.MustCompile("at least one permission must be specified for the service account"),
	}
}

func accTestStepServiceAccountResource_SinglePermission(resourceName string) resource.TestStep {
	displayName := accTestName("service-account", "service-account-1")
	permissionNames := []string{"Modify Policies"}
	config, check := getAccTestStepForServiceAccountResourceFullConfig(
		resourceName,
		displayName,
		permissionNames,
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}

func accTestStepServiceAccountResource_DuplicatedPermission(resourceName string) resource.TestStep {
	displayName := accTestName("service-account", "service-account-1")
	permissionNames := []string{"Modify Policies", "Modify Policies"}
	config, _ := getAccTestStepForServiceAccountResourceFullConfig(
		resourceName,
		displayName,
		permissionNames,
	)
	resourceFullName := fmt.Sprintf("cyral_service_account.%s", resourceName)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			resourceFullName,
			serviceAccountResourceDisplayNameKey,
			displayName,
		),
		resource.TestCheckResourceAttr(
			resourceFullName,
			fmt.Sprintf("%s.#", serviceAccountResourcePermissionIDsKey),
			"1",
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

func accTestStepServiceAccountResource_AllPermissions(resourceName string) resource.TestStep {
	displayName := accTestName("service-account", "service-account-1")
	permissionNames := []string{
		"Approval Management",
		"Modify Policies",
		"Modify Roles",
		"Modify Sidecars and Repositories",
		"Modify Users",
		"Repo Crawler",
		"View Audit Logs",
		"View Datamaps",
		"View Integrations",
		"View Policies",
		"View Roles",
		"View Users",
		"Modify Integrations",
	}
	config, check := getAccTestStepForServiceAccountResourceFullConfig(
		resourceName,
		displayName,
		permissionNames,
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}

func accTestStepServiceAccountResource_UpdatedFields(resourceName string) resource.TestStep {
	displayName := accTestName("service-account", "service-account-1-updated")
	permissionNames := []string{
		"Approval Management",
		"Modify Roles",
		"Modify Users",
		"View Audit Logs",
		"View Integrations",
		"View Roles",
		"Modify Integrations",
	}
	config, check := getAccTestStepForServiceAccountResourceFullConfig(
		resourceName,
		displayName,
		permissionNames,
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}

func getAccTestStepForServiceAccountResourceFullConfig(
	resourceName string,
	displayName string,
	permissionNames []string,
) (string, resource.TestCheckFunc) {
	config := formatBasicDataSourcePermissionIntoConfig("permissions")
	config += fmt.Sprintf(`
		locals {
			serviceAccountPermissions = %s
		}
	`, listToStr(permissionNames),
	)
	config += fmt.Sprintf(`
		resource "cyral_service_account" "%s" {
			display_name = %q
			permission_ids = [
				for permission in data.cyral_permission.permissions.%s: permission.id 
				if contains(local.serviceAccountPermissions, permission.name)
			]
		}
	`, resourceName, displayName, PermissionDataSourcePermissionListKey,
	)
	resourceFullName := fmt.Sprintf("cyral_service_account.%s", resourceName)
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			resourceFullName,
			serviceAccountResourceDisplayNameKey,
			displayName,
		),
		resource.TestCheckResourceAttr(
			resourceFullName,
			fmt.Sprintf("%s.#", serviceAccountResourcePermissionIDsKey),
			fmt.Sprintf("%d", len(permissionNames)),
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			serviceAccountResourceClientIDKey,
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			serviceAccountResourceClientSecretKey,
		),
	}
	return config, resource.ComposeTestCheckFunc(checks...)
}
