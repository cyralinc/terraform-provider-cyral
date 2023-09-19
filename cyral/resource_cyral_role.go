package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/client"
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
	Permissions []*PermissionInfo `json:"roles,omitempty"`
}

type PermissionInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func resourceRole() *schema.Resource {
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
	log.Printf("[DEBUG] Init resourceRoleCreate")
	c := m.(*client.Client)

	resourceData, err := getRoleDataFromResource(c, d)
	if err != nil {
		return createError("Unable to create role", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/users/groups", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create role", fmt.Sprintf("%v", err))
	}

	response := RoleDataResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.Id)

	log.Printf("[DEBUG] End resourceRoleCreate")

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRoleRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/users/groups/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read role. Role Id: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := RoleDataResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)

	if len(response.Permissions) > 0 {
		flatPermissions := flattenPermissions(response.Permissions)
		log.Printf("[DEBUG] resourceRoleRead - flatPermissions: %s", flatPermissions)

		if err := d.Set("permissions", flatPermissions); err != nil {
			return createError(fmt.Sprintf("Unable to read role. Role Id: %s",
				d.Id()), fmt.Sprintf("%v", err))
		}
	}

	log.Printf("[DEBUG] End resourceRoleRead")

	return diag.Diagnostics{}
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRoleUpdate")
	c := m.(*client.Client)

	resourceData, err := getRoleDataFromResource(c, d)
	if err != nil {
		return createError("Unable to update role", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/users/groups/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update role", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRoleUpdate")

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRoleDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/users/groups/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete role", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRoleDelete")

	return diag.Diagnostics{}
}

func getRoleDataFromResource(c *client.Client, d *schema.ResourceData) (RoleDataRequest, error) {
	var permissionIds []string
	var err error

	if permissionsInterface, ok := d.GetOk("permissions"); ok {
		permissionIds, err = NewPermissionIDsFromInterface(permissionsInterface, c)
		if err != nil {
			return RoleDataRequest{}, err
		}
	}

	return RoleDataRequest{
		Name:          d.Get("name").(string),
		PermissionIDs: permissionIds,
	}, nil
}

func flattenPermissions(permissions []*PermissionInfo) []interface{} {
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

func getPermissionsFromAPI(c *client.Client) ([]*PermissionInfo, error) {
	url := fmt.Sprintf("https://%s/v1/users/roles", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return []*PermissionInfo{}, err
	}

	response := RoleDataResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return []*PermissionInfo{}, err
	}

	return response.Permissions, nil
}
