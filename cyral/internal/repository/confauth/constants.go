package confauth

import "github.com/cyralinc/terraform-provider-cyral/cyral/utils"

const (
	resourceName        = "cyral_repository_conf_auth"
	AccessTokenAuthType = "ACCESS_TOKEN"
	AwsIAMAuthType      = "AWS_IAM"
	DefaultAuthType     = AccessTokenAuthType
)

const (
	TLSEnable              = TLSType("enable")
	TLSEnableAndVerifyCert = TLSType("enableAndVerifyCert")
	TLSDisable             = TLSType("disable")
)

type TLSType string

func ClientTLSTypes() []TLSType {
	return []TLSType{
		TLSEnable,
		TLSDisable,
	}
}

func ClientTLSTypesAsString() []string {
	return utils.ToSliceOfString[TLSType](ClientTLSTypes(), func(t TLSType) string {
		return string(t)
	})
}

func RepoTLSTypes() []TLSType {
	return []TLSType{
		TLSEnable,
		TLSEnableAndVerifyCert,
		TLSDisable,
	}
}

func RepoTLSTypesAsString() []string {
	return utils.ToSliceOfString[TLSType](RepoTLSTypes(), func(t TLSType) string {
		return string(t)
	})
}
