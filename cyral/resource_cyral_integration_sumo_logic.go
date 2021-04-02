package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateSumoLogicIntegrationResponse struct {
	ID string `json:"id"`
}

func (response CreateSumoLogicIntegrationResponse) WriteResourceData(d *schema.ResourceData) {
	d.SetId(response.ID)
}

func (response *CreateSumoLogicIntegrationResponse) ReadResourceData(d *schema.ResourceData) {
	response.ID = d.Id()
}

type SumoLogicIntegrationData struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (data SumoLogicIntegrationData) WriteResourceData(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("address", data.Address)

}

func (data *SumoLogicIntegrationData) ReadResourceData(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.Address = d.Get("address").(string)
}

var CreateSumoLogicFunctionConfig = FunctionConfig{
	Name:       "SumoLogicResourceCreate",
	HttpMethod: http.MethodPost,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/sumologic", c.ControlPlane)
	},
	ResourceData:       &SumoLogicIntegrationData{},
	ResponseData:       &CreateSumoLogicIntegrationResponse{},
	ReadFunctionConfig: &ReadSumoLogicFunctionConfig,
}

var ReadSumoLogicFunctionConfig = FunctionConfig{
	Name:       "SumoLogicResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &SumoLogicIntegrationData{},
}

var UpdateSumoLogicFunctionConfig = FunctionConfig{
	Name:       "SumoLogicResourceUpdate",
	HttpMethod: http.MethodPut,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())
	},
	ResourceData:       &SumoLogicIntegrationData{},
	ReadFunctionConfig: &ReadSumoLogicFunctionConfig,
}

var DeleteSumoLogicFunctionConfig = FunctionConfig{
	Name:       "SumoLogicResourceDelete",
	HttpMethod: http.MethodDelete,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())
	},
}

func resourceIntegrationSumoLogic() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateSumoLogicFunctionConfig.Create,
		ReadContext:   ReadSumoLogicFunctionConfig.Read,
		UpdateContext: UpdateSumoLogicFunctionConfig.Update,
		DeleteContext: DeleteSumoLogicFunctionConfig.Delete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
