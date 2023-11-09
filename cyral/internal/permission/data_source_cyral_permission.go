package permission

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

const (
	// Schema keys
	PermissionDataSourcePermissionListKey = "permission_list"
)

type PermissionDataSourceResponse struct {
	// Permissions correspond to Roles in API.
	Permissions []Permission `json:"roles"`
}

func (response *PermissionDataSourceResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(uuid.New().String())
	d.Set(PermissionDataSourcePermissionListKey, permissionsToInterfaceList(response.Permissions))
	return nil
}

func DataSourcePermission() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Cyral permissions. See also resource " +
			"[`cyral_service_account`](../resources/service_account.md).",
		ReadContext: core.ReadResource(
			core.ResourceOperationConfig{
				Name:       "PermissionDataSourceRead",
				HttpMethod: http.MethodGet,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/users/roles", c.ControlPlane)
				},
				NewResponseData: func(d *schema.ResourceData) core.ResponseData {
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