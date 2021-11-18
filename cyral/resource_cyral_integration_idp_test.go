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

func TestAccIdPIntegrationResource(t *testing.T) {
	idpDisplayName := "tf-test-idp-integration"
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccIdPIntegrationConfig_EmptySamlp(),
				ExpectError: regexp.MustCompile(`At least 1 "samlp" blocks are required`),
			},
			{
				Config:      testAccIdPIntegrationConfig_EmptyConfig(),
				ExpectError: regexp.MustCompile(`At least 1 "config" blocks are required`),
			},
			{
				Config: testAccIdPIntegrationConfig_EmptySSOUrl(),
				ExpectError: regexp.MustCompile(
					`The argument "single_sign_on_service_url" is required`),
			},
			{
				Config: testAccIdPIntegrationConfig_ADFS_DefaultValues(),
				Check:  testAccIdPIntegrationCheck_ADFS_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_AAD_DefaultValues(),
				Check:  testAccIdPIntegrationCheck_AAD_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_Forgerock_DefaultValues(),
				Check:  testAccIdPIntegrationCheck_Forgerock_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_GSuite_DefaultValues(),
				Check:  testAccIdPIntegrationCheck_GSuite_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_PingOne_DefaultValues(),
				Check:  testAccIdPIntegrationCheck_PingOne_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_Okta_DefaultValues(),
				Check:  testAccIdPIntegrationCheck_Okta_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_Updated(idpDisplayName),
				Check:  testAccIdPIntegrationCheck_Updated(idpDisplayName),
			},
			{
				Config: testAccIdPIntegrationConfig_NotEmptyAlias(),
				Check:  testAccIdPIntegrationCheck_NotEmptyAlias(),
			},
		},
	})
}

func testAccIdPIntegrationConfig_EmptySamlp() string {
	return `
	resource "cyral_integration_idp_okta" "test_idp_integration" {
	}
	`
}

func testAccIdPIntegrationConfig_EmptyConfig() string {
	return `
	resource "cyral_integration_idp_okta" "test_idp_integration" {
		samlp {
		}
	}
	`
}

func testAccIdPIntegrationConfig_EmptySSOUrl() string {
	return `
	resource "cyral_integration_idp_okta" "test_idp_integration" {
		samlp {
			config {
			}
		}
	}
	`
}

func testAccIdPIntegrationConfig_ADFS_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_adfs" "test_idp_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIdPIntegrationCheck_ADFS_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_adfs.test_idp_integration",
			"id", regexp.MustCompile(`adfs.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_adfs.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_AAD_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_aad" "test_idp_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIdPIntegrationCheck_AAD_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_aad.test_idp_integration",
			"id", regexp.MustCompile(`aad.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_aad.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_Forgerock_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_forgerock" "test_idp_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIdPIntegrationCheck_Forgerock_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_forgerock.test_idp_integration",
			"id", regexp.MustCompile(`forgerock.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_forgerock.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_GSuite_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_gsuite" "test_idp_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIdPIntegrationCheck_GSuite_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_gsuite.test_idp_integration",
			"id", regexp.MustCompile(`gsuite.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_gsuite.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_PingOne_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_ping_one" "test_idp_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIdPIntegrationCheck_PingOne_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_ping_one.test_idp_integration",
			"id", regexp.MustCompile(`pingone.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_ping_one.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_Okta_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "test_idp_integration" {
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIdPIntegrationCheck_Okta_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"id", regexp.MustCompile(`okta.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_Updated(idpDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "test_idp_integration" {
		samlp {
			display_name = "%s"
			disabled = true
			config {
				single_sign_on_service_url = "%s"
				back_channel_supported = true
			}
		}
	}
	`, idpDisplayName, testSingleSignOnURL)
}

func testAccIdPIntegrationCheck_Updated(idpDisplayName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"id", regexp.MustCompile(`okta.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"samlp.0.display_name", idpDisplayName),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"samlp.0.disabled", "true"),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"samlp.0.config.0.back_channel_supported", "true"),
	)
}

func testAccIdPIntegrationConfig_NotEmptyAlias() string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "test_idp_integration" {
		draft_alias = "test-alias"
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, testSingleSignOnURL)
}

func testAccIdPIntegrationCheck_NotEmptyAlias() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"draft_alias", "test-alias"),
		resource.TestCheckResourceAttrPair(
			"cyral_integration_idp_okta.test_idp_integration", "id",
			"cyral_integration_idp_okta.test_idp_integration", "draft_alias"),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", testSingleSignOnURL),
	)
}
