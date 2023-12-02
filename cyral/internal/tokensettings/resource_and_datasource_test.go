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
	resourceName           = "token_settings"
	completeResourceName   = "cyral_access_token_settings." + resourceName
	completeDataSourceName = "data." + completeResourceName
)

// Given that the access token setting is a global resource, we test both resource and
// data source in this single test function, so that the tests are not run in parallel and
// don't compromise each other.
func TestAccAccessTokenSettingsResource(t *testing.T) {
	testSteps := []resource.TestStep{
		{
			Config: testAccAccessTokenSettingsConfig_MaxValidityValidation(),
			ExpectError: regexp.MustCompile(fmt.Sprintf(
				`expected %s to end with a 's' suffix.`,
				tokensettings.MaxValidityKey,
			)),
		},
		{
			Config: testAccAccessTokenSettingsConfig_DefaultValidityValidation(),
			ExpectError: regexp.MustCompile(fmt.Sprintf(
				`expected %s to end with a 's' suffix.`,
				tokensettings.DefaultValidityKey,
			)),
		},
		{
			Config: testAccAccessTokenSettingsConfig_MaxNumberOfTokensPerUserValidation(),
			ExpectError: regexp.MustCompile(fmt.Sprintf(
				`expected %s to be in the range \(1 - 5\)`,
				tokensettings.MaxNumberOfTokensPerUserKey,
			)),
		},
		{
			Config: testAccAccessTokenSettingsConfig_TokenLengthValidation(),
			ExpectError: regexp.MustCompile(fmt.Sprintf(
				`expected %s to be one of \[8 12 16\]`,
				tokensettings.TokenLengthKey,
			)),
		},
		{
			Config: testAccAccessTokenSettingsConfig_NoFieldsSet(),
			Check:  testAccAccessTokenSettingsCheck(accessTokenSettingsDefaultValues),
		},
		{
			Config: testAccAccessTokenSettingsConfig_AllFieldsSet(accessTokenSettingsAllFieldsSet),
			Check:  testAccAccessTokenSettingsCheck(accessTokenSettingsAllFieldsSet),
		},
		{
			// Needed to delete the resource before the data source is used to retrieve the current
			// config in the next step.
			Config: testAccAccessTokenSettingsConfig_ResourceDeletion(),
		},
		{
			Config: testAccAccessTokenSettingsConfig_TokenSettingsResetToDefault(),
			Check:  testAccAccessTokenSettingsDataSourceCheck(accessTokenSettingsDefaultValues),
		},
		{
			Config: testAccAccessTokenSettingsConfig_OnlySomeFieldsSet(accessTokenSettingsOnlySomeFieldsSet),
			Check:  testAccAccessTokenSettingsCheck(accessTokenSettingsOnlySomeFieldsSet),
		},
		{
			ImportState:       true,
			ImportStateVerify: true,
			ResourceName:      completeResourceName,
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps:             testSteps,
	})
}

type AccessTokenSettingsTestParameters struct {
	settings tokensettings.AccessTokenSettings
}

var (
	accessTokenSettingsAllFieldsSet = AccessTokenSettingsTestParameters{
		settings: tokensettings.AccessTokenSettings{
			MaxValidity:              "86400s",
			DefaultValidity:          "43200s",
			MaxNumberOfTokensPerUser: 1,
			OfflineTokenValidation:   false,
			TokenLength:              12,
		},
	}
	accessTokenSettingsDefaultValues = AccessTokenSettingsTestParameters{
		settings: tokensettings.AccessTokenSettings{
			MaxValidity:              "36000s",
			DefaultValidity:          "36000s",
			MaxNumberOfTokensPerUser: 3,
			OfflineTokenValidation:   true,
			TokenLength:              16,
		},
	}
	accessTokenSettingsOnlySomeFieldsSet = AccessTokenSettingsTestParameters{
		settings: tokensettings.AccessTokenSettings{
			MaxValidity:              "72000s",
			DefaultValidity:          accessTokenSettingsDefaultValues.settings.DefaultValidity,
			MaxNumberOfTokensPerUser: accessTokenSettingsDefaultValues.settings.MaxNumberOfTokensPerUser,
			OfflineTokenValidation:   accessTokenSettingsDefaultValues.settings.OfflineTokenValidation,
			TokenLength:              8,
		},
	}
)

func testAccAccessTokenSettingsConfig_MaxValidityValidation() string {
	return fmt.Sprintf(`
	resource "cyral_access_token_settings" "%s" {
		max_validity = "300"
	}
	`, resourceName)
}

func testAccAccessTokenSettingsConfig_DefaultValidityValidation() string {
	return fmt.Sprintf(`
	resource "cyral_access_token_settings" "%s" {
		max_validity = "36000s"
		default_validity = "300"
	}
	`, resourceName)
}

