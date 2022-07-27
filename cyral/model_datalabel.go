package cyral

const (
	dataLabelTypeUnknown    = "UNKNOWN"
	dataLabelTypePredefined = "PREDEFINED"
	dataLabelTypeCustom     = "CUSTOM"
	defaultDataLabelType    = dataLabelTypeUnknown
)

func dataLabelTypes() []string {
	return []string{
		dataLabelTypeUnknown,
		dataLabelTypePredefined,
		dataLabelTypeCustom,
	}
}

type DataLabel struct {
	Name        string   `json:"name,omitempty"`
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (dl *DataLabel) TagsAsInterface() []interface{} {
	var tagIfaces []interface{}
	for _, tag := range dl.Tags {
		tagIfaces = append(tagIfaces, tag)
	}
	return tagIfaces
}
