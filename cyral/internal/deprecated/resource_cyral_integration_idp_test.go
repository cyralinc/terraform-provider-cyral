package deprecated_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdPIntegrationResource(t *testing.T) {
	idpDisplayName := utils.AccTestName(utils.IntegrationIdPResourceName, "integration")

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
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
				Config: testAccIdPIntegrationConfig_ADFS_DefaultValues(idpDisplayName),
				Check:  testAccIdPIntegrationCheck_ADFS_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_AAD_DefaultValues(idpDisplayName),
				Check:  testAccIdPIntegrationCheck_AAD_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_Forgerock_DefaultValues(idpDisplayName),
				Check:  testAccIdPIntegrationCheck_Forgerock_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_GSuite_DefaultValues(idpDisplayName),
				Check:  testAccIdPIntegrationCheck_GSuite_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_PingOne_DefaultValues(idpDisplayName),
				Check:  testAccIdPIntegrationCheck_PingOne_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_Okta_DefaultValues(idpDisplayName),
				Check:  testAccIdPIntegrationCheck_Okta_DefaultValues(),
			},
			{
				Config: testAccIdPIntegrationConfig_Updated(idpDisplayName),
				Check:  testAccIdPIntegrationCheck_Updated(idpDisplayName),
			},
			{
				Config: testAccIdPIntegrationConfig_NotEmptyAlias(idpDisplayName),
				Check:  testAccIdPIntegrationCheck_NotEmptyAlias(),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"draft_alias"},
				ResourceName:            "cyral_integration_idp_okta.test_idp_integration",
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

func testAccIdPIntegrationConfig_ADFS_DefaultValues(idpDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_adfs" "test_idp_integration" {
		samlp {
			display_name = "%s"
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, idpDisplayName, utils.TestSingleSignOnURL)
}

func testAccIdPIntegrationCheck_ADFS_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_adfs.test_idp_integration",
			"id", regexp.MustCompile(`adfs.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_adfs.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", utils.TestSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_AAD_DefaultValues(idpDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_aad" "test_idp_integration" {
		samlp {
			display_name = "%s"
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, idpDisplayName, utils.TestSingleSignOnURL)
}

func testAccIdPIntegrationCheck_AAD_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_aad.test_idp_integration",
			"id", regexp.MustCompile(`aad.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_aad.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", utils.TestSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_Forgerock_DefaultValues(idpDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_forgerock" "test_idp_integration" {
		samlp {
			display_name = "%s"
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, idpDisplayName, utils.TestSingleSignOnURL)
}

func testAccIdPIntegrationCheck_Forgerock_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_forgerock.test_idp_integration",
			"id", regexp.MustCompile(`forgerock.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_forgerock.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", utils.TestSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_GSuite_DefaultValues(idpDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_gsuite" "test_idp_integration" {
		samlp {
			display_name = "%s"
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, idpDisplayName, utils.TestSingleSignOnURL)
}

func testAccIdPIntegrationCheck_GSuite_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_gsuite.test_idp_integration",
			"id", regexp.MustCompile(`gsuite.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_gsuite.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", utils.TestSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_PingOne_DefaultValues(idpDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_ping_one" "test_idp_integration" {
		samlp {
			display_name = "%s"
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, idpDisplayName, utils.TestSingleSignOnURL)
}

func testAccIdPIntegrationCheck_PingOne_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_ping_one.test_idp_integration",
			"id", regexp.MustCompile(`pingone.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_ping_one.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", utils.TestSingleSignOnURL),
	)
}

func testAccIdPIntegrationConfig_Okta_DefaultValues(idpDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "test_idp_integration" {
		samlp {
			display_name = "%s"
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, idpDisplayName, utils.TestSingleSignOnURL)
}

func testAccIdPIntegrationCheck_Okta_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"id", regexp.MustCompile(`okta.`)),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", utils.TestSingleSignOnURL),
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
	`, idpDisplayName, utils.TestSingleSignOnURL)
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
			"samlp.0.config.0.single_sign_on_service_url", utils.TestSingleSignOnURL),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"samlp.0.config.0.back_channel_supported", "true"),
	)
}

func testAccIdPIntegrationConfig_NotEmptyAlias(idpDisplayName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "test_idp_integration" {
		draft_alias = "%s"
		samlp {
			display_name = "%s"
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, utils.AccTestName(utils.IntegrationIdPResourceName, "test-alias"), idpDisplayName, utils.TestSingleSignOnURL)
}

func testAccIdPIntegrationCheck_NotEmptyAlias() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"draft_alias", utils.AccTestName(utils.IntegrationIdPResourceName, "test-alias")),
		resource.TestCheckResourceAttrPair(
			"cyral_integration_idp_okta.test_idp_integration", "id",
			"cyral_integration_idp_okta.test_idp_integration", "draft_alias"),
		resource.TestCheckResourceAttr("cyral_integration_idp_okta.test_idp_integration",
			"samlp.0.config.0.single_sign_on_service_url", utils.TestSingleSignOnURL),
	)
}
