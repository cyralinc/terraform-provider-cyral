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

type IdentifiedSidecarInfo struct {
	ID      string      `json:"id"`
	Sidecar SidecarData `json:"sidecar"`
}

func dataSourceSidecarID() *schema.Resource {
	return &schema.Resource{
		Description: "Given a sidecar name, retrieves the respective sidecar ID.",
		ReadContext: dataSourceSidecarIDRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the sidecar.",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"sidecar_name": {
				Description: "The name of the sidecar.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceSidecarIDRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	log.Printf("[DEBUG] Init dataSourceSidecarIDRead")
	c := m.(*client.Client)

	sidecarsInfo, err := listSidecars(c)
	if err != nil {
		return createError("Unable to retrieve the list of existent sidecars.", err.Error())
	}

	sidecarName := d.Get("sidecar_name").(string)
	for _, sidecarInfo := range sidecarsInfo {
		if sidecarName == sidecarInfo.Sidecar.Name {
			d.SetId(sidecarInfo.ID)
			break
		}
	}

	if d.Id() == "" {
		return createError("Sidecar not found.",
			fmt.Sprintf("No sidecar found for name '%s'.", sidecarName))
	}

	log.Printf("[DEBUG] Sidecar ID: %s", d.Id())
	log.Printf("[DEBUG] End dataSourceSidecarIDRead")

	return diag.Diagnostics{}
}

func listSidecars(c *client.Client) ([]IdentifiedSidecarInfo, error) {
	log.Printf("[DEBUG] Init listSidecars")
	url := fmt.Sprintf("https://%s/v1/sidecars", c.ControlPlane)
	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var sidecarsInfo []IdentifiedSidecarInfo
	if err := json.Unmarshal(body, &sidecarsInfo); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", sidecarsInfo)
	log.Printf("[DEBUG] End listSidecars")

	return sidecarsInfo, nil
}
