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

type CreateLookerIntegrationResponse struct {
	ID string `json:"ID"`
}

type LookerIntegrationData struct {
	Name         string `json:"name"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Url          string `json:"url"`
}

func resourceIntegrationLooker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationLookerCreate,
		ReadContext:   resourceIntegrationLookerRead,
		UpdateContext: resourceIntegrationLookerUpdate,
		DeleteContext: resourceIntegrationLookerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceIntegrationLookerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationLookerCreate")
	c := m.(*client.Client)

	resourceData := getLookerIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/looker", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create integration", fmt.Sprintf("%v", err))
	}

	response := CreateLookerIntegrationResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceIntegrationLookerCreate")

	return resourceIntegrationLookerRead(ctx, d, m)
}

func resourceIntegrationLookerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationLookerRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := LookerIntegrationData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("client_id", response.ClientId)
	d.Set("client_secret", response.ClientSecret)
	d.Set("url", response.Url)

	log.Printf("[DEBUG] End resourceIntegrationLookerRead")

	return diag.Diagnostics{}
}

func resourceIntegrationLookerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationLookerUpdate")
	c := m.(*client.Client)

	resourceData := getLookerIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationLookerUpdate")

	return resourceIntegrationLookerRead(ctx, d, m)
}

func resourceIntegrationLookerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationLookerDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/looker/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationLookerDelete")

	return diag.Diagnostics{}
}

func getLookerIntegrationDataFromResource(c *client.Client, d *schema.ResourceData) LookerIntegrationData {
	return LookerIntegrationData{
		Name:         d.Get("name").(string),
		ClientId:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		Url:          d.Get("url").(string),
	}
}
