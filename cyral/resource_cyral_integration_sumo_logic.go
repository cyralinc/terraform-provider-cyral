package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateSumoLogicIntegrationResponse struct {
	ID string `json:"id"`
}

type SumoLogicIntegrationData struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func resourceIntegrationSumoLogic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationSumoLogicCreate,
		ReadContext:   resourceIntegrationSumoLogicRead,
		UpdateContext: resourceIntegrationSumoLogicUpdate,
		DeleteContext: resourceIntegrationSumoLogicDelete,

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

func resourceIntegrationSumoLogicCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSumoLogicCreate")
	c := m.(*client.Client)

	resourceData := getSumoLogicIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/sumologic", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create integration", fmt.Sprintf("%v", err))
	}

	response := CreateSumoLogicIntegrationResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceIntegrationSumoLogicCreate")

	return resourceIntegrationSumoLogicRead(ctx, d, m)
}

func resourceIntegrationSumoLogicRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSumoLogicRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := SumoLogicIntegrationData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("address", response.Address)

	log.Printf("[DEBUG] End resourceIntegrationSumoLogicRead")

	return diag.Diagnostics{}
}

func resourceIntegrationSumoLogicUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSumoLogicUpdate")
	c := m.(*client.Client)

	resourceData := getSumoLogicIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationSumoLogicUpdate")

	return resourceIntegrationSumoLogicRead(ctx, d, m)
}

func resourceIntegrationSumoLogicDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSumoLogicDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/sumologic/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationSumoLogicDelete")

	return diag.Diagnostics{}
}

func getSumoLogicIntegrationDataFromResource(c *client.Client, d *schema.ResourceData) SumoLogicIntegrationData {
	return SumoLogicIntegrationData{
		Name:    d.Get("name").(string),
		Address: d.Get("address").(string),
	}
}
