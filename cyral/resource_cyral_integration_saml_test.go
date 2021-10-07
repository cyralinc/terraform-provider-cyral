package cyral

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	EnvVarSAMLMetadataURL = "SAML_METADATA_URL"
)

func TestAccSAMLIntegrationResource(t *testing.T) {
	samlDisplayName := "tf-test-saml-integration"
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck: func() {
			if v := os.Getenv(EnvVarSAMLMetadataURL); v == "" {
				t.Skip(fmt.Sprintf(
					"Acceptance tests skipped unless env '%s' set", EnvVarSAMLMetadataURL))
			}
		},
		Steps: []resource.TestStep{
			/* {
				Config:      testAccSAMLIntegrationConfig_Error(samlDisplayName),
				ExpectError: regexp.MustCompile(`Error running pre-apply refresh: exit status 1`),
			}, */
			{
				Config: testAccSAMLIntegrationConfig_Okta_DefaultValues(),
				Check:  testAccSAMLIntegrationCheck_Okta_DefaultValues(),
			},
			{
				Config: testAccSAMLIntegrationConfig_GSuite_DefaultValues(),
				Check:  testAccSAMLIntegrationCheck_GSuite_DefaultValues(),
			},
			{
				Config: testAccSAMLIntegrationConfig_Updated(samlDisplayName),
				Check:  testAccSAMLIntegrationCheck_Updated(samlDisplayName),
			},
		},
	})
}

func testAccSAMLIntegrationConfig_Okta_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml_okta" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, os.Getenv(EnvVarSAMLMetadataURL))
}

func testAccSAMLIntegrationCheck_Okta_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_saml_okta.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", os.Getenv(EnvVarSAMLMetadataURL)),
	)
}

func testAccSAMLIntegrationConfig_GSuite_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml_gsuite" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, os.Getenv(EnvVarSAMLMetadataURL))
}

func testAccSAMLIntegrationCheck_GSuite_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_saml_gsuite.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", os.Getenv(EnvVarSAMLMetadataURL)),
	)
}

func testAccSAMLIntegrationConfig_Updated(samlDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml_okta" "test_saml_integration" {
		samlp {
			display_name = "%s"
			disabled = true
			config {
				single_sign_on_service_url = "%s"
				back_channel_supported = true
			}
		}
	}
	`, samlDisplayName, os.Getenv(EnvVarSAMLMetadataURL))
}

func testAccSAMLIntegrationCheck_Updated(samlDisplayName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_saml_okta.test_saml_integration",
			"samlp.0.display_name", samlDisplayName),
		resource.TestCheckResourceAttr("cyral_integration_saml_okta.test_saml_integration",
			"samlp.0.disabled", "true"),
		resource.TestCheckResourceAttr("cyral_integration_saml_okta.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", os.Getenv(EnvVarSAMLMetadataURL)),
		resource.TestCheckResourceAttr("cyral_integration_saml_okta.test_saml_integration",
			"samlp.0.config.0.back_channel_supported", "true"),
	)
}
