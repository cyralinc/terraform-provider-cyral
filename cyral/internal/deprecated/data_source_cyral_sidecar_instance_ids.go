package deprecated

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SidecarDetails struct {
	Instances []DeprecatedSidecarInstances `json:"instances,omitempty"`
}

type DeprecatedSidecarInstances struct {
	ASGInstanceID string `json:"asg_instance,omitempty"`
}

func DataSourceSidecarInstanceIDs() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated. It will be removed in the next major version of " +
			"the provider. Use the data source `cyral_sidecar_instance` instead",
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
		return utils.CreateError(fmt.Sprintf("Unable to retrieve sidecar details. SidecarID: %s",
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
	url := fmt.Sprintf("https://%s/sidecars/%s/details", c.ControlPlane, sidecarID)
	body, err := c.DoRequest(context.Background(), url, http.MethodGet, nil)
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
