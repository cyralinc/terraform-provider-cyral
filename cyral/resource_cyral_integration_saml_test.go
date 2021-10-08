package cyral

import (
	"fmt"
	"os"
	"regexp"
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
			{
				Config:      testAccSAMLIntegrationConfig_EmptySamlp(),
				ExpectError: regexp.MustCompile(`At least 1 "samlp" blocks are required`),
			},
			{
				Config:      testAccSAMLIntegrationConfig_EmptyConfig(),
				ExpectError: regexp.MustCompile(`At least 1 "config" blocks are required`),
			},
			{
				Config: testAccSAMLIntegrationConfig_EmptySSOUrl(),
				ExpectError: regexp.MustCompile(
					`The argument "single_sign_on_service_url" is required`),
			},
			{
				Config: testAccSAMLIntegrationConfig_ADFS_DefaultValues(),
				Check:  testAccSAMLIntegrationCheck_ADFS_DefaultValues(),
			},
			{
				Config: testAccSAMLIntegrationConfig_AAD_DefaultValues(),
				Check:  testAccSAMLIntegrationCheck_AAD_DefaultValues(),
			},
			{
				Config: testAccSAMLIntegrationConfig_Forgerock_DefaultValues(),
				Check:  testAccSAMLIntegrationCheck_Forgerock_DefaultValues(),
			},
			{
				Config: testAccSAMLIntegrationConfig_GSuite_DefaultValues(),
				Check:  testAccSAMLIntegrationCheck_GSuite_DefaultValues(),
			},
			{
				Config: testAccSAMLIntegrationConfig_Okta_DefaultValues(),
				Check:  testAccSAMLIntegrationCheck_Okta_DefaultValues(),
			},
			{
				Config: testAccSAMLIntegrationConfig_Pingone_DefaultValues(),
				Check:  testAccSAMLIntegrationCheck_Pingone_DefaultValues(),
			},
			{
				Config: testAccSAMLIntegrationConfig_NotEmptyAlias(),
				Check:  testAccSAMLIntegrationCheck_NotEmptyAlias(),
			},
			{
				Config: testAccSAMLIntegrationConfig_Updated(samlDisplayName),
				Check:  testAccSAMLIntegrationCheck_Updated(samlDisplayName),
			},
		},
	})
}

func testAccSAMLIntegrationConfig_EmptySamlp() string {
	return `
	resource "cyral_integration_saml_okta" "test_saml_integration" {
	}
	`
}

func testAccSAMLIntegrationConfig_EmptyConfig() string {
	return `
	resource "cyral_integration_saml_okta" "test_saml_integration" {
		samlp {
		}
	}
	`
}

func testAccSAMLIntegrationConfig_EmptySSOUrl() string {
	return `
	resource "cyral_integration_saml_okta" "test_saml_integration" {
		samlp {
			config {
			}
		}
	}
	`
}

func testAccSAMLIntegrationConfig_ADFS_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml_adfs" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, os.Getenv(EnvVarSAMLMetadataURL))
}

func testAccSAMLIntegrationCheck_ADFS_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_saml_adfs.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", os.Getenv(EnvVarSAMLMetadataURL)),
	)
}

func testAccSAMLIntegrationConfig_AAD_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml_aad" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, os.Getenv(EnvVarSAMLMetadataURL))
}

func testAccSAMLIntegrationCheck_AAD_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_saml_aad.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", os.Getenv(EnvVarSAMLMetadataURL)),
	)
}

func testAccSAMLIntegrationConfig_Forgerock_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml_forgerock" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, os.Getenv(EnvVarSAMLMetadataURL))
}

func testAccSAMLIntegrationCheck_Forgerock_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_saml_forgerock.test_saml_integration",
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

func testAccSAMLIntegrationConfig_Pingone_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml_pingone" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, os.Getenv(EnvVarSAMLMetadataURL))
}

func testAccSAMLIntegrationCheck_Pingone_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_saml_pingone.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", os.Getenv(EnvVarSAMLMetadataURL)),
	)
}

func testAccSAMLIntegrationConfig_NotEmptyAlias() string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml_okta" "test_saml_integration" {
		draft_alias = "test-alias"
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, os.Getenv(EnvVarSAMLMetadataURL))
}

func testAccSAMLIntegrationCheck_NotEmptyAlias() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_saml_okta.test_saml_integration",
			"draft_alias", "test-alias"),
		resource.TestCheckResourceAttr("cyral_integration_saml_okta.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", os.Getenv(EnvVarSAMLMetadataURL)),
	)
}

func testAccSAMLIntegrationConfig_Updated(samlDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml_okta" "test_saml_integration" {
		draft_alias = "test-alias"
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
			"draft_alias", "test-alias"),
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
