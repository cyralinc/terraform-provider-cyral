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
	Description string   `json:"description"`
	Permissions []string `json:"roles"`
}

type RoleDataResponse struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
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
			"description": {
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

	setRoleDataToResource(d, response)

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
	permissionsData, err := getPermissionsFromAPI(c)
	if err != nil {
		return RoleDataRequest{}, err
	}

	var resourcePermissionsIds []string
	for _, permission := range permissionsData {
		resourcePermission := d.Get(formatPermissionName(permission.Name))
		if resourcePermission == true {
			resourcePermissionsIds = append(resourcePermissionsIds, permission.Id)
		}
	}

	return RoleDataRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Permissions: resourcePermissionsIds,
	}, nil
}

func setRoleDataToResource(d *schema.ResourceData, roleData RoleDataResponse) {
	d.Set("name", roleData.Name)
	d.Set("description", roleData.Description)

	for _, permission := range roleData.Permissions {
		d.Set(formatPermissionName(permission.Name), true)
	}
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
