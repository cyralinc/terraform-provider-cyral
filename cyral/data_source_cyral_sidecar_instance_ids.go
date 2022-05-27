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

type SidecarDetails struct {
	Instances []SidecarInstance `json:"instances,omitempty"`
}

type SidecarInstance struct {
	ASGInstanceID string `json:"asg_instance,omitempty"`
}

func dataSourceSidecarInstanceIDs() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves the IDs of all the current instances of a given sidecar.",
		ReadContext: dataSourceSidecarInstanceIDsRead,
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
			"instance_ids": {
				Description: "All the current instance IDs of the sidecar.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceSidecarInstanceIDsRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	log.Printf("[DEBUG] Init dataSourceSidecarInstanceIDsRead")
	c := m.(*client.Client)

	var instanceIDs []string

	sidecarID := d.Get("sidecar_id").(string)
	sidecarDetails, err := getSidecarDetails(c, sidecarID)
	if err != nil {
		return createError(fmt.Sprintf("Unable to retrieve sidecar details. SidecarID: %s",
			sidecarID), err.Error())
	}

	for _, instance := range sidecarDetails.Instances {
		instanceIDs = append(instanceIDs, instance.ASGInstanceID)
	}
	// Sorts instance IDs so that we can have a more organized and deterministic output.
	sort.Strings(instanceIDs)

	d.SetId(uuid.New().String())
	d.Set("instance_ids", instanceIDs)

	log.Printf("[DEBUG] Sidecar instance IDs: %v", instanceIDs)
	log.Printf("[DEBUG] End dataSourceSidecarInstanceIDsRead")

	return diag.Diagnostics{}
}

func getSidecarDetails(c *client.Client, sidecarID string) (SidecarDetails, error) {
	log.Printf("[DEBUG] Init getSidecarDetails")
	// Remove port from control plane to make request to Jeeves server
	controlPlaneWithoutPort := removePortFromURL(c.ControlPlane)
	url := fmt.Sprintf("https://%s/sidecars/%s/details", controlPlaneWithoutPort, sidecarID)
	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return SidecarDetails{}, err
	}

	var sidecarDetails SidecarDetails
	if err := json.Unmarshal(body, &sidecarDetails); err != nil {
		return SidecarDetails{}, err
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", sidecarDetails)
	log.Printf("[DEBUG] End getSidecarDetails")

	return sidecarDetails, nil
}
