package cyral

// HCVaultIntegration defines the necessary data for Hashicorp Vault integration
type HCVaultIntegration struct {
	AuthMethod string `json:"authMethod" tfgen:"auth_method,required"`
	ID         string `json:"id" tfgen:"id,computed"`
	AuthType   string `json:"authType" tfgen:"auth_type,required"`
	Name       string `json:"name" tfgen:"name,required"`
	Server     string `json:"server" tfgen:"server,required,sensitive"`
}
