package model

type DataLabelClassificationRule struct {
	RuleType   string `json:"ruleType"`
	RuleCode   string `json:"ruleCode"`
	RuleStatus string `json:"status"`
}

func (dl *DataLabelClassificationRule) AsInterface() []interface{} {
	if dl == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"rule_type":   dl.RuleType,
		"rule_code":   dl.RuleCode,
		"rule_status": dl.RuleStatus,
	}}
}
