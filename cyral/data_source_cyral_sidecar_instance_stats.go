package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Schema keys
	InstanceIDKey        = "instance_id"
	QueriesPerSecondKey  = "queries_per_second"
	ActiveConnectionsKey = "active_connections"
)

type SidecarInstanceStats struct {
	QueriesPerSecond  string `json:"queriesPerSecond"`
	ActiveConnections string `json:"activeConnections"`
}

func (stats *SidecarInstanceStats) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get(InstanceIDKey).(string))
	d.Set(QueriesPerSecondKey, stats.QueriesPerSecond)
	d.Set(ActiveConnectionsKey, stats.ActiveConnections)

	return nil
}

func dataSourceSidecarInstanceStats() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve sidecar instance statistics. See also data source " +
			"[`cyral_sidecar_instance`](../data-sources/sidecar_instance.md).",
		ReadContext: ReadResource(ResourceOperationConfig{
			Name:       "SidecarInstanceStatsDataSourceRead",
			HttpMethod: http.MethodGet,
			CreateURL: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf(
					"https://%s/v2/sidecars/%s/instances/%s/stats",
					c.ControlPlane,
					d.Get(SidecarIDKey),
					d.Get(InstanceIDKey),
				)
			},
			NewResponseData: func(_ *schema.ResourceData) ResponseData {
				return &SidecarInstanceStats{}
			},
		}),
		Schema: map[string]*schema.Schema{
			SidecarIDKey: {
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
			IDKey: {
				Description: fmt.Sprintf("Data source identifier. It's equal to `%s`.", InstanceIDKey),
				Type:        schema.TypeString,
				Computed:    true,
			},
			QueriesPerSecondKey: {
				Description: "Amount of queries that the sidecar instance receives per second.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			ActiveConnectionsKey: {
				Description: "Number of active connections.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
