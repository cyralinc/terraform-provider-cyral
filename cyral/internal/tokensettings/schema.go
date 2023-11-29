package tokensettings

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	MaxValidityKey              = "max_validity"
	DefaultValidityKey          = "default_validity"
	MaxNumberOfTokensPerUserKey = "max_number_of_tokens_per_user"
	OfflineTokenValidationKey   = "offline_token_validation"
	TokenLengthKey              = "token_length"
)

func getAccessTokenSettingsSchema(isDataSource bool) map[string]*schema.Schema {
	idDescriptionFmt := fmt.Sprintf("The %s identifier. It's set as `%s`.", "%s", accessTokenSettingsID)
	resourceSchema := map[string]*schema.Schema{
		utils.IDKey: {
			Description: fmt.Sprintf(idDescriptionFmt, "resource"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		MaxValidityKey: {
			Description: fmt.Sprintf(
				utils.DurationFieldDescriptionFmt,
				"The maximum duration that a user can request for access token validity. Defaults to `36000s`",
			),
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "36000s",
			ValidateFunc: utils.ValidationDurationString,
		},
		DefaultValidityKey: {
			Description: fmt.Sprintf(
				utils.DurationFieldDescriptionFmt,
				"The default duration used for access token validity. Defaults to `36000s`",
			),
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "36000s",
			ValidateFunc: utils.ValidationDurationString,
		},
		MaxNumberOfTokensPerUserKey: {
			Description: "The maximum number of access tokens that a user can have at the same time. " +
				"Must be between `1` and `5` (inclusive). Defaults to `3`.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      3,
			ValidateFunc: validation.IntBetween(1, 5),
		},
		OfflineTokenValidationKey: {
			Description: "The configuration that determines if the sidecar should perform access token " +
				"validation independently using cached token values. If this is `true`, the sidecar will be " +
				"able to validate and authenticate database access even when it cannot reach the Control Plane. " +
				"Defaults to `true`.",
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		TokenLengthKey: {
			Description: "The number of characters of the access token plaintext value. Valid values are `8`, " +
				"`12` and `16`. Defaults to `16`.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      16,
			ValidateFunc: validation.IntInSlice([]int{8, 12, 16}),
		},
	}

	if isDataSource {
		dataSourceSchema := utils.ConvertSchemaFieldsToComputed(resourceSchema)
		dataSourceSchema[utils.IDKey].Description = fmt.Sprintf(idDescriptionFmt, "data source")
		return dataSourceSchema
	}
	return resourceSchema
}
