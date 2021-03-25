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

type CreateELKIntegrationResponse struct {
	ID string `json:"ID"`
}

type ELKIntegrationData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	KibanaURL string `json:"kibanaUrl"`
	ESURL     string `json:"esUrl"`
}

func resourceIntegrationELK() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationELKCreate,
		ReadContext:   resourceIntegrationELKRead,
		UpdateContext: resourceIntegrationELKUpdate,
		DeleteContext: resourceIntegrationELKDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kibana_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"es_url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceIntegrationELKCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationELKCreate")
	c := m.(*client.Client)

	resourceData := getELKIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/elk", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create integration", fmt.Sprintf("%v", err))
	}

	response := CreateELKIntegrationResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceIntegrationELKCreate")

	return resourceIntegrationELKRead(ctx, d, m)
}

func resourceIntegrationELKRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationELKRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := ELKIntegrationData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("kibana_url", response.KibanaURL)
	d.Set("es_url", response.ESURL)

	log.Printf("[DEBUG] End resourceIntegrationELKRead")

	return diag.Diagnostics{}
}

func resourceIntegrationELKUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationELKUpdate")
	c := m.(*client.Client)

	resourceData := getELKIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationELKUpdate")

	return resourceIntegrationELKRead(ctx, d, m)
}

func resourceIntegrationELKDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationELKDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/elk/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationELKDelete")

	return diag.Diagnostics{}
}

func getELKIntegrationDataFromResource(c *client.Client, d *schema.ResourceData) ELKIntegrationData {
	return ELKIntegrationData{
		ID:        d.Id(),
		Name:      d.Get("name").(string),
		KibanaURL: d.Get("kibana_url").(string),
		ESURL:     d.Get("es_url").(string),
	}
}
