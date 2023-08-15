package classificationrule

type Type string

const (
	Unknown = Type("UNKNOWN")
	Rego    = Type("REGO")
)

type Status string

const (
	Enabled  = "ENABLED"
	Disabled = "DISABLED"
)

func Types() []Type {
	return []Type{
		Unknown,
		Rego,
	}
}

func Statuses() []Status {
	return []Status{
		Enabled,
		Disabled,
	}
}

func StatusesAsString() []string {
	statuses := Statuses()
	result := make([]string, len(statuses))
	for i, v := range statuses {
		result[i] = string(v)
	}
	return result
}
