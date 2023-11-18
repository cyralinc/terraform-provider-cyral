package role

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/permission"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Roles correspond to Groups in API.
type RoleDataRequest struct {
	Name string `json:"name,omitempty"`
	// Permissions correspond to Roles in API.
	PermissionIDs []string `json:"roles,omitempty"`
}

// Roles correspond to Groups in API.
type RoleDataResponse struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// Permissions correspond to Roles in API.
	Permissions []*permission.Permission `json:"roles,omitempty"`
}

func ResourceRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [roles for Cyral control plane users](https://cyral.com/docs/account-administration/acct-manage-cyral-roles/#create-and-manage-administrator-roles-for-cyral-control-plane-users). See also: [Role SSO Groups](./role_sso_groups.md).",
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the role.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"permissions": {
				Description: "A block responsible for configuring the role permissions.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: permissionsSchema,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRoleCreate")
	c := m.(*client.Client)

	resourceData, err := getRoleDataFromResource(c, d)
	if err != nil {
		return utils.CreateError("Unable to create role", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/users/groups", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return utils.CreateError("Unable to create role", fmt.Sprintf("%v", err))
	}

	response := RoleDataResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	d.SetId(response.Id)

	tflog.Debug(ctx, "End resourceRoleCreate")

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRoleRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/users/groups/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return utils.CreateError(fmt.Sprintf("Unable to read role. Role Id: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := RoleDataResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	d.Set("name", response.Name)

	if len(response.Permissions) > 0 {
		flatPermissions := flattenPermissions(response.Permissions)
		tflog.Debug(ctx, fmt.Sprintf("resourceRoleRead - flatPermissions: %s", flatPermissions))

		if err := d.Set("permissions", flatPermissions); err != nil {
			return utils.CreateError(fmt.Sprintf("Unable to read role. Role Id: %s",
				d.Id()), fmt.Sprintf("%v", err))
		}
	}

	tflog.Debug(ctx, "End resourceRoleRead")

	return diag.Diagnostics{}
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRoleUpdate")
	c := m.(*client.Client)

	resourceData, err := getRoleDataFromResource(c, d)
	if err != nil {
		return utils.CreateError("Unable to update role", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/users/groups/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return utils.CreateError("Unable to update role", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, "End resourceRoleUpdate")

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRoleDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/users/groups/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return utils.CreateError("Unable to delete role", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, "End resourceRoleDelete")

	return diag.Diagnostics{}
}

func getRoleDataFromResource(c *client.Client, d *schema.ResourceData) (RoleDataRequest, error) {
	var resourcePermissionsIds []string

	if permissions, ok := d.GetOk("permissions"); ok {
		permissions := permissions.(*schema.Set).List()

		resourcePermissions := permissions[0].(map[string]interface{})

		apiPermissions, err := getPermissionsFromAPI(c)
		if err != nil {
			return RoleDataRequest{}, err
		}

		for _, apiPermission := range apiPermissions {
			resourcePermission := resourcePermissions[formatPermissionName(apiPermission.Name)]
			if v, ok := resourcePermission.(bool); ok && v {
				resourcePermissionsIds = append(resourcePermissionsIds, apiPermission.Id)
			}
		}
	}

	return RoleDataRequest{
		Name:          d.Get("name").(string),
		PermissionIDs: resourcePermissionsIds,
	}, nil
}

func flattenPermissions(permissions []*permission.Permission) []interface{} {
	flatPermissions := make([]interface{}, 1)

	permissionsMap := make(map[string]interface{})
	for _, permission := range permissions {
		permissionsMap[formatPermissionName(permission.Name)] = true
	}

	flatPermissions[0] = permissionsMap

	return flatPermissions
}

func formatPermissionName(permissionName string) string {
	permissionName = strings.ToLower(permissionName)
	permissionName = strings.ReplaceAll(permissionName, " ", "_")
	return permissionName
}

func getPermissionsFromAPI(c *client.Client) ([]*permission.Permission, error) {
	url := fmt.Sprintf("https://%s/v1/users/roles", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return []*permission.Permission{}, err
	}

	response := RoleDataResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return []*permission.Permission{}, err
	}

	return response.Permissions, nil
}
