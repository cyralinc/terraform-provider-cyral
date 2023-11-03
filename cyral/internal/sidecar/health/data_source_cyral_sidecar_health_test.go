package health_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarHealthDataSourceFullNameFmt = "data.cyral_sidecar_health.%s"
)

func TestAccSidecarHealthDataSource(t *testing.T) {
	dataSourceName := "sidecar_health"
	testSteps := []resource.TestStep{
		accTestStepSidecarHealthDataSource_RequiredArgumentSidecarID(dataSourceName),
		accTestStepSidecarHealthDataSource_ListAllSidecarHealthInfo(dataSourceName),
	}
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps:             testSteps,
	})
}

func accTestStepSidecarHealthDataSource_RequiredArgumentSidecarID(dataSourceName string) resource.TestStep {
	config := fmt.Sprintf(`
		data "cyral_sidecar_health" "%s" {
		}
	`, dataSourceName)
	return resource.TestStep{
		Config: config,
		ExpectError: regexp.MustCompile(
			fmt.Sprintf(
				`The argument "%s" is required, but no definition was found.`,
				utils.SidecarIDKey,
			),
		),
	}
}

func accTestStepSidecarHealthDataSource_ListAllSidecarHealthInfo(dataSourceName string) resource.TestStep {
	config := utils.FormatBasicSidecarIntoConfig(
		utils.BasicSidecarResName,
		utils.AccTestName("data-sidecar-health", "sidecar"),
		"cft-ec2",
		"",
	)
	config += fmt.Sprintf(`
		data "cyral_sidecar_health" "%s" {
			sidecar_id = %s
		}
	`, dataSourceName, utils.BasicSidecarID)
	dataSourceFullName := fmt.Sprintf(sidecarHealthDataSourceFullNameFmt, dataSourceName)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			dataSourceFullName,
			utils.IDKey,
		),
		resource.TestCheckResourceAttr(
			dataSourceFullName,
			utils.StatusKey,
			"UNKNOWN",
		),
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}
