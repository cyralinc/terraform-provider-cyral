package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	TestSingleSignOnURL = "https://some-test-sso-url.com"
)

func TestAccSSOIntegrationResource(t *testing.T) {
	samlDisplayName := "tf-test-saml-integration"
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSSOIntegrationConfig_EmptySamlp(),
				ExpectError: regexp.MustCompile(`At least 1 "samlp" blocks are required`),
			},
			{
				Config:      testAccSSOIntegrationConfig_EmptyConfig(),
				ExpectError: regexp.MustCompile(`At least 1 "config" blocks are required`),
			},
			{
				Config: testAccSSOIntegrationConfig_EmptySSOUrl(),
				ExpectError: regexp.MustCompile(
					`The argument "single_sign_on_service_url" is required`),
			},
			{
				Config: testAccSSOIntegrationConfig_ADFS_DefaultValues(),
				Check:  testAccSSOIntegrationCheck_ADFS_DefaultValues(),
			},
			{
				Config: testAccSSOIntegrationConfig_AAD_DefaultValues(),
				Check:  testAccSSOIntegrationCheck_AAD_DefaultValues(),
			},
			{
				Config: testAccSSOIntegrationConfig_Forgerock_DefaultValues(),
				Check:  testAccSSOIntegrationCheck_Forgerock_DefaultValues(),
			},
			{
				Config: testAccSSOIntegrationConfig_GSuite_DefaultValues(),
				Check:  testAccSSOIntegrationCheck_GSuite_DefaultValues(),
			},
			{
				Config: testAccSSOIntegrationConfig_PingOne_DefaultValues(),
				Check:  testAccSSOIntegrationCheck_PingOne_DefaultValues(),
			},
			{
				Config: testAccSSOIntegrationConfig_Okta_DefaultValues(),
				Check:  testAccSSOIntegrationCheck_Okta_DefaultValues(),
			},
			{
				Config: testAccSSOIntegrationConfig_Updated(samlDisplayName),
				Check:  testAccSSOIntegrationCheck_Updated(samlDisplayName),
			},
			{
				Config: testAccSSOIntegrationConfig_NotEmptyAlias(),
				Check:  testAccSSOIntegrationCheck_NotEmptyAlias(),
			},
		},
	})
}

func testAccSSOIntegrationConfig_EmptySamlp() string {
	return `
	resource "cyral_integration_sso_okta" "test_saml_integration" {
	}
	`
}

func testAccSSOIntegrationConfig_EmptyConfig() string {
	return `
	resource "cyral_integration_sso_okta" "test_saml_integration" {
		samlp {
		}
	}
	`
}

func testAccSSOIntegrationConfig_EmptySSOUrl() string {
	return `
	resource "cyral_integration_sso_okta" "test_saml_integration" {
		samlp {
			config {
			}
		}
	}
	`
}

func testAccSSOIntegrationConfig_ADFS_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_sso_adfs" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, TestSingleSignOnURL)
}

func testAccSSOIntegrationCheck_ADFS_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sso_adfs.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", TestSingleSignOnURL),
	)
}

func testAccSSOIntegrationConfig_AAD_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_sso_aad" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, TestSingleSignOnURL)
}

func testAccSSOIntegrationCheck_AAD_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sso_aad.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", TestSingleSignOnURL),
	)
}

func testAccSSOIntegrationConfig_Forgerock_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_sso_forgerock" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, TestSingleSignOnURL)
}

func testAccSSOIntegrationCheck_Forgerock_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sso_forgerock.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", TestSingleSignOnURL),
	)
}

func testAccSSOIntegrationConfig_GSuite_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_sso_gsuite" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, TestSingleSignOnURL)
}

func testAccSSOIntegrationCheck_GSuite_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sso_gsuite.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", TestSingleSignOnURL),
	)
}

func testAccSSOIntegrationConfig_PingOne_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_sso_ping_one" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, TestSingleSignOnURL)
}

func testAccSSOIntegrationCheck_PingOne_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sso_ping_one.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", TestSingleSignOnURL),
	)
}

func testAccSSOIntegrationConfig_Okta_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_sso_okta" "test_saml_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, TestSingleSignOnURL)
}

func testAccSSOIntegrationCheck_Okta_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sso_okta.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", TestSingleSignOnURL),
	)
}

func testAccSSOIntegrationConfig_Updated(samlDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_sso_okta" "test_saml_integration" {
		samlp {
			display_name = "%s"
			disabled = true
			config {
				single_sign_on_service_url = "%s"
				back_channel_supported = true
			}
		}
	}
	`, samlDisplayName, TestSingleSignOnURL)
}

func testAccSSOIntegrationCheck_Updated(samlDisplayName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sso_okta.test_saml_integration",
			"samlp.0.display_name", samlDisplayName),
		resource.TestCheckResourceAttr("cyral_integration_sso_okta.test_saml_integration",
			"samlp.0.disabled", "true"),
		resource.TestCheckResourceAttr("cyral_integration_sso_okta.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", TestSingleSignOnURL),
		resource.TestCheckResourceAttr("cyral_integration_sso_okta.test_saml_integration",
			"samlp.0.config.0.back_channel_supported", "true"),
	)
}

func testAccSSOIntegrationConfig_NotEmptyAlias() string {
	return fmt.Sprintf(`
	resource "cyral_integration_sso_okta" "test_saml_integration" {
		draft_alias = "test-alias"
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, TestSingleSignOnURL)
}

func testAccSSOIntegrationCheck_NotEmptyAlias() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_sso_okta.test_saml_integration",
			"draft_alias", "test-alias"),
		resource.TestCheckResourceAttrPair(
			"cyral_integration_sso_okta.test_saml_integration", "id",
			"cyral_integration_sso_okta.test_saml_integration", "draft_alias"),
		resource.TestCheckResourceAttr("cyral_integration_sso_okta.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", TestSingleSignOnURL),
	)
}
