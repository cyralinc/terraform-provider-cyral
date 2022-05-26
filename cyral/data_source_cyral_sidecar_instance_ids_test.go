package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSidecarInstanceIDsDataSource(t *testing.T) {
	nonExistentSidecarID := "some-non-existent-sidecar-id"

	resource.Test(t, resource.TestCase{
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
	return `
	// Creates a sidecar that doesn't have any instances,
	// since it was not deployed.
	resource "cyral_sidecar" "sidecar_1" {
		name = "tf-provider-sidecar-1"
		deployment_method = "cloudFormation"
		labels = ["terraform-provider", "sidecar-instance-ids"]
	}
	data "cyral_sidecar_instance_ids" "instance_ids" {
		sidecar_id = cyral_sidecar.sidecar_1.id
	}
	`
}

func testAccSidecarInstanceIDsCheck_NoSidecarInstances() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_instance_ids.instance_ids",
			"instance_ids.#", "0",
		),
	)
}
