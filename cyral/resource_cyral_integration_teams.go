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

type CreateMsTeamsIntegrationResponse struct {
	ID string `json:"id"`
}

type MsTeamsIntegrationData struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func resourceIntegrationMsTeams() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationMsTeamsCreate,
		ReadContext:   resourceIntegrationMsTeamsRead,
		UpdateContext: resourceIntegrationMsTeamsUpdate,
		DeleteContext: resourceIntegrationMsTeamsDelete,

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

func resourceIntegrationMsTeamsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationTeamsCreate")
	c := m.(*client.Client)

	resourceData := getMsTeamsIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/teams", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create integration", fmt.Sprintf("%v", err))
	}

	response := CreateMsTeamsIntegrationResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceIntegrationTeamsCreate")

	return resourceIntegrationMsTeamsRead(ctx, d, m)
}

func resourceIntegrationMsTeamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationTeamsRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := MsTeamsIntegrationData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("url", response.URL)

	log.Printf("[DEBUG] End resourceIntegrationTeamsRead")

	return diag.Diagnostics{}
}

func resourceIntegrationMsTeamsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationTeamsUpdate")
	c := m.(*client.Client)

	resourceData := getMsTeamsIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationTeamsUpdate")

	return resourceIntegrationMsTeamsRead(ctx, d, m)
}

func resourceIntegrationMsTeamsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationTeamsDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationTeamsDelete")

	return diag.Diagnostics{}
}

func getMsTeamsIntegrationDataFromResource(c *client.Client, d *schema.ResourceData) MsTeamsIntegrationData {
	return MsTeamsIntegrationData{
		Name: d.Get("name").(string),
		URL:  d.Get("url").(string),
	}
}
