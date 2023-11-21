package tokensettings

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MaxValidityKey              = "max_validity"
	DefaultValidityKey          = "default_validity"
	MaxNumberOfTokensPerUserKey = "max_number_of_tokens_per_user"
	OfflineTokenValidationKey   = "offline_token_validation"
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
				"The maximum duration that a user can request for access token validity",
			),
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: utils.ValidationDurationString,
		},
		DefaultValidityKey: {
			Description: fmt.Sprintf(
				utils.DurationFieldDescriptionFmt,
				"The default duration used for access token validity",
			),
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: utils.ValidationDurationString,
		},
		MaxNumberOfTokensPerUserKey: {
			Description: "The maximum number of access tokens that a user can have at the same time.",
			Type:        schema.TypeInt,
			Required:    true,
		},
		OfflineTokenValidationKey: {
			Description: "The configuration that determines if the sidecar should perform access token " +
				"validation independently using cached token values. If this is true, the sidecar will be " +
				"able to validate and authenticate database access even when it cannot reach the Control Plane.",
			Type:     schema.TypeBool,
			Required: true,
		},
	}

	if isDataSource {
		dataSourceSchema := utils.ConvertSchemaFieldsToComputed(resourceSchema)
		dataSourceSchema[utils.IDKey].Description = fmt.Sprintf(idDescriptionFmt, "data source")
		return dataSourceSchema
	}
	return resourceSchema
}
