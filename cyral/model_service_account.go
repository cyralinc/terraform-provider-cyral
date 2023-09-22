package cyral

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServiceAccount struct {
	DisplayName  string `json:"displayName"`
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	// Permissions correspond to Roles in Cyral APIs.
	PermissionIDs []string `json:"roleIds"`
}

func (serviceAccount *ServiceAccount) ReadFromSchema(d *schema.ResourceData, c *client.Client) error {
	serviceAccount.DisplayName = d.Get(serviceAccountResourceDisplayNameKey).(string)
	permissionIDs, err := NewPermissionIDsFromInterface(d.Get(serviceAccountResourcePermissionsKey), c)
	if err != nil {
		return err
	}
	if len(permissionIDs) == 0 {
		return fmt.Errorf("at least one permission must be specified for the service account")
	}
	serviceAccount.PermissionIDs = permissionIDs
	return nil
}

func (serviceAccount *ServiceAccount) WriteToSchema(d *schema.ResourceData, c *client.Client) error {
	d.SetId(serviceAccount.ClientID)
	d.Set(serviceAccountResourceDisplayNameKey, serviceAccount.DisplayName)
	d.Set(serviceAccountResourceClientIDKey, serviceAccount.ClientID)
	isCreateResponse := serviceAccount.ClientSecret != ""
	if isCreateResponse {
		d.Set(serviceAccountResourceClientSecretKey, serviceAccount.ClientSecret)
	}
	permissionsInterfaceList, err := PermissionIDsToInterfaceList(serviceAccount.PermissionIDs, c)
	if err != nil {
		return err
	}
	d.Set(serviceAccountResourcePermissionsKey, permissionsInterfaceList)
	return nil
}

func NewPermissionIDsFromInterface(permissionsInterface any, c *client.Client) ([]string, error) {
	if permissionsInterface == nil {
		return nil, nil
	}
	permissionsList := permissionsInterface.(*schema.Set).List()
	if len(permissionsList) == 0 {
		return nil, nil
	}
	resourcePermissions := permissionsList[0].(map[string]any)
	apiPermissions, err := getPermissionsFromAPI(c)
	if err != nil {
		return nil, fmt.Errorf("error getting permissions from API")
	}
	var resourcePermissionIds []string
	for _, apiPermission := range apiPermissions {
		resourcePermission := resourcePermissions[formatPermissionName(apiPermission.Name)]
		if hasPermission, ok := resourcePermission.(bool); ok && hasPermission {
			resourcePermissionIds = append(resourcePermissionIds, apiPermission.Id)
		}
	}
	return resourcePermissionIds, nil
}

func PermissionIDsToInterfaceList(permissionIDs []string, c *client.Client) ([]any, error) {
	if permissionIDs == nil {
		return nil, nil
	}
	apiPermissions, err := getPermissionsFromAPI(c)
	if err != nil {
		return nil, fmt.Errorf("error getting permissions from API")
	}
	permissionsInterfaceList := make([]any, 1)
	permissionsMap := make(map[string]any)
	for _, permissionID := range permissionIDs {
		for _, apiPermission := range apiPermissions {
			if permissionID == apiPermission.Id {
				permissionsMap[formatPermissionName(apiPermission.Name)] = true
			}
		}
	}
	permissionsInterfaceList[0] = permissionsMap
	return permissionsInterfaceList, nil
}
