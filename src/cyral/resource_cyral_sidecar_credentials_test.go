package cyral

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarCredentialsResourceName = "sidecar-credentials"
)

func TestAccSidecarCredentialsResource(t *testing.T) {
	testConfig, testFunc := setupSidecarCredentialsTest()

	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: testConfig,
					Check:  testFunc,
				},
				{
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"client_secret"},
					ResourceName:            "cyral_sidecar_credentials.test_sidecar_credentials",
				},
			},
		},
	)
}

func setupSidecarCredentialsTest() (string, resource.TestCheckFunc) {
	configuration := createSidecarCredentialsConfig()

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"cyral_sidecar_credentials.test_sidecar_credentials", "sidecar_id",
			"cyral_sidecar.test_sidecar", "id",
		),
		resource.TestCheckResourceAttrSet(
			"cyral_sidecar_credentials.test_sidecar_credentials",
			"client_id",
		),
		resource.TestCheckResourceAttrSet(
			"cyral_sidecar_credentials.test_sidecar_credentials",
			"client_secret",
		),
		resource.TestCheckResourceAttrPair(
			"cyral_sidecar_credentials.test_sidecar_credentials", "id",
			"cyral_sidecar_credentials.test_sidecar_credentials", "client_id",
		),
	)

	return configuration, testFunction
}

func createSidecarCredentialsConfig() string {
	var config string
	config += utils.FormatBasicSidecarIntoConfig(
		utils.BasicSidecarResName,
		utils.AccTestName(sidecarCredentialsResourceName, "sidecar"),
		"docker", "",
	)
	config += fmt.Sprintf(
		`
	resource "cyral_sidecar_credentials" "test_sidecar_credentials" {
		sidecar_id = %s
	}`, utils.BasicSidecarID,
	)
	return config
}
