package tokensettings_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/tokensettings"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	resourceName         = "token_settings"
	completeResourceName = "cyral_access_token_settings." + resourceName
)

func TestAccAccessTokenSettingsResource(t *testing.T) {
	testSteps := []resource.TestStep{}
	testSteps = append(testSteps, getRequiredArgumentTestSteps()...)
	testSteps = append(
		testSteps,
		[]resource.TestStep{
			{
				Config: testAccAccessTokenSettingsConfig_OnlyRequiredArguments(accessTokenSettingsOnlyRequiredArguments),
				Check:  testAccAccessTokenSettingsCheck(accessTokenSettingsOnlyRequiredArguments),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      completeResourceName,
			},
		}...,
	)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps:             testSteps,
	})
}

func getRequiredArgumentTestSteps() []resource.TestStep {
	requiredArgumentsTestSteps := []resource.TestStep{}
	requiredArguments := []string{
		tokensettings.MaxValidityKey,
		tokensettings.DefaultValidityKey,
		tokensettings.MaxNumberOfTokensPerUserKey,
		tokensettings.OfflineTokenValidationKey,
	}
	for _, argument := range requiredArguments {
		requiredArgumentsTestSteps = append(requiredArgumentsTestSteps, resource.TestStep{
			Config: testAccAccessTokenSettingsConfig_EmptyConfig(),
			ExpectError: regexp.MustCompile(
				fmt.Sprintf(`The argument "%s" is required`, argument),
			),
		})
	}
	return requiredArgumentsTestSteps
}

func testAccAccessTokenSettingsConfig_EmptyConfig() string {
	return `
	resource "cyral_access_token_settings" "token_settings" {
	}
	`
}

type AccessTokenSettingsTestParameters struct {
	settings tokensettings.AccessTokenSettings
}

var (
	accessTokenSettingsOnlyRequiredArguments = AccessTokenSettingsTestParameters{
		settings: tokensettings.AccessTokenSettings{
			MaxValidity:              "72000s",
			DefaultValidity:          "36000s",
			MaxNumberOfTokensPerUser: 5,
			OfflineTokenValidation:   true,
		},
	}
)

func testAccAccessTokenSettingsConfig_OnlyRequiredArguments(
	parameters AccessTokenSettingsTestParameters,
) string {
	return fmt.Sprintf(`
	resource "cyral_access_token_settings" "token_settings" {
		max_validity = %q
		default_validity = %q
		max_number_of_tokens_per_user = %d
		offline_token_validation = %t
	}
	`,
		parameters.settings.MaxValidity,
		parameters.settings.DefaultValidity,
		parameters.settings.MaxNumberOfTokensPerUser,
		parameters.settings.OfflineTokenValidation,
	)
}

func testAccAccessTokenSettingsCheck(
	parameters AccessTokenSettingsTestParameters,
) resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(completeResourceName, utils.IDKey),
		resource.TestCheckResourceAttr(
			completeResourceName,
			tokensettings.MaxValidityKey,
			parameters.settings.MaxValidity,
		),
		resource.TestCheckResourceAttrSet(completeResourceName, utils.IDKey),
		resource.TestCheckResourceAttr(
			completeResourceName,
			tokensettings.DefaultValidityKey,
			parameters.settings.DefaultValidity,
		),
		resource.TestCheckResourceAttrSet(completeResourceName, utils.IDKey),
		resource.TestCheckResourceAttr(
			completeResourceName,
			tokensettings.MaxNumberOfTokensPerUserKey,
			fmt.Sprintf("%d", parameters.settings.MaxNumberOfTokensPerUser),
		),
		resource.TestCheckResourceAttrSet(completeResourceName, utils.IDKey),
		resource.TestCheckResourceAttr(
			completeResourceName,
			tokensettings.OfflineTokenValidationKey,
			fmt.Sprintf("%t", parameters.settings.OfflineTokenValidation),
		),
	}
	return resource.ComposeTestCheckFunc(
		testCheckFuncs...,
	)
}
