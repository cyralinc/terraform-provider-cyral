package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSidecarInstanceStatsDataSource(t *testing.T) {
	dataSourceName := "instance_stats"
	testSteps := []resource.TestStep{
		accTestStepSidecarInstanceStatsDataSource_EmptySidecarID(dataSourceName),
		accTestStepSidecarInstanceStatsDataSource_EmptyInstanceID(dataSourceName),
		accTestStepSidecarInstanceStatsDataSource_NoSidecarFoundForGivenID(dataSourceName),
		accTestStepSidecarInstanceStatsDataSource_NoInstanceFoundForGivenID(dataSourceName),
	}
	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: providerFactories,
			Steps:             testSteps,
		},
	)
}

func accTestStepSidecarInstanceStatsDataSource_EmptySidecarID(dataSourceName string) resource.TestStep {
	config := fmt.Sprintf(`
	data "cyral_sidecar_instance_stats" "%s" {
	}
	`, dataSourceName)
	return resource.TestStep{
		Config:      config,
		ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, SidecarIDKey)),
	}
}

func accTestStepSidecarInstanceStatsDataSource_EmptyInstanceID(dataSourceName string) resource.TestStep {
	config := fmt.Sprintf(`
	data "cyral_sidecar_instance_stats" "%s" {
		sidecar_id = %q
	}
	`, dataSourceName, "some-sidecar-id")
	return resource.TestStep{
		Config:      config,
		ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, InstanceIDKey)),
	}
}

func accTestStepSidecarInstanceStatsDataSource_NoSidecarFoundForGivenID(dataSourceName string) resource.TestStep {
	nonExistentSidecarID := "some-non-existent-sidecar-id"
	config := fmt.Sprintf(`
	data "cyral_sidecar_instance_stats" "%s" {
		sidecar_id = %q
		instance_id = %q
	}
	`, dataSourceName, nonExistentSidecarID, "some-instance-id")
	return resource.TestStep{
		Config:      config,
		ExpectError: regexp.MustCompile(fmt.Sprintf("sidecar with id '%s' does not exist", nonExistentSidecarID)),
	}
}

func accTestStepSidecarInstanceStatsDataSource_NoInstanceFoundForGivenID(dataSourceName string) resource.TestStep {
	// Creates a sidecar that doesn't have any instances, since it was not
	// deployed.
	config := formatBasicSidecarIntoConfig(
		basicSidecarResName,
		accTestName("data-sidecar-instance", "sidecar"),
		"cft-ec2",
		"",
	)
	config += fmt.Sprintf(`
	data "cyral_sidecar_instance_stats" "%s" {
		sidecar_id = %s
		instance_id = %q
	}
	`, dataSourceName, basicSidecarID, "some-non-existent-instance-id")
	return resource.TestStep{
		Config:      config,
		ExpectError: regexp.MustCompile("instance not found"),
	}
}
