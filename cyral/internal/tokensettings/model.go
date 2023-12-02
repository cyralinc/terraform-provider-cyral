package tokensettings

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

type AccessTokenSettings struct {
	MaxValidity              string `json:"maxValidity"`
	DefaultValidity          string `json:"defaultValidity"`
	MaxNumberOfTokensPerUser uint32 `json:"maxNumberOfTokensPerUser"`
	OfflineTokenValidation   bool   `json:"offlineTokenValidation"`
	TokenLength              uint32 `json:"tokenLength"`
}

func (settings *AccessTokenSettings) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set(MaxValidityKey, settings.MaxValidity); err != nil {
		return fmt.Errorf(utils.ErrorSettingFieldFmt, MaxValidityKey, err)
	}
	if err := d.Set(DefaultValidityKey, settings.DefaultValidity); err != nil {
		return fmt.Errorf(utils.ErrorSettingFieldFmt, DefaultValidityKey, err)
	}
	if err := d.Set(MaxNumberOfTokensPerUserKey, settings.MaxNumberOfTokensPerUser); err != nil {
		return fmt.Errorf(utils.ErrorSettingFieldFmt, MaxNumberOfTokensPerUserKey, err)
	}
	if err := d.Set(OfflineTokenValidationKey, settings.OfflineTokenValidation); err != nil {
		return fmt.Errorf(utils.ErrorSettingFieldFmt, OfflineTokenValidationKey, err)
	}
	if err := d.Set(TokenLengthKey, settings.TokenLength); err != nil {
		return fmt.Errorf(utils.ErrorSettingFieldFmt, TokenLengthKey, err)
	}
	d.SetId(accessTokenSettingsID)
	return nil
}

func (settings *AccessTokenSettings) ReadFromSchema(d *schema.ResourceData) error {
	settings.MaxValidity = d.Get(MaxValidityKey).(string)
	settings.DefaultValidity = d.Get(DefaultValidityKey).(string)
	settings.MaxNumberOfTokensPerUser = uint32(d.Get(MaxNumberOfTokensPerUserKey).(int))
	settings.OfflineTokenValidation = d.Get(OfflineTokenValidationKey).(bool)
	settings.TokenLength = uint32(d.Get(TokenLengthKey).(int))
	return nil
}
