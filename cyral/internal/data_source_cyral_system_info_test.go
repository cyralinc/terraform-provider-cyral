package internal_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	systemInfoDataSourceFullNameFmt = "data.cyral_system_info.%s"
)

func TestAccSystemInfoDataSource(t *testing.T) {
	dataSourceName := "system_info"
	testSteps := []resource.TestStep{
		accTestStepSystemInfoDataSource_ListAllSystemInfo(dataSourceName),
	}
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps:             testSteps,
	})
}

func accTestStepSystemInfoDataSource_ListAllSystemInfo(dataSourceName string) resource.TestStep {
	dataSourceFullName := fmt.Sprintf(systemInfoDataSourceFullNameFmt, dataSourceName)
	config := fmt.Sprintf(`
		data "cyral_system_info" "%s" {
		}
	`, dataSourceName)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			dataSourceFullName,
			utils.IDKey,
		),
		resource.TestCheckResourceAttrSet(
			dataSourceFullName,
			internal.ControlPlaneVersionKey,
		),
		resource.TestCheckResourceAttrSet(
			dataSourceFullName,
			internal.SidecarLatestVersionKey,
		),
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}
