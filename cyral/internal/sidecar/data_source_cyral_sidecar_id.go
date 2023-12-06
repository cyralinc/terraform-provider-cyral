package sidecar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IdentifiedSidecarInfo struct {
	ID      string      `json:"id"`
	Sidecar SidecarData `json:"sidecar"`
}

func DataSourceSidecarID() *schema.Resource {
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
	tflog.Debug(ctx, "Init dataSourceSidecarIDRead")
	c := m.(*client.Client)

	sidecarsInfo, err := ListSidecars(c)
	if err != nil {
		return utils.CreateError("Unable to retrieve the list of existent sidecars.", err.Error())
	}

	sidecarName := d.Get("sidecar_name").(string)
	for _, sidecarInfo := range sidecarsInfo {
		if sidecarName == sidecarInfo.Sidecar.Name {
			d.SetId(sidecarInfo.ID)
			break
		}
	}

	if d.Id() == "" {
		return utils.CreateError("Sidecar not found.",
			fmt.Sprintf("No sidecar found for name '%s'.", sidecarName))
	}

	tflog.Debug(ctx, fmt.Sprintf("Sidecar ID: %s", d.Id()))
	tflog.Debug(ctx, "End dataSourceSidecarIDRead")

	return diag.Diagnostics{}
}

func ListSidecars(c *client.Client) ([]IdentifiedSidecarInfo, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Init listSidecars")
	url := fmt.Sprintf("https://%s/v1/sidecars", c.ControlPlane)
	body, err := c.DoRequest(ctx, url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var sidecarsInfo []IdentifiedSidecarInfo
	if err := json.Unmarshal(body, &sidecarsInfo); err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", sidecarsInfo))
	tflog.Debug(ctx, "End listSidecars")

	return sidecarsInfo, nil
}
