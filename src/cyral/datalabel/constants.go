package datalabel

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
	types := Types()
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}
