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

type CreateLogstashIntegrationResponse struct {
	ID string `json:"ID"`
}

type LogstashIntegrationData struct {
	Endpoint                   string `json:"endpoint"`
	Name                       string `json:"name"`
	UseMutualAuthentication    bool   `json:"useMutualAuthentication"`
	UsePrivateCertificateChain bool   `json:"usePrivateCertificateChain"`
	UseTLS                     bool   `json:"useTLS"`
}

func resourceIntegrationLogstash() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationLogstashCreate,
		ReadContext:   resourceIntegrationLogstashRead,
		UpdateContext: resourceIntegrationLogstashUpdate,
		DeleteContext: resourceIntegrationLogstashDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"use_mutual_authentication": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"use_private_certificate_chain": {
				Type:     schema.TypeBool,
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

func resourceIntegrationLogstashCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationLogstashCreate")
	c := m.(*client.Client)

	resourceData := getLogstashIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/logstash", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create integration", fmt.Sprintf("%v", err))
	}

	response := CreateLogstashIntegrationResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceIntegrationLogstashCreate")

	return resourceIntegrationLogstashRead(ctx, d, m)
}

func resourceIntegrationLogstashRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationLogstashRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/logstash/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := LogstashIntegrationData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("endpoint", response.Endpoint)
	d.Set("use_mutual_authentication", response.UseMutualAuthentication)
	d.Set("use_private_certificate_chain", response.UsePrivateCertificateChain)
	d.Set("use_tls", response.UseTLS)

	log.Printf("[DEBUG] End resourceIntegrationLogstashRead")

	return diag.Diagnostics{}
}

func resourceIntegrationLogstashUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationLogstashUpdate")
	c := m.(*client.Client)

	resourceData := getLogstashIntegrationDataFromResource(c, d)

	url := fmt.Sprintf("https://%s/v1/integrations/logstash/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationLogstashUpdate")

	return resourceIntegrationLogstashRead(ctx, d, m)
}

func resourceIntegrationLogstashDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceIntegrationLogstashDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/integrations/logstash/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceIntegrationLogstashDelete")

	return diag.Diagnostics{}
}

func getLogstashIntegrationDataFromResource(c *client.Client, d *schema.ResourceData) LogstashIntegrationData {
	return LogstashIntegrationData{
		Endpoint:                   d.Get("endpoint").(string),
		Name:                       d.Get("name").(string),
		UseMutualAuthentication:    d.Get("use_mutual_authentication").(bool),
		UsePrivateCertificateChain: d.Get("use_private_certificate_chain").(bool),
		UseTLS:                     d.Get("use_tls").(bool),
	}
}
