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

type CreateSidecarCredentialsResponse struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func resourceSidecarCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSidecarCredentialsCreate,
		ReadContext:   resourceSidecarCredentialsRead,
		UpdateContext: resourceSidecarCredentialsUpdate,
		DeleteContext: resourceSidecarCredentialsDelete,

		Schema: map[string]*schema.Schema{
			"sidecar_id": {
				Type:     schema.TypeString,
				Required: true,
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

	response := CreateSidecarCredentialsResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ClientID)

	return diag.Diagnostics{}
}

func resourceSidecarCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceSidecarCredentialsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceSidecarCredentialsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}
