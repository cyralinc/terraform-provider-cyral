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

type CreateSidecarCredentialsRequest struct {
	SidecarID string `json:"sidecarId"`
}

type SidecarCredentialsData struct {
	SidecarID    string `json:"sidecarId"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func resourceSidecarCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSidecarCredentialsCreate,
		ReadContext:   resourceSidecarCredentialsRead,
		DeleteContext: resourceSidecarCredentialsDelete,

		Schema: map[string]*schema.Schema{
			"sidecar_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceSidecarCredentialsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceSidecarCredentialsCreate")
	c := m.(*client.Client)

	payload := CreateSidecarCredentialsRequest{d.Get("sidecar_id").(string)}

	url := fmt.Sprintf("https://%s/v1/users/sidecarAccounts", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, payload)
	if err != nil {
		return createError("Unable to create sidecar credentials", fmt.Sprintf("%v", err))
	}

	response := SidecarCredentialsData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ClientID)
	d.Set("client_id", response.ClientID)
	d.Set("client_secret", response.ClientSecret) // TODO: encrypt client secret

	return resourceSidecarCredentialsRead(ctx, d, m)
}

func resourceSidecarCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceSidecarCredentialsRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/users/sidecarAccounts/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read sidecar credentials. ClientID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := SidecarCredentialsData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("sidecar_id", response.SidecarID)
	d.Set("client_id", response.ClientID)

	log.Printf("[DEBUG] End resourceSidecarCredentialsRead")

	return diag.Diagnostics{}
}

func resourceSidecarCredentialsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceSidecarCredentialsDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/users/sidecarAccounts/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete sidecar credentials", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceSidecarCredentialsDelete")

	return diag.Diagnostics{}
}
