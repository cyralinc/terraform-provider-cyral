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

type CreateSlackAlertsIntegrationResponse struct {
	ID string `json:"id"`
}

type SlackAlertsIntegrationData struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func resourceIntegrationSlackAlerts() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationSlackAlertsCreate,
		ReadContext:   resourceIntegrationSlackAlertsRead,
		UpdateContext: resourceIntegrationSlackAlertsUpdate,
		DeleteContext: resourceIntegrationSlackAlertsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
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

func resourceIntegrationSlackAlertsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSlackAlertsCreate")
	c := m.(*client.Client)

	resourceData := getSlackAlertsIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/slack", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create integration", fmt.Sprintf("%v", err))
	}

	response := CreateSlackAlertsIntegrationResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceIntegrationSlackAlertsCreate")

	return resourceIntegrationSlackAlertsRead(ctx, d, m)
}

func resourceIntegrationSlackAlertsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSlackAlertsRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := SlackAlertsIntegrationData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("url", response.Url)

	log.Printf("[DEBUG] End resourceIntegrationSlackAlertsRead")

	return diag.Diagnostics{}
}

func resourceIntegrationSlackAlertsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSlackAlertsUpdate")
	c := m.(*client.Client)

	resourceData := getSlackAlertsIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationSlackAlertsUpdate")

	return resourceIntegrationSlackAlertsRead(ctx, d, m)
}

func resourceIntegrationSlackAlertsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSlackAlertsDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/slack/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationSlackAlertsDelete")

	return diag.Diagnostics{}
}

func getSlackAlertsIntegrationDataFromResource(c *client.Client, d *schema.ResourceData) SlackAlertsIntegrationData {
	return SlackAlertsIntegrationData{
		Name: d.Get("name").(string),
		Url:  d.Get("url").(string),
	}
}
