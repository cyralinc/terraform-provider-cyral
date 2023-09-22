package cyral

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

const (
	// Schema keys
	PermissionDataSourcePermissionNamesKey = "permission_names"
	PermissionDataSourcePermissionListKey  = "permission_list"
)

type PermissionDataSourceResponse struct {
	// Permissions correspond to Roles in API.
	Permissions []Permission `json:"roles"`
}

func (response *PermissionDataSourceResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(uuid.New().String())
	d.Set(PermissionDataSourcePermissionListKey, getFilteredPermissionsInterfaceList(d, response.Permissions))
	return nil
}

func getFilteredPermissionsInterfaceList(d *schema.ResourceData, permissions []Permission) []any {
	var filteredPermissionsInterfaceList []any
	permissionNamesFilterInterfaceList := d.Get(PermissionDataSourcePermissionNamesKey).([]any)
	for _, permissionNameFilter := range permissionNamesFilterInterfaceList {
		permissionNameFilter := permissionNameFilter.(string)
		for _, permission := range permissions {
			if permission.Name == permissionNameFilter {
				filteredPermissionsInterfaceList = append(filteredPermissionsInterfaceList, map[string]any{
					IDKey:          permission.Id,
					NameKey:        permission.Name,
					DescriptionKey: permission.Description,
				})
			}
		}
	}
	return filteredPermissionsInterfaceList
}

func dataSourcePermission() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter Cyral permissions. See also resources " +
			"[`cyral_role`](../resources/role.md) and [`cyral_service_account`](../resources/service_account.md).",
		ReadContext: ReadResource(
			ResourceOperationConfig{
				Name:       "PermissionDataSourceRead",
				HttpMethod: http.MethodGet,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/users/roles", c.ControlPlane)
				},
				NewResponseData: func(d *schema.ResourceData) ResponseData {
					return &PermissionDataSourceResponse{}
				},
			},
		),
		Schema: map[string]*schema.Schema{
			IDKey: {
				Description: "The data source identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			PermissionDataSourcePermissionNamesKey: {
				Description: "Filter to retrieve only the permissions that match any of the names present in this " +
					"list. Valid values are: " + supportedTypesMarkdown(permissionNames),
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(permissionNames, false),
				},
			},
			PermissionDataSourcePermissionListKey: {
				Description: "List of existing Cyral permissions satisfying the filter criteria.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						IDKey: {
							Description: "Permission identifier.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						NameKey: {
							Description: "Permission name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						DescriptionKey: {
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
