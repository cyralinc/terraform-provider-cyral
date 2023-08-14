package model

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	DataLabelTypeUnknown    = "UNKNOWN"
	DataLabelTypePredefined = "PREDEFINED"
	DataLabelTypeCustom     = "CUSTOM"
	DefaultDataLabelType    = DataLabelTypeUnknown
)

func DataLabelTypes() []string {
	return []string{
		DataLabelTypeUnknown,
		DataLabelTypePredefined,
		DataLabelTypeCustom,
	}
}

type DataLabel struct {
	Name               string                       `json:"name,omitempty"`
	Type               string                       `json:"type,omitempty"`
	Description        string                       `json:"description,omitempty"`
	Tags               DataLabelTags                `json:"tags,omitempty"`
	ClassificationRule *DataLabelClassificationRule `json:"classificationRule,omitempty"`
	Implicit           bool                         `json:"implicit,omitempty"`
}

func (dl *DataLabel) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("description", dl.Description); err != nil {
		return err
	}

	if err := d.Set("tags", dl.Tags.AsInterface()); err != nil {
		return err
	}

	if err := d.Set("classification_rule", dl.ClassificationRule.AsInterface()); err != nil {
		return err
	}

	return nil
}

func (dl *DataLabel) ReadFromSchema(d *schema.ResourceData) error {
	var tags DataLabelTags
	tagIfaces := d.Get("tags").([]interface{})
	for _, tagIface := range tagIfaces {
		tags = append(tags, tagIface.(string))
	}
	dl.Name = d.Get("name").(string)
	dl.Type = DataLabelTypeCustom
	dl.Description = d.Get("description").(string)
	dl.Tags = tags
	//	dl.ClassificationRule = d.Get("")

	return nil
}
