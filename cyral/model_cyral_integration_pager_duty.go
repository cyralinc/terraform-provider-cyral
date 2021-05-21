package cyral

//go:generate go run ../tfgen PagerDutyIntegration "https://%s/v1/integrations/authorizationPolicies" --output resource_cyral_integration_pager_duty.go
type PagerDutyIntegration struct {
	ID       string `json:"id" tfgen:"id,computed"`
	Name     string `json:"name" tfgen:"name,required"`
	APIToken string `json:"api_token" tfgen:"api_token,required,sensitive"`
}
