package datalabel

import (
	"fmt"

	cs "github.com/cyralinc/terraform-provider-cyral/src/cyral/datalabel/classificationrule"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Tags []string

func (dlt Tags) AsInterface() []interface{} {
	var tagIfaces []interface{}
	for _, tag := range dlt {
		tagIfaces = append(tagIfaces, tag)
	}
	return tagIfaces
}

type DataLabel struct {
	Name               string                 `json:"name,omitempty"`
	Type               Type                   `json:"type,omitempty"`
	Description        string                 `json:"description,omitempty"`
	Tags               Tags                   `json:"tags,omitempty"`
	ClassificationRule *cs.ClassificationRule `json:"classificationRule,omitempty"`
	Implicit           bool                   `json:"implicit,omitempty"`
}

func (dl *DataLabel) WriteToSchema(d *schema.ResourceData) error {
	// // If the `description` field is not in the root of the schema, then
	// // it is a data source that is returning multiple entries.
	// if _, ok := d.GetOk("description"); ok {
	// 	if err := writeDataLabelsToDataSourceSchema([]*DataLabel{(*DataLabel)(dl)}, d); err != nil {
	// 		return err
	// 	}
	// 	d.SetId(uuid.New().String())
	// } else {
	if err := d.Set("description", dl.Description); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
	}

	if err := d.Set("tags", dl.Tags.AsInterface()); err != nil {
		return fmt.Errorf("error setting 'tags' field: %w", err)
	}

	if err := d.Set("classification_rule", dl.ClassificationRule.AsInterface()); err != nil {
		return fmt.Errorf("error setting 'classification_rule' field: %w", err)
	}

	d.SetId(dl.Name)
	// }

	return nil
}

func (dl *DataLabel) ReadFromSchema(d *schema.ResourceData) error {
	var tags []string
	tagIfaces := d.Get("tags").([]any)
	for _, tagIface := range tagIfaces {
		tags = append(tags, tagIface.(string))
	}

	var classificationRule *cs.ClassificationRule
	classificationRuleList := d.Get("classification_rule").(*schema.Set).List()
	if len(classificationRuleList) > 0 {
		classificationRuleMap := classificationRuleList[0].(map[string]any)
		classificationRule = &cs.ClassificationRule{
			RuleType:   classificationRuleMap["rule_type"].(string),
			RuleCode:   classificationRuleMap["rule_code"].(string),
			RuleStatus: classificationRuleMap["rule_status"].(string),
		}
	}
	dl.Name = d.Get("name").(string)
	dl.Type = Custom
	dl.Description = d.Get("description").(string)
	dl.Tags = tags
	dl.ClassificationRule = classificationRule

	return nil
}
