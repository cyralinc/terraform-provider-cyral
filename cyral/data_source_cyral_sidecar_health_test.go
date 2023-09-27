package cyral

import (
	"fmt"
	"regexp"
	"testing"

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
		ProviderFactories: providerFactories,
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
				SidecarIDKey,
			),
		),
	}
}

func accTestStepSidecarHealthDataSource_ListAllSidecarHealthInfo(dataSourceName string) resource.TestStep {
	config := formatBasicSidecarIntoConfig(
		basicSidecarResName,
		accTestName("data-sidecar-health", "sidecar"),
		"cft-ec2",
		"",
	)
	config += fmt.Sprintf(`
		data "cyral_sidecar_health" "%s" {
			sidecar_id = %s
		}
	`, dataSourceName, basicSidecarID)
	dataSourceFullName := fmt.Sprintf(sidecarHealthDataSourceFullNameFmt, dataSourceName)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			dataSourceFullName,
			IDKey,
		),
		resource.TestCheckResourceAttr(
			dataSourceFullName,
			StatusKey,
			"UNKNOWN",
		),
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}
