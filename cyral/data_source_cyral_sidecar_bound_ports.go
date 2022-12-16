package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/google/uuid"
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
	Address *NetworkAddress `json:"address,omitempty"`
}

type ListComposeBindingsRequest struct {
	PageSize  uint32 `json:"pageSize,omitempty"`
	PageAfter string `json:"pageAfter,omitempty"`
}

func dataSourceSidecarBoundPorts() *schema.Resource {
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
	log.Printf("[DEBUG] Init dataSourceSidecarBoundPortsRead")
	c := m.(*client.Client)

	var boundPorts []uint32

	sidecarID := d.Get("sidecar_id").(string)
	composedBindings, err := getComposedBindings(c, sidecarID)
	if err != nil {
		return createError(fmt.Sprintf("Unable to retrieve repo IDs bound to sidecar. SidecarID: %s",
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

	log.Printf("[DEBUG] Sidecar bound ports: %v", boundPorts)
	log.Printf("[DEBUG] End dataSourceSidecarBoundPortsRead")

	return diag.Diagnostics{}
}

func getComposedBindings(c *client.Client, sidecarID string) ([]*ComposedBinding, error) {
	log.Printf("[DEBUG] Init getComposedBindings")

	var composedBindings []*ComposedBinding
	pageSize := 100
	pageAfter := ""

	for {
		req := &ListComposeBindingsRequest{
			PageSize:  uint32(pageSize),
			PageAfter: pageAfter,
		}

		url := fmt.Sprintf("https://%s/v1/sidecars/%s/composedBindings/filter", c.ControlPlane, sidecarID)
		body, err := c.DoRequest(url, http.MethodPost, req)
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
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", composedBindings)
	log.Printf("[DEBUG] End getComposedBindings")

	return composedBindings, nil
}
