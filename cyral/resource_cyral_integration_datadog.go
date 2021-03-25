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

type CreateDatadogIntegrationResponse struct {
	ID string `json:"ID"`
}

type DatadogIntegrationData struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

func resourceIntegrationDatadog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationDatadogCreate,
		ReadContext:   resourceIntegrationDatadogRead,
		UpdateContext: resourceIntegrationDatadogUpdate,
		DeleteContext: resourceIntegrationDatadogDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceIntegrationDatadogCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationDatadogCreate")
	c := m.(*client.Client)

	resourceData := getDatadogIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/datadog", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create integration", fmt.Sprintf("%v", err))
	}

	response := CreateDatadogIntegrationResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceIntegrationDatadogCreate")

	return resourceIntegrationDatadogRead(ctx, d, m)
}

func resourceIntegrationDatadogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationDatadogRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := DatadogIntegrationData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("api_key", response.APIKey)

	log.Printf("[DEBUG] End resourceIntegrationDatadogRead")

	return diag.Diagnostics{}
}

func resourceIntegrationDatadogUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationDatadogUpdate")
	c := m.(*client.Client)

	resourceData := getDatadogIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationDatadogUpdate")

	return resourceIntegrationDatadogRead(ctx, d, m)
}

func resourceIntegrationDatadogDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationDatadogDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/datadog/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationDatadogDelete")

	return diag.Diagnostics{}
}

func getDatadogIntegrationDataFromResource(c *client.Client, d *schema.ResourceData) DatadogIntegrationData {
	return DatadogIntegrationData{
		ID:     d.Id(),
		Name:   d.Get("name").(string),
		APIKey: d.Get("api_key").(string),
	}
}
