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
	//samlDisplayName := "tf-test-saml-integration"
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck: func() {
			if v := os.Getenv(EnvVarSAMLMetadataURL); v == "" {
				t.Fatalf("%q must be set for TestAccSAMLIntegrationResource acceptance tests",
					EnvVarSAMLMetadataURL)
			}
		},
		Steps: []resource.TestStep{
			/* {
				Config:      testAccRepoConfAnalysisConfig_ErrorRedact(repoName),
				ExpectError: regexp.MustCompile(`Error running pre-apply refresh: exit status 1`),
			}, */
			{
				Config: testAccSAMLIntegrationConfig_DefaultValues(),
				Check:  testAccSAMLIntegrationCheck_DefaultValues(),
			},
			/* {
				Config: testAccRepoConfAnalysisConfig_Updated(repoName),
				Check:  testAccRepoConfAnalysisCheck_Updated(),
			}, */
		},
	})
}

func testAccSAMLIntegrationConfig_DefaultValues() string {
	return fmt.Sprintf(`
	resource "cyral_integration_saml" "test_saml_integration" {
		identity_provider = "okta"
		samlp {
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}
	`, os.Getenv(EnvVarSAMLMetadataURL))
}

func testAccSAMLIntegrationCheck_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_integration_saml.test_saml_integration",
			"identity_provider", "okta"),
		resource.TestCheckResourceAttr("cyral_integration_saml.test_saml_integration",
			"samlp.0.config.0.single_sign_on_service_url", os.Getenv(EnvVarSAMLMetadataURL)),
	)
}
