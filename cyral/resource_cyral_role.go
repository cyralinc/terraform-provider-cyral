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

type RoleDataRequest struct {
	Name        string   `json:"name"`
	Permissions []string `json:"roles"`
}

type RoleDataResponse struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Permissions []RolePermission `json:"roles"`
}

type RolePermission struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"modify_sidecars_and_repositories": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"modify_users": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"modify_policies": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"view_sidecars_and_repositories": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"view_audit_logs": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"modify_integrations": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"modify_roles": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
	var resourcePermissionsIds []string

	if permissions, ok := d.GetOk("permissions"); ok {
		permissions := permissions.(*schema.Set).List()

		if err := client.ValidateRolePermissions(permissions); err != nil {
			return RoleDataRequest{}, err
		}

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
		Name:        d.Get("name").(string),
		Permissions: resourcePermissionsIds,
	}, nil
}

func flattenPermissions(permissions []RolePermission) []interface{} {
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

func getPermissionsFromAPI(c *client.Client) ([]RolePermission, error) {
	url := fmt.Sprintf("https://%s/v1/users/roles", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return []RolePermission{}, err
	}

	response := RoleDataResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return []RolePermission{}, err
	}

	return response.Permissions, nil
}
