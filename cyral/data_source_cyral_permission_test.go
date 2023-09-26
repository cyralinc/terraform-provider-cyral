package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPermissionDataSource(t *testing.T) {
	testSteps := []resource.TestStep{}
	dataSourceName1 := "permissions_1"
	testSteps = append(
		testSteps,
		[]resource.TestStep{
			accTestStepPermissionDataSource_RetrieveAllPermissions(dataSourceName1),
		}...,
	)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps:             testSteps,
	})
}

func accTestStepPermissionDataSource_RetrieveAllPermissions(dataSourceName string) resource.TestStep {
	dataSourceFullName := fmt.Sprintf("data.cyral_permission.%s", dataSourceName)
	config := fmt.Sprintf(`
		data "cyral_permission" "%s" {
		}
	`, dataSourceName)
	var checks []resource.TestCheckFunc
	for index, expectedPermissionName := range allPermissionNames {
		checks = append(checks,
			[]resource.TestCheckFunc{
				resource.TestCheckResourceAttrSet(
					dataSourceFullName,
					fmt.Sprintf(
						"%s.%d.%s",
						PermissionDataSourcePermissionListKey,
						index,
						IDKey,
					),
				),
				resource.TestCheckTypeSetElemNestedAttrs(
					dataSourceFullName,
					fmt.Sprintf("%s.*", PermissionDataSourcePermissionListKey),
					map[string]string{NameKey: expectedPermissionName},
				),
				resource.TestCheckTypeSetElemNestedAttrs(
					dataSourceFullName,
					fmt.Sprintf("%s.*", PermissionDataSourcePermissionListKey),
					map[string]string{DescriptionKey: expectedPermissionName},
				),
			}...,
		)
	}
	return resource.TestStep{
		Config: config,
		Check:  resource.ComposeTestCheckFunc(checks...),
	}
}