package permission

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Cyral permissions. See also resource " +
			"[`cyral_service_account`](../resources/service_account.md).",
		// The DefaultContextHandler is NOT used here as this data source intentionally
		// does not handle 404 errors, returning them to the user.
		ReadContext: core.ReadResource(
			core.ResourceOperationConfig{
				ResourceName: "PermissionDataSourceRead",
				Type:         operationtype.Read,
				HttpMethod:   http.MethodGet,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/users/roles", c.ControlPlane)
				},
				SchemaWriterFactory: func(d *schema.ResourceData) core.SchemaWriter {
					return &PermissionDataSourceResponse{}
				},
			},
		),
		Schema: map[string]*schema.Schema{
			utils.IDKey: {
				Description: "The data source identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			PermissionDataSourcePermissionListKey: {
				Description: "List of all existing Cyral permissions.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						utils.IDKey: {
							Description: "Permission identifier.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						utils.NameKey: {
							Description: "Permission name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						utils.DescriptionKey: {
							Description: "Permission description.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
