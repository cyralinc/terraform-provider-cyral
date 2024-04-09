package stats

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
