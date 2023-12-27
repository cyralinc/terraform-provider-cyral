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
			ProviderFactories: provider.ProviderFactories,
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
		ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, utils.SidecarIDKey)),
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
		ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, instance.InstanceIDKey)),
	}
}

func accTestStepSidecarInstanceStatsDataSource_NoSidecarFoundForGivenID(dataSourceName string) resource.TestStep {
	nonExistentSidecarID := "id"
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
	config := utils.FormatBasicSidecarIntoConfig(
		utils.BasicSidecarResName,
		utils.AccTestName("data-sidecar-instance-stats", "sidecar"),
		"cft-ec2",
		"",
	)
	config += fmt.Sprintf(`
	data "cyral_sidecar_instance_stats" "%s" {
		sidecar_id = %s
		instance_id = %q
	}
	`, dataSourceName, utils.BasicSidecarID, "some-non-existent-instance-id")
	return resource.TestStep{
		Config:      config,
		ExpectError: regexp.MustCompile("instance not found"),
	}
}
