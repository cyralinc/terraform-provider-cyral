package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testSingleSignOnURL = "https://some-test-sso-url.com"
)

func TestAccIDPIntegrationResource(t *testing.T) {
	samlDisplayName := "tf-test-saml-integration"
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccIDPIntegrationConfig_EmptySamlp(),
				ExpectError: regexp.MustCompile(`At least 1 "samlp" blocks are required`),
			},
			{
				Config:      testAccIDPIntegrationConfig_EmptyConfig(),
				ExpectError: regexp.MustCompile(`At least 1 "config" blocks are required`),
			},
			{
				Config: testAccIDPIntegrationConfig_EmptySSOUrl(),
				ExpectError: regexp.MustCompile(
					`The argument "single_sign_on_service_url" is required`),
			},
			{
				Config: testAccIDPIntegrationConfig_ADFS_DefaultValues(),
				Check:  testAccIDPIntegrationCheck_ADFS_DefaultValues(),
			},
			{
				Config: testAccIDPIntegrationConfig_AAD_DefaultValues(),
				Check:  testAccIDPIntegrationCheck_AAD_DefaultValues(),
			},
			{
				Config: testAccIDPIntegrationConfig_Forgerock_DefaultValues(),
				Check:  testAccIDPIntegrationCheck_Forgerock_DefaultValues(),
			},
			{
				Config: testAccIDPIntegrationConfig_GSuite_DefaultValues(),
				Check:  testAccIDPIntegrationCheck_GSuite_DefaultValues(),
			},
			{
				Config: testAccIDPIntegrationConfig_PingOne_DefaultValues(),
				Check:  testAccIDPIntegrationCheck_PingOne_DefaultValues(),
			},
			{
				Config: testAccIDPIntegrationConfig_Okta_DefaultValues(),
				Check:  testAccIDPIntegrationCheck_Okta_DefaultValues(),
			},
			{
				Config: testAccIDPIntegrationConfig_Updated(samlDisplayName),
				Check:  testAccIDPIntegrationCheck_Updated(samlDisplayName),
			},
			{
				Config: testAccIDPIntegrationConfig_NotEmptyAlias(),
				Check:  testAccIDPIntegrationCheck_NotEmptyAlias(),
			},
		},
	})
}

func testAccIDPIntegrationConfig_EmptySamlp() string {
	return `
	resource "cyral_integration_idp_okta" "test_saml_integration" {
	}
	`
}

func testAccIDPIntegrationConfig_EmptyConfig() string {
	return `
	resource "cyral_integration_idp_okta" "test_saml_integration" {
		samlp {
		}
	}
	`
}

func testAccIDPIntegrationConfig_EmptySSOUrl() string {
	return `
	resource "cyral_integration_idp_okta" "test_saml_integration" {
		samlp {
			config {
			}
		}
	}
	`
}

func testAccIDPIntegrationConfig_ADFS_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_adfs" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIDPIntegrationCheck_ADFS_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_adfs.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIDPIntegrationConfig_AAD_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_aad" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIDPIntegrationCheck_AAD_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_aad.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIDPIntegrationConfig_Forgerock_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_forgerock" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIDPIntegrationCheck_Forgerock_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_forgerock.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIDPIntegrationConfig_GSuite_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_gsuite" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIDPIntegrationCheck_GSuite_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_gsuite.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIDPIntegrationConfig_PingOne_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_ping_one" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIDPIntegrationCheck_PingOne_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_ping_one.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIDPIntegrationConfig_Okta_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIDPIntegrationCheck_Okta_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIDPIntegrationConfig_Updated(samlDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "test_saml_integration" {
		samlp {
			display_name = "%s"
			disabled = true
			config {
				single_sign_on_service_url = "%s"
				back_channel_supported = true
			}
		}
	}
	`, samlDisplayName, testSingleSignOnURL)
}

func testAccIDPIntegrationCheck_Updated(samlDisplayName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_saml_integration",
			"samlp.0.display_name", samlDisplayName),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_saml_integration",
			"samlp.0.disabled", "true"),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_saml_integration",
			"samlp.0.config.0.back_channel_supported", "true"),
	)
}

func testAccIDPIntegrationConfig_NotEmptyAlias() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "test_saml_integration" {
		draft_alias = "test-alias"
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIDPIntegrationCheck_NotEmptyAlias() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_saml_integration",
			"draft_alias", "test-alias"),
		resource.TestCheckResourceAttrPair(
			"cyral_integration_idp_okta.test_saml_integration", "id",
			"cyral_integration_idp_okta.test_saml_integration", "draft_alias"),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}
