package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSidecarIDDataSource(t *testing.T) {
	nonExistentSidecarName := "some-non-existent-sidecar-name"

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSidecarIDConfig_EmptySidecarName(),
				ExpectError: regexp.MustCompile(`The argument "sidecar_name" is required`),
			},
			{
				Config: testAccSidecarIDConfig_NoSidecarFoundForGivenName(nonExistentSidecarName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("No sidecar found for name '%s'.",
					nonExistentSidecarName)),
			},
			{
				Config: testAccSidecarIDConfig_ExistentSidecar(),
				Check:  testAccSidecarIDCheck_ExistentSidecar(),
			},
		},
	})
}

func testAccSidecarIDConfig_EmptySidecarName() string {
	return `
	data "cyral_sidecar_id" "sidecar_id" {
	}
	`
}

func testAccSidecarIDConfig_NoSidecarFoundForGivenName(nonExistentSidecarName string) string {
	return fmt.Sprintf(`
	data "cyral_sidecar_id" "sidecar_id" {
		sidecar_name = "%s"
	}
	`, nonExistentSidecarName)
}

func testAccSidecarIDConfig_ExistentSidecar() string {
	return `
	resource "cyral_sidecar" "sidecar_1" {
		name = "tf-provider-sidecar-1"
		deployment_method = "cloudFormation"
		labels = ["terraform-provider", "sidecar-id"]
	}

	data "cyral_sidecar_id" "sidecar_id" {
		sidecar_name = cyral_sidecar.sidecar_1.name
	}
	`
}

func testAccSidecarIDCheck_ExistentSidecar() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"cyral_sidecar.sidecar_1", "id",
			"data.cyral_sidecar_id.sidecar_id", "id",
		),
	)
}
