package health

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SidecarHealth struct {
	Status string `json:"status"`
}

func (health *SidecarHealth) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(uuid.New().String())
	d.Set(utils.StatusKey, health.Status)

	return nil
}

func DataSourceSidecarHealth() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve aggregated information about the " +
			"[sidecar's health](https://cyral.com/docs/sidecars/sidecar-manage/#check-sidecar-cluster-status), " +
			"considering all instances of the sidecar.",
		ReadContext: core.ReadResource(core.ResourceOperationConfig{
			Name:       "SidecarHealthDataSourceRead",
			HttpMethod: http.MethodGet,
			CreateURL: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf(
					"https://%s/v2/sidecars/%s/health", c.ControlPlane, d.Get(utils.SidecarIDKey),
				)
			},
			NewResponseData: func(_ *schema.ResourceData) core.ResponseData {
				return &SidecarHealth{}
			},
		}),
		Schema: map[string]*schema.Schema{
			utils.SidecarIDKey: {
				Description: "ID of the Sidecar that will be used to retrieve health information.",
				Type:        schema.TypeString,
				Required:    true,
			},
			utils.IDKey: {
				Description: "Data source identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			utils.StatusKey: {
				Description: "Sidecar health status. Possible values are: `HEALTHY`, `DEGRADED`, `UNHEALTHY` " +
					"and `UNKNOWN`. For more information, see " +
					"[Sidecar Status](https://cyral.com/docs/sidecars/sidecar-manage/#check-sidecar-cluster-status).",
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
