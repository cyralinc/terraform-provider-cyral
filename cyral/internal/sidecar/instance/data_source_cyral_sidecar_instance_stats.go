package instance

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Schema keys
	InstanceIDKey        = "instance_id"
	QueriesPerSecondKey  = "queries_per_second"
	ActiveConnectionsKey = "active_connections"
)

type SidecarInstanceStats struct {
	QueriesPerSecond  float32 `json:"queriesPerSecond"`
	ActiveConnections uint32  `json:"activeConnections"`
}

func (stats *SidecarInstanceStats) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get(InstanceIDKey).(string))
	d.Set(QueriesPerSecondKey, stats.QueriesPerSecond)
	d.Set(ActiveConnectionsKey, stats.ActiveConnections)

	return nil
}

func DataSourceSidecarInstanceStats() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve sidecar instance statistics. See also data source " +
			"[`cyral_sidecar_instance`](../data-sources/sidecar_instance.md).",
		ReadContext: core.ReadResource(core.ResourceOperationConfig{
			ResourceName: "SidecarInstanceStatsDataSourceRead",
			Type:         operationtype.Read,
			HttpMethod:   http.MethodGet,
			URLFactory: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf(
					"https://%s/v2/sidecars/%s/instances/%s/stats",
					c.ControlPlane,
					d.Get(utils.SidecarIDKey),
					d.Get(InstanceIDKey),
				)
			},
			SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
				return &SidecarInstanceStats{}
			},
		}),
		Schema: map[string]*schema.Schema{
			utils.SidecarIDKey: {
				Description: "Sidecar identifier.",
				Type:        schema.TypeString,
				Required:    true,
			},
			InstanceIDKey: {
				Description: "Sidecar instance identifier. See also data source " +
					"[`cyral_sidecar_instance`](../data-sources/sidecar_instance.md).",
				Type:     schema.TypeString,
				Required: true,
			},
			utils.IDKey: {
				Description: fmt.Sprintf("Data source identifier. It's equal to `%s`.", InstanceIDKey),
				Type:        schema.TypeString,
				Computed:    true,
			},
			QueriesPerSecondKey: {
				Description: "Amount of queries that the sidecar instance receives per second.",
				Type:        schema.TypeFloat,
				Computed:    true,
			},
			ActiveConnectionsKey: {
				Description: "Number of active connections.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}
