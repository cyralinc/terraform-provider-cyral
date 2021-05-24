package cyral

//go:generate go run ../tfgen HCVaultIntegration https://%s/v1/integrations/secretProviders/hcvault --output resource_cyral_integration_hcvault.go
type HCVaultIntegration struct {
	AuthMethod string `json:"authMethod" tfgen:"auth_method,required"`
	ID         string `json:"id" tfgen:"id,computed"`
	AuthType   string `json:"authType" tfgen:"auth_type,required"`
	Name       string `json:"name" tfgen:"name,required"`
	Server     string `json:"server" tfgen:"server,required,sensitive"`
}
