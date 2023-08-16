package cyral

const (
	typeUnknown    = "UNKNOWN"
	statusEnabled  = "ENABLED"
	statusDisabled = "DISABLED"
	// Data label type
	dataLabelTypePredefined = "PREDEFINED"
	dataLabelTypeCustom     = "CUSTOM"
	defaultDataLabelType    = typeUnknown
	// Classification rule type
	classificationRuleTypeRego = "REGO"
)

func dataLabelTypes() []string {
	return []string{
		typeUnknown,
		dataLabelTypePredefined,
		dataLabelTypeCustom,
	}
}

func classificationRuleTypes() []string {
	return []string{
		typeUnknown,
		classificationRuleTypeRego,
	}
}

func classificationRuleStatus() []string {
	return []string{
		statusEnabled,
		statusDisabled,
	}
}

type DataLabel struct {
	Name               string                       `json:"name,omitempty"`
	Type               string                       `json:"type,omitempty"`
	Description        string                       `json:"description,omitempty"`
	Tags               []string                     `json:"tags,omitempty"`
	ClassificationRule *DataLabelClassificationRule `json:"classificationRule,omitempty"`
	Implicit           bool                         `json:"implicit,omitempty"`
}

func (dl *DataLabel) TagsAsInterface() []interface{} {
	var tagIfaces []interface{}
	for _, tag := range dl.Tags {
		tagIfaces = append(tagIfaces, tag)
	}
	return tagIfaces
}

func (dl *DataLabel) ClassificationRuleAsInterface() []interface{} {
	if dl.ClassificationRule == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"rule_type":   dl.ClassificationRule.RuleType,
		"rule_code":   dl.ClassificationRule.RuleCode,
		"rule_status": dl.ClassificationRule.RuleStatus,
	}}
}

type DataLabelClassificationRule struct {
	RuleType   string `json:"ruleType"`
	RuleCode   string `json:"ruleCode"`
	RuleStatus string `json:"status"`
}
