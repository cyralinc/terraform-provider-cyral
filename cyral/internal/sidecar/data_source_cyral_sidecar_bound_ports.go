package sidecar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/listener"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ListComposeBindingsResponse struct {
	ComposedBindings []*ComposedBinding `json:"composedBindings,omitempty"`
	TotalCount       uint32             `json:"totalCount,omitempty"`
}

type ComposedBinding struct {
	Binding   *BindingComponent    `json:"binding,omitempty"`
	Listeners []*ListenerComponent `json:"listeners,omitempty"`
}

type BindingComponent struct {
	Id string `json:"id,omitempty"`
}

type ListenerComponent struct {
	Address *listener.NetworkAddress `json:"address,omitempty"`
}

func DataSourceSidecarBoundPorts() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves all the ports of a given sidecar that are currently bound to repositories.",
		ReadContext: dataSourceSidecarBoundPortsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Computed ID for this data source (locally computed to be used in Terraform state).",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"sidecar_id": {
				Description: "The ID of the sidecar.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"bound_ports": {
				Description: "All the sidecar ports that are currently bound to repositories.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceSidecarBoundPortsRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	tflog.Debug(ctx, "Init dataSourceSidecarBoundPortsRead")
	c := m.(*client.Client)

	var boundPorts []uint32

	sidecarID := d.Get("sidecar_id").(string)
	composedBindings, err := getComposedBindings(ctx, c, sidecarID)
	if err != nil {
		return utils.CreateError(fmt.Sprintf("Unable to retrieve repo IDs bound to sidecar. SidecarID: %s",
			sidecarID), err.Error())
	}

	for _, composedBinding := range composedBindings {
		for _, listener := range composedBinding.Listeners {
			boundPorts = append(boundPorts, uint32(listener.Address.Port))
		}
	}
	// Sorts ports so that we can have a more organized and deterministic output.
	sort.Slice(boundPorts, func(i, j int) bool { return boundPorts[i] < boundPorts[j] })

	d.SetId(uuid.New().String())
	d.Set("bound_ports", boundPorts)

	tflog.Debug(ctx, fmt.Sprintf("Sidecar bound ports: %v", boundPorts))
	tflog.Debug(ctx, "End dataSourceSidecarBoundPortsRead")

	return diag.Diagnostics{}
}

func getComposedBindings(ctx context.Context, c *client.Client, sidecarID string) ([]*ComposedBinding, error) {
	tflog.Debug(ctx, "Init getComposedBindings")

	var composedBindings []*ComposedBinding
	pageSize := 100
	pageAfter := ""

	for {
		url := fmt.Sprintf("https://%s/v1/sidecars/%s/composedBindings/filter?pageSize=%d",
			c.ControlPlane, sidecarID, pageSize)
		if pageAfter != "" {
			url = url + fmt.Sprintf("&pageAfter=%s", pageAfter)
		}
		body, err := c.DoRequest(ctx, url, http.MethodPost, nil)
		if err != nil {
			return nil, err
		}

		var resp ListComposeBindingsResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}
		composedBindings = append(composedBindings, resp.ComposedBindings...)
		if len(composedBindings) < int(resp.TotalCount) {
			pageAfter = resp.ComposedBindings[len(resp.ComposedBindings)-1].Binding.Id
		} else {
			break
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshaled): %#v", composedBindings))
	tflog.Debug(ctx, "End getComposedBindings")

	return composedBindings, nil
}
