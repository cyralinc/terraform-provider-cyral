package permission_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/permission"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
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
		ProviderFactories: provider.ProviderFactories,
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
	for index, expectedPermissionName := range permission.AllPermissionNames {
		checks = append(checks,
			[]resource.TestCheckFunc{
				resource.TestCheckResourceAttrSet(
					dataSourceFullName,
					fmt.Sprintf(
						"%s.%d.%s",
						permission.PermissionDataSourcePermissionListKey,
						index,
						utils.IDKey,
					),
				),
				resource.TestCheckTypeSetElemNestedAttrs(
					dataSourceFullName,
					fmt.Sprintf("%s.*", permission.PermissionDataSourcePermissionListKey),
					map[string]string{utils.NameKey: expectedPermissionName},
				),
				resource.TestCheckTypeSetElemNestedAttrs(
					dataSourceFullName,
					fmt.Sprintf("%s.*", permission.PermissionDataSourcePermissionListKey),
					map[string]string{utils.DescriptionKey: expectedPermissionName},
				),
			}...,
		)
	}
	return resource.TestStep{
		Config: config,
		Check:  resource.ComposeTestCheckFunc(checks...),
	}
}
