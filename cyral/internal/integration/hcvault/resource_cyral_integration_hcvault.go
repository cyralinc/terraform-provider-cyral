package hcvault

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

var ReadHCVaultIntegrationConfig = core.ResourceOperationConfig{
	Name:       "HCVaultIntegrationResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/secretProviders/hcvault/%s", c.ControlPlane, d.Id())
	},
	NewResponseData:     func(_ *schema.ResourceData) core.SchemaWriter { return &HCVaultIntegration{} },
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Integration hcvault"},
}

func ResourceIntegrationHCVault() *schema.Resource {
	return &schema.Resource{
		Description: "Manages integration with Hashicorp Vault to store secrets.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "HCVaultIntegrationResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/secretProviders/hcvault", c.ControlPlane)
				},
				NewResourceData: func() core.SchemaReader { return &HCVaultIntegration{} },
			}, ReadHCVaultIntegrationConfig,
		),
		ReadContext: core.ReadResource(ReadHCVaultIntegrationConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "HCVaultIntegrationResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/secretProviders/hcvault/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() core.SchemaReader { return &HCVaultIntegration{} },
			}, ReadHCVaultIntegrationConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "HCVaultIntegrationResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/secretProviders/hcvault/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"auth_method": {
				Description: "Authentication method for the integration.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"auth_type": {
				Description: "Authentication type for the integration.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"name": {
				Description: "Integration name that will be used internally in the control plane.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"server": {
				Description: "Server on which the vault service is running.",
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
