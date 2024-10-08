package health

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve aggregated information about the " +
			"[sidecar's health](https://cyral.com/docs/sidecars/manage/#check-sidecar-cluster-status), " +
			"considering all instances of the sidecar.",
		ReadContext: core.ReadResource(core.ResourceOperationConfig{
			ResourceName: "SidecarHealthDataSourceRead",
			Type:         operationtype.Read,
			HttpMethod:   http.MethodGet,
			URLFactory: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf(
					"https://%s/v2/sidecars/%s/health", c.ControlPlane, d.Get(utils.SidecarIDKey),
				)
			},
			SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
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
					"[Sidecar Status](https://cyral.com/docs/sidecars/manage/#check-sidecar-cluster-status).",
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
