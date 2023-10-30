package internal_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServiceAccountResource(t *testing.T) {
	testSteps := []resource.TestStep{}
	resourceName1 := "service_account_1"
	testSteps = append(
		testSteps,
		[]resource.TestStep{
			accTestStepServiceAccountResource_RequiredArgumentDisplayName(resourceName1),
			accTestStepServiceAccountResource_RequiredArgumentPermissions(resourceName1),
			accTestStepServiceAccountResource_EmptyPermissions(resourceName1),
			accTestStepServiceAccountResource_SinglePermission(resourceName1),
			accTestStepServiceAccountResource_DuplicatedPermission(resourceName1),
			accTestStepServiceAccountResource_AllPermissions(resourceName1),
			accTestStepServiceAccountResource_UpdatedFields(resourceName1),
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
				internal.ServiceAccountResourceDisplayNameKey,
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
				internal.ServiceAccountResourcePermissionIDsKey,
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
	displayName := utils.AccTestName("service-account", "service-account-1")
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
	displayName := utils.AccTestName("service-account", "service-account-1")
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
			internal.ServiceAccountResourceDisplayNameKey,
			displayName,
		),
		resource.TestCheckResourceAttr(
			resourceFullName,
			fmt.Sprintf("%s.#", internal.ServiceAccountResourcePermissionIDsKey),
			"1",
		),
		resource.TestCheckResourceAttrPair(
			resourceFullName,
			utils.IDKey,
			resourceFullName,
			internal.ServiceAccountResourceClientIDKey,
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			internal.ServiceAccountResourceClientIDKey,
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			internal.ServiceAccountResourceClientSecretKey,
		),
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}

func accTestStepServiceAccountResource_AllPermissions(resourceName string) resource.TestStep {
	displayName := utils.AccTestName("service-account", "service-account-1")
	config, check := getAccTestStepForServiceAccountResourceFullConfig(
		resourceName,
		displayName,
		internal.AllPermissionNames,
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}

func accTestStepServiceAccountResource_UpdatedFields(resourceName string) resource.TestStep {
	displayName := utils.AccTestName("service-account", "service-account-1-updated")
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
	config := utils.FormatBasicDataSourcePermissionIntoConfig("permissions")
	config += fmt.Sprintf(`
		locals {
			serviceAccountPermissions = %s
		}
	`, utils.ListToStr(permissionNames),
	)
	config += fmt.Sprintf(`
		resource "cyral_service_account" "%s" {
			display_name = %q
			permission_ids = [
				for permission in data.cyral_permission.permissions.%s: permission.id
				if contains(local.serviceAccountPermissions, permission.name)
			]
		}
	`, resourceName, displayName, internal.PermissionDataSourcePermissionListKey,
	)
	resourceFullName := fmt.Sprintf("cyral_service_account.%s", resourceName)
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			resourceFullName,
			internal.ServiceAccountResourceDisplayNameKey,
			displayName,
		),
		resource.TestCheckResourceAttr(
			resourceFullName,
			fmt.Sprintf("%s.#", internal.ServiceAccountResourcePermissionIDsKey),
			fmt.Sprintf("%d", len(permissionNames)),
		),
		resource.TestCheckResourceAttrPair(
			resourceFullName,
			utils.IDKey,
			resourceFullName,
			internal.ServiceAccountResourceClientIDKey,
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			internal.ServiceAccountResourceClientIDKey,
		),
		resource.TestCheckResourceAttrSet(
			resourceFullName,
			internal.ServiceAccountResourceClientSecretKey,
		),
	}
	return config, resource.ComposeTestCheckFunc(checks...)
}
