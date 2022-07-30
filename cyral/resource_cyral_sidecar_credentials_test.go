package cyral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testSidecarCredentialsImportStateCheck(states []*terraform.InstanceState) error {
	for _, state := range states {
		id := state.ID
		attributes := state.Attributes
		if client_id, ok := attributes["client_id"]; !ok {
			return importErrorf(id, "client ID not found in state attributes")
		} else {
			if id != client_id {
				return importErrorf(id, "expected client ID to be equal to ID")
			}
		}
	}
	return nil
}

func TestAccSidecarCredentialsResource(t *testing.T) {
	testConfig, testFunc := setupSidecarCredentialsTest()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				// The check for sidecar credentials needs to be
				// manual, because ImportStateVerify will expect
				// that client_secret is set, but the GET API
				// does not return the sidecar account's client
				// secret for security reasons.
				ImportState:      true,
				ImportStateCheck: testSidecarCredentialsImportStateCheck,
				ResourceName:     "cyral_sidecar_credentials.test_sidecar_credentials",
			},
		},
	})
}

func setupSidecarCredentialsTest() (string, resource.TestCheckFunc) {
	configuration := createSidecarCredentialsConfig()

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"cyral_sidecar_credentials.test_sidecar_credentials", "sidecar_id",
			"cyral_sidecar.test_sidecar", "id"),
		resource.TestCheckResourceAttrSet(
			"cyral_sidecar_credentials.test_sidecar_credentials",
			"client_id"),
		resource.TestCheckResourceAttrSet(
			"cyral_sidecar_credentials.test_sidecar_credentials",
			"client_secret"),
		resource.TestCheckResourceAttrPair(
			"cyral_sidecar_credentials.test_sidecar_credentials", "id",
			"cyral_sidecar_credentials.test_sidecar_credentials", "client_id"),
	)

	return configuration, testFunction
}

func createSidecarCredentialsConfig() string {
	return `
	resource "cyral_sidecar" "test_sidecar" {
		name = "sidecar-test"
		deployment_method = "docker"
	}

	resource "cyral_sidecar_credentials" "test_sidecar_credentials" {
		sidecar_id = cyral_sidecar.test_sidecar.id
	}
	`
}
