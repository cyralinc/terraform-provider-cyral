package instance_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/instance"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarInstanceDataSourceFullNameFmt = "data.cyral_sidecar_instance.%s"
)

func TestAccSidecarInstanceDataSource(t *testing.T) {
	dataSourceName := "instances"
	testSteps := []resource.TestStep{
		accTestStepSidecarInstanceDataSource_EmptySidecarID(dataSourceName),
		accTestStepSidecarInstanceDataSource_NoSidecarFoundForGivenID(dataSourceName),
		accTestStepSidecarInstanceDataSource_NoSidecarInstances(dataSourceName),
	}
	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: provider.ProviderFactories,
			Steps:             testSteps,
		},
	)
}

func accTestStepSidecarInstanceDataSource_EmptySidecarID(dataSourceName string) resource.TestStep {
	config := fmt.Sprintf(`
	data "cyral_sidecar_instance" "%s" {
	}
	`, dataSourceName)
	return resource.TestStep{
		Config:      config,
		ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, utils.SidecarIDKey)),
	}
}

func accTestStepSidecarInstanceDataSource_NoSidecarFoundForGivenID(dataSourceName string) resource.TestStep {
	nonExistentSidecarID := "id"
	config := fmt.Sprintf(`
	data "cyral_sidecar_instance" "%s" {
		sidecar_id = %q
	}
	`, dataSourceName, nonExistentSidecarID)
	return resource.TestStep{
		Config:      config,
		ExpectError: regexp.MustCompile(fmt.Sprintf("sidecar with id '%s' does not exist", nonExistentSidecarID)),
	}
}

func accTestStepSidecarInstanceDataSource_NoSidecarInstances(dataSourceName string) resource.TestStep {
	// Creates a sidecar that doesn't have any instances, since it was not
	// deployed.
	config := utils.FormatBasicSidecarIntoConfig(
		utils.BasicSidecarResName,
		utils.AccTestName("data-sidecar-instance", "sidecar"),
		"cft-ec2",
		"",
	)
	config += fmt.Sprintf(`
	data "cyral_sidecar_instance" "%s" {
		sidecar_id = %s
	}
	`, dataSourceName, utils.BasicSidecarID)
	dataSourceFullName := fmt.Sprintf(sidecarInstanceDataSourceFullNameFmt, dataSourceName)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			dataSourceFullName,
			utils.SidecarIDKey,
		),
		resource.TestCheckResourceAttrSet(
			dataSourceFullName,
			utils.IDKey,
		),
		resource.TestCheckResourceAttr(
			dataSourceFullName,
			fmt.Sprintf("%s.#", instance.SidecarInstanceListKey),
			"0",
		),
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}
