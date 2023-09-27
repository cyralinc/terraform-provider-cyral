package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Schema keys
	StatusKey = "status"
)

type SidecarHealth struct {
	Status string `json:"status"`
}

func (health *SidecarHealth) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(uuid.New().String())
	d.Set(StatusKey, health.Status)

	return nil
}

func dataSourceSidecarHealth() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve aggregated information about the " +
			"[sidecar's health](https://cyral.com/docs/sidecars/sidecar-manage/#check-sidecar-cluster-status), " +
			"considering all instances of the sidecar.",
		ReadContext: ReadResource(ResourceOperationConfig{
			Name:       "SidecarHealthDataSourceRead",
			HttpMethod: http.MethodGet,
			CreateURL: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf(
					"https://%s/v2/sidecars/%s/health", c.ControlPlane, d.Get(SidecarIDKey),
				)
			},
			NewResponseData: func(_ *schema.ResourceData) ResponseData {
				return &SidecarHealth{}
			},
		}),
		Schema: map[string]*schema.Schema{
			SidecarIDKey: {
				Description: "ID of the Sidecar that will be used to retrieve health information.",
				Type:        schema.TypeString,
				Required:    true,
			},
			IDKey: {
				Description: "Data source identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			StatusKey: {
				Description: "Sidecar health status. Possible values are: `HEALTHY`, `DEGRADED`, `UNHEALTHY` " +
					"and `UNKNOWN`. For more information, see " +
					"[Sidecar Status](https://cyral.com/docs/sidecars/sidecar-manage/#check-sidecar-cluster-status).",
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
