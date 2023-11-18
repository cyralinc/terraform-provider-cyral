package systeminfo

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Schema keys
	ControlPlaneVersionKey  = "control_plane_version"
	SidecarLatestVersionKey = "sidecar_latest_version"
)

type SystemInfo struct {
	ControlPlaneVersion  string `json:"controlPlaneVersion"`
	SidecarLatestVersion string `json:"sidecarLatestVersion"`
}

func (systemInfo *SystemInfo) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(uuid.New().String())
	d.Set(ControlPlaneVersionKey, systemInfo.ControlPlaneVersion)
	d.Set(SidecarLatestVersionKey, systemInfo.SidecarLatestVersion)

	return nil
}

func DataSourceSystemInfo() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve information from Cyral system.",
		ReadContext: core.ReadResource(core.ResourceOperationConfig{
			ResourceName: "SystemInfoDataSourceRead",
			HttpMethod:   http.MethodGet,
			URLFactory: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf("https://%s/v1/systemInfo", c.ControlPlane)
			},
			SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
				return &SystemInfo{}
			},
		}),
		Schema: map[string]*schema.Schema{
			utils.IDKey: {
				Description: "Data source identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			ControlPlaneVersionKey: {
				Description: "Control Plane version.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			SidecarLatestVersionKey: {
				Description: "Latest Sidecar version available to this Control Plane.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
