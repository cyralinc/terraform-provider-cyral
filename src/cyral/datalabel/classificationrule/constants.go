package classificationrule

import "github.com/cyralinc/terraform-provider-cyral/src/utils"

type Type string

const (
	Unknown = Type("UNKNOWN")
	Rego    = Type("REGO")
)

type Status string

const (
	Enabled  = Status("ENABLED")
	Disabled = Status("DISABLED")
)

func Types() []Type {
	return []Type{
		Unknown,
		Rego,
	}
}

func TypesAsString() []string {
	return utils.ToSliceOfString[Type](Types(), func(t Type) string {
		return string(t)
	})
}

func Statuses() []Status {
	return []Status{
		Enabled,
		Disabled,
	}
}

func StatusesAsString() []string {
	return utils.ToSliceOfString[Status](Statuses(), func(s Status) string {
		return string(s)
	})
}
