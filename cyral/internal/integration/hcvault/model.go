package hcvault

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HCVaultIntegration defines the necessary data for Hashicorp Vault integration
type HCVaultIntegration struct {
	AuthMethod string `json:"authMethod" tfgen:"auth_method,required"`
	ID         string `json:"id" tfgen:"id,computed"`
	AuthType   string `json:"authType" tfgen:"auth_type,required"`
	Name       string `json:"name" tfgen:"name,required"`
	Server     string `json:"server" tfgen:"server,required,sensitive"`
}

func (data HCVaultIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("auth_method", data.AuthMethod)
	d.Set("id", data.ID)
	d.Set("auth_type", data.AuthType)
	d.Set("name", data.Name)
	d.Set("server", data.Server)
	return nil
}

func (data *HCVaultIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.AuthMethod = d.Get("auth_method").(string)
	data.ID = d.Get("id").(string)
	data.AuthType = d.Get("auth_type").(string)
	data.Name = d.Get("name").(string)
	data.Server = d.Get("server").(string)
	return nil
}
