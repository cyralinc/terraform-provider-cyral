package datalabel

import "github.com/cyralinc/terraform-provider-cyral/src/utils"

type Type string

const (
	TypeUnknown = Type("UNKNOWN")
	Predefined  = Type("PREDEFINED")
	Custom      = Type("CUSTOM")
	Default     = TypeUnknown
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
