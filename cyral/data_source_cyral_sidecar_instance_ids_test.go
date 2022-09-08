package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarInstanceIDsDataSourceName = "data-sidecar-instance-ids"
)

func TestAccSidecarInstanceIDsDataSource(t *testing.T) {
	nonExistentSidecarID := "some-non-existent-sidecar-id"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSidecarInstanceIDsConfig_EmptySidecarID(),
				ExpectError: regexp.MustCompile(`The argument "sidecar_id" is required`),
			},
			{
				Config: testAccSidecarInstanceIDsConfig_NoSidecarFoundForGivenID(
					nonExistentSidecarID,
				),
				ExpectError: regexp.MustCompile(fmt.Sprintf(
					"Unable to retrieve sidecar details. SidecarID: %s",
					nonExistentSidecarID,
				)),
			},
			{
				Config: testAccSidecarInstanceIDsConfig_NoSidecarInstances(),
				Check:  testAccSidecarInstanceIDsCheck_NoSidecarInstances(),
			},
		},
	})
}

func testAccSidecarInstanceIDsConfig_EmptySidecarID() string {
	return `
	data "cyral_sidecar_instance_ids" "instance_ids" {
	}
	`
}

func testAccSidecarInstanceIDsConfig_NoSidecarFoundForGivenID(
	nonExistentSidecarID string,
) string {
	return fmt.Sprintf(`
	data "cyral_sidecar_instance_ids" "instance_ids" {
		sidecar_id = "%s"
	}
	`, nonExistentSidecarID)
}

func testAccSidecarInstanceIDsConfig_NoSidecarInstances() string {
	// Creates a sidecar that doesn't have any instances, since it was not
	// deployed.
	var config string
	config += formatBasicSidecarIntoConfig(
		basicSidecarResName,
		accTestName(sidecarInstanceIDsDataSourceName, "sidecar"),
		"cloudFormation",
	)

	config += fmt.Sprintf(`
	data "cyral_sidecar_instance_ids" "instance_ids" {
		sidecar_id = %s
	}`, basicSidecarID)
	return config
}

func testAccSidecarInstanceIDsCheck_NoSidecarInstances() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_instance_ids.instance_ids",
			"instance_ids.#", "0",
		),
	)
}
