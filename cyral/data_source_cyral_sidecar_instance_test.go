package cyral

import (
	"fmt"
	"regexp"
	"testing"

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
			ProviderFactories: providerFactories,
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
		ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, SidecarIDKey)),
	}
}

func accTestStepSidecarInstanceDataSource_NoSidecarFoundForGivenID(dataSourceName string) resource.TestStep {
	nonExistentSidecarID := "some-non-existent-sidecar-id"
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
	config := formatBasicSidecarIntoConfig(
		basicSidecarResName,
		accTestName("data-sidecar-instance", "sidecar"),
		"cft-ec2",
		"",
	)
	config += fmt.Sprintf(`
	data "cyral_sidecar_instance" "%s" {
		sidecar_id = %s
	}
	`, dataSourceName, basicSidecarID)
	dataSourceFullName := fmt.Sprintf(sidecarInstanceDataSourceFullNameFmt, dataSourceName)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			dataSourceFullName,
			SidecarIDKey,
		),
		resource.TestCheckResourceAttrSet(
			dataSourceFullName,
			IDKey,
		),
		resource.TestCheckResourceAttr(
			dataSourceFullName,
			fmt.Sprintf("%s.#", SidecarInstanceListKey),
			"0",
		),
	)
	return resource.TestStep{
		Config: config,
		Check:  check,
	}
}
