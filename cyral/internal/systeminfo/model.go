package systeminfo

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
