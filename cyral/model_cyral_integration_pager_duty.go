package cyral

type PagerDutyIntegration struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Parameters   string `json:"parameters,omitempty"`
	Purpose      string `json:"purpose,omitempty"`
	Category     string `json:"category,omitempty"`
	TemplateType string `json:"templateType,omitempty"`
}
