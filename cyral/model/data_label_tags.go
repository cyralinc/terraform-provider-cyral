package model

type DataLabelTags []string

func (dlt DataLabelTags) AsInterface() []interface{} {
	var tagIfaces []interface{}
	for _, tag := range dlt {
		tagIfaces = append(tagIfaces, tag)
	}
	return tagIfaces
}
