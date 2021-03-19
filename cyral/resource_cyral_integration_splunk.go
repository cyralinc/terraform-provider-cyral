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

type CreateSplunkIntegrationResponse struct {
	ID string `json:"id"`
}

type SplunkIntegrationData struct {
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
	Port        int    `json:"hecPort,string"`
	Host        string `json:"host"`
	Index       string `json:"index"`
	UseTLS      bool   `json:"useTLS"`
}

func resourceIntegrationSplunk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationSplunkCreate,
		ReadContext:   resourceIntegrationSplunkRead,
		UpdateContext: resourceIntegrationSplunkUpdate,
		DeleteContext: resourceIntegrationSplunkDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"index": {
				Type:     schema.TypeString,
				Required: true,
			},
			"use_tls": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceIntegrationSplunkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSplunkCreate")
	c := m.(*client.Client)

	resourceData := getSplunkIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/splunk", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create integration", fmt.Sprintf("%v", err))
	}

	response := CreateSplunkIntegrationResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceIntegrationSplunkCreate")

	return resourceIntegrationSplunkRead(ctx, d, m)
}

func resourceIntegrationSplunkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSplunkRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := SplunkIntegrationData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("access_token", response.AccessToken)
	d.Set("port", response.Port)
	d.Set("host", response.Host)
	d.Set("index", response.Index)
	d.Set("use_tls", response.UseTLS)

	log.Printf("[DEBUG] End resourceIntegrationSplunkRead")

	return diag.Diagnostics{}
}

func resourceIntegrationSplunkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSplunkUpdate")
	c := m.(*client.Client)

	resourceData := getSplunkIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationSplunkUpdate")

	return resourceIntegrationSplunkRead(ctx, d, m)
}

func resourceIntegrationSplunkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationSplunkDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationSplunkDelete")

	return diag.Diagnostics{}
}

func getSplunkIntegrationDataFromResource(c *client.Client, d *schema.ResourceData) SplunkIntegrationData {
	return SplunkIntegrationData{
		Name:        d.Get("name").(string),
		AccessToken: d.Get("access_token").(string),
		Port:        d.Get("port").(int),
		Host:        d.Get("host").(string),
		Index:       d.Get("index").(string),
		UseTLS:      d.Get("use_tls").(bool),
	}
}
