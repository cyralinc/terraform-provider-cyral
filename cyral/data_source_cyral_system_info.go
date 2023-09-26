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

func dataSourceSystemInfo() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve information from Cyral systems.",
		ReadContext: ReadResource(ResourceOperationConfig{
			Name:       "SystemInfoDataSourceRead",
			HttpMethod: http.MethodGet,
			CreateURL: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf("https://%s/v1/systemInfo", c.ControlPlane)
			},
			NewResponseData: func(_ *schema.ResourceData) ResponseData {
				return &SystemInfo{}
			},
		}),
		Schema: map[string]*schema.Schema{
			IDKey: {
				Description: "The data source identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			ControlPlaneVersionKey: {
				Description: "The Control Plane version.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			SidecarLatestVersionKey: {
				Description: "The latest Sidecar version available to this Control Plane.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
