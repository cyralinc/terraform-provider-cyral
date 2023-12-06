package datalabel

import "github.com/cyralinc/terraform-provider-cyral/cyral/utils"

type Type string

const (
	TypeUnknown = Type("UNKNOWN")
	Predefined  = Type("PREDEFINED")
	Custom      = Type("CUSTOM")
	Default     = TypeUnknown

	resourceName   = "cyral_datalabel"
	dataSourceName = "cyral_datalabel"
)

func Types() []Type {
	return []Type{
		TypeUnknown,
		Predefined,
		Custom,
	}
}

func TypesAsString() []string {
	return utils.ToSliceOfString[Type](Types(), func(t Type) string {
		return string(t)
	})
}
