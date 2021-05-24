package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (data HCVaultIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("auth_method", data.AuthMethod)
	d.Set("id", data.ID)
	d.Set("auth_type", data.AuthType)
	d.Set("name", data.Name)
	d.Set("server", data.Server)

}

func (data *HCVaultIntegration) ReadFromSchema(d *schema.ResourceData) {
	data.AuthMethod = d.Get("auth_method").(string)
	data.ID = d.Get("id").(string)
	data.AuthType = d.Get("auth_type").(string)
	data.Name = d.Get("name").(string)
	data.Server = d.Get("server").(string)

}

var ReadHCVaultIntegrationConfig = ResourceOperationConfig{
	Name:       "HCVaultIntegrationResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/secretProviders/hcvault/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &HCVaultIntegration{},
}

func resourceIntegrationHCVault() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "HCVaultIntegrationResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/secretProviders/hcvault", c.ControlPlane)
				},
				ResourceData: &HCVaultIntegration{},
				ResponseData: &IDBasedResponse{},
			}, ReadHCVaultIntegrationConfig,
		),
		ReadContext: ReadResource(ReadHCVaultIntegrationConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "HCVaultIntegrationResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/secretProviders/hcvault/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &HCVaultIntegration{},
			}, ReadHCVaultIntegrationConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "HCVaultIntegrationResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/secretProviders/hcvault/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"auth_method": {
				Required: true,
				Type:     schema.TypeString,
			},
			"id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"auth_type": {
				Required: true,
				Type:     schema.TypeString,
			},
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"server": {
				Required:  true,
				Sensitive: true,
				Type:      schema.TypeString,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
