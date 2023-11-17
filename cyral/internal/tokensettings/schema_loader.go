package tokensettings

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MaxValidityKey              = "max_validity"
	DefaultValidityKey          = "default_validity"
	MaxNumberOfTokensPerUserKey = "max_number_of_tokens_per_user"
	OfflineTokenValidationKey   = "offline_token_validation"
)

func getAccessTokenSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MaxValidityKey: {
			Description:  "",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: utils.ValidationDurationString,
		},
		DefaultValidityKey: {
			Description:  "",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: utils.ValidationDurationString,
		},
		MaxNumberOfTokensPerUserKey: {
			Description: "",
			Type:        schema.TypeInt,
			Required:    true,
		},
		OfflineTokenValidationKey: {
			Description: "",
			Type:        schema.TypeBool,
			Required:    true,
		},
	}
}

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "tokensettings"
}

func (p *packageSchema) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{
		{
			Name:   "cyral_access_token_settings",
			Type:   core.DataSourceSchemaType,
			Schema: dataSourceSchema,
		},
		{
			Name:   "cyral_access_token_settings",
			Type:   core.ResourceSchemaType,
			Schema: resourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema {
	return &packageSchema{}
}
