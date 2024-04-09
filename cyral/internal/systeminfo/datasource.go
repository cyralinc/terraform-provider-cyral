package systeminfo

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
		Description: "Retrieve information from Cyral system.",
		// The DefaultContextHandler is NOT used here as this data source intentionally
		// does not handle 404 errors, returning them to the user.
		ReadContext: core.ReadResource(core.ResourceOperationConfig{
			ResourceName: "SystemInfoDataSourceRead",
			Type:         operationtype.Read,
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