func testAccAccessTokenSettingsConfig_MaxNumberOfTokensPerUserValidation() string {
	return fmt.Sprintf(`
	resource "cyral_access_token_settings" "%s" {
		max_validity = "36000s"
		default_validity = "36000s"
		max_number_of_tokens_per_user = 0
	}
	`, resourceName)
}

func testAccAccessTokenSettingsConfig_TokenLengthValidation() string {
	return fmt.Sprintf(`
	resource "cyral_access_token_settings" "%s" {
		max_validity = "36000s"
		default_validity = "36000s"
		max_number_of_tokens_per_user = 3
		token_length = 1
	}
	`, resourceName)
}

func testAccAccessTokenSettingsConfig_NoFieldsSet() string {
	return fmt.Sprintf(`
	resource "cyral_access_token_settings" "%s" {
	}
	`, resourceName)
}

func testAccAccessTokenSettingsConfig_AllFieldsSet(
	parameters AccessTokenSettingsTestParameters,
) string {
	return fmt.Sprintf(`
	resource "cyral_access_token_settings" "%s" {
		max_validity = %q
		default_validity = %q
		max_number_of_tokens_per_user = %d
		offline_token_validation = %t
		token_length = %d
	}
	`,
		resourceName,
		parameters.settings.MaxValidity,
		parameters.settings.DefaultValidity,
		parameters.settings.MaxNumberOfTokensPerUser,
		parameters.settings.OfflineTokenValidation,
		parameters.settings.TokenLength,
	)
}

func testAccAccessTokenSettingsConfig_ResourceDeletion() string {
	return `
	output "dummy_output" {
		value = "dummy-value"
	}
	`
}

func testAccAccessTokenSettingsConfig_TokenSettingsResetToDefault() string {
	return fmt.Sprintf(`
	data "cyral_access_token_settings" "%s" {
	}
	`, resourceName)
}

func testAccAccessTokenSettingsConfig_OnlySomeFieldsSet(
	parameters AccessTokenSettingsTestParameters,
) string {
	return fmt.Sprintf(`
	resource "cyral_access_token_settings" "%s" {
		max_validity = %q
		token_length = %d
	}
	`,
		resourceName,
		parameters.settings.MaxValidity,
		parameters.settings.TokenLength,
	)
}

func testAccAccessTokenSettingsCheck(
	parameters AccessTokenSettingsTestParameters,
) resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			completeResourceName,
			utils.IDKey,
			"settings/access_token",
		),
		resource.TestCheckResourceAttr(
			completeResourceName,
			tokensettings.MaxValidityKey,
			parameters.settings.MaxValidity,
		),
		resource.TestCheckResourceAttr(
			completeResourceName,
			tokensettings.DefaultValidityKey,
			parameters.settings.DefaultValidity,
		),
		resource.TestCheckResourceAttr(
			completeResourceName,
			tokensettings.MaxNumberOfTokensPerUserKey,
			fmt.Sprintf("%d", parameters.settings.MaxNumberOfTokensPerUser),
		),
		resource.TestCheckResourceAttr(
			completeResourceName,
			tokensettings.OfflineTokenValidationKey,
			fmt.Sprintf("%t", parameters.settings.OfflineTokenValidation),
		),
		resource.TestCheckResourceAttr(
			completeResourceName,
			tokensettings.TokenLengthKey,
			fmt.Sprintf("%d", parameters.settings.TokenLength),
		),
	}
	return resource.ComposeTestCheckFunc(
		testCheckFuncs...,
	)
}

func testAccAccessTokenSettingsDataSourceCheck(
	parameters AccessTokenSettingsTestParameters,
) resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			completeDataSourceName,
			utils.IDKey,
			"settings/access_token",
		),
		resource.TestCheckResourceAttr(
			completeDataSourceName,
			tokensettings.MaxValidityKey,
			parameters.settings.MaxValidity,
		),
		resource.TestCheckResourceAttr(
			completeDataSourceName,
			tokensettings.DefaultValidityKey,
			parameters.settings.DefaultValidity,
		),
		resource.TestCheckResourceAttr(
			completeDataSourceName,
			tokensettings.MaxNumberOfTokensPerUserKey,
			fmt.Sprintf("%d", parameters.settings.MaxNumberOfTokensPerUser),
		),
		resource.TestCheckResourceAttr(
			completeDataSourceName,
			tokensettings.OfflineTokenValidationKey,
			fmt.Sprintf("%t", parameters.settings.OfflineTokenValidation),
		),
		resource.TestCheckResourceAttr(
			completeDataSourceName,
			tokensettings.TokenLengthKey,
			fmt.Sprintf("%d", parameters.settings.TokenLength),
		),
	}
	return resource.ComposeTestCheckFunc(
		testCheckFuncs...,
	)
}
