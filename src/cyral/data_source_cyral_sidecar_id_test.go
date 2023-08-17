package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarIDDataSourceName = "data-sidecar-id"
)

func TestAccSidecarIDDataSource(t *testing.T) {
	nonExistentSidecarName := "some-non-existent-sidecar-name"

	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config:      testAccSidecarIDConfig_EmptySidecarName(),
					ExpectError: regexp.MustCompile(`The argument "sidecar_name" is required`),
				},
				{
					Config: testAccSidecarIDConfig_NoSidecarFoundForGivenName(nonExistentSidecarName),
					ExpectError: regexp.MustCompile(
						fmt.Sprintf(
							"No sidecar found for name '%s'.",
							nonExistentSidecarName,
						),
					),
				},
				{
					Config: testAccSidecarIDConfig_ExistentSidecar(),
					Check:  testAccSidecarIDCheck_ExistentSidecar(),
				},
			},
		},
	)
}

func testAccSidecarIDConfig_EmptySidecarName() string {
	return `
	data "cyral_sidecar_id" "sidecar_id" {
	}
	`
}

func testAccSidecarIDConfig_NoSidecarFoundForGivenName(nonExistentSidecarName string) string {
	return fmt.Sprintf(
		`
	data "cyral_sidecar_id" "sidecar_id" {
		sidecar_name = "%s"
	}
	`, nonExistentSidecarName,
	)
}

func testAccSidecarIDConfig_ExistentSidecar() string {
	var config string
	config += utils.FormatBasicSidecarIntoConfig(
		BasicSidecarResName,
		utils.AccTestName(sidecarIDDataSourceName, "sidecar"),
		"cloudFormation", "",
	)
	config += fmt.Sprintf(
		`
	data "cyral_sidecar_id" "sidecar_id" {
		sidecar_name = cyral_sidecar.%s.name
	}`, BasicSidecarResName,
	)
	return config
}

func testAccSidecarIDCheck_ExistentSidecar() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			fmt.Sprintf("cyral_sidecar.%s", BasicSidecarResName), "id",
			"data.cyral_sidecar_id.sidecar_id", "id",
		),
	)
}
