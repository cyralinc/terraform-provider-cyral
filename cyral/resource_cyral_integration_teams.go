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

type CreateTeamsIntegrationResponse struct {
	ID string `json:"id"`
}

type TeamsIntegrationData struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func resourceIntegrationTeams() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationTeamsCreate,
		ReadContext:   resourceIntegrationTeamsRead,
		UpdateContext: resourceIntegrationTeamsUpdate,
		DeleteContext: resourceIntegrationTeamsDelete,

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

func resourceIntegrationTeamsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationTeamsCreate")
	c := m.(*client.Client)

	resourceData := getTeamsIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/teams", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create integration", fmt.Sprintf("%v", err))
	}

	response := CreateTeamsIntegrationResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceIntegrationTeamsCreate")

	return resourceIntegrationTeamsRead(ctx, d, m)
}

func resourceIntegrationTeamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationTeamsRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := TeamsIntegrationData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("url", response.Url)

	log.Printf("[DEBUG] End resourceIntegrationTeamsRead")

	return diag.Diagnostics{}
}

func resourceIntegrationTeamsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationTeamsUpdate")
	c := m.(*client.Client)

	resourceData := getTeamsIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationTeamsUpdate")

	return resourceIntegrationTeamsRead(ctx, d, m)
}

func resourceIntegrationTeamsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationTeamsDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/notifications/teams/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationTeamsDelete")

	return diag.Diagnostics{}
}

func getTeamsIntegrationDataFromResource(c *client.Client, d *schema.ResourceData) TeamsIntegrationData {
	return TeamsIntegrationData{
		Name: d.Get("name").(string),
		Url:  d.Get("url").(string),
	}
}
