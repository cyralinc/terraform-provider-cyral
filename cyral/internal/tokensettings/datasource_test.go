package tokensettings_test

import (
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/tokensettings"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	completeDataSourceName = "data." + completeResourceName
)

func TestAccAccessTokenSettingsDataSource(t *testing.T) {
	testSteps := []resource.TestStep{
		{
			Config: testAccAccessTokenSettingsDataSourceConfig(),
			Check:  testAccAccessTokenSettingsDataSourceCheck(),
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps:             testSteps,
	})
}

func testAccAccessTokenSettingsDataSourceConfig() string {
	return `
	data "cyral_access_token_settings" "token_settings" {
	}
	`
}

func testAccAccessTokenSettingsDataSourceCheck() resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(completeDataSourceName, utils.IDKey),
		resource.TestCheckResourceAttrSet(completeDataSourceName, tokensettings.MaxValidityKey),
		resource.TestCheckResourceAttrSet(completeDataSourceName, tokensettings.DefaultValidityKey),
		resource.TestCheckResourceAttrSet(completeDataSourceName, tokensettings.MaxNumberOfTokensPerUserKey),
		resource.TestCheckResourceAttrSet(completeDataSourceName, tokensettings.OfflineTokenValidationKey),
	}
	return resource.ComposeTestCheckFunc(
		testCheckFuncs...,
	)
}
