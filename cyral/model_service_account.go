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

const (
	// Schema keys
	modifySidecarAndRepositoriesPermissionKey = "modify_sidecars_and_repositories"
	modifyPoliciesPermissionKey               = "modify_policies"
	modifyIntegrationsPermissionKey           = "modify_integrations"
	modifyUsersPermissionKey                  = "modify_users"
	modifyRolesPermissionKey                  = "modify_roles"
	viewUsersPermissionKey                    = "view_users"
	viewAuditLogsPermissionKey                = "view_audit_logs"
	repoCrawlerPermissionKey                  = "repo_crawler"
	viewDatamapsPermissionKey                 = "view_datamaps"
	viewRolesPermissionKey                    = "view_roles"
	viewPoliciesPermissionKey                 = "view_policies"
	approvalManagementPermissionKey           = "approval_management"
	viewIntegrationsPermissionKey             = "view_integrations"
)

var permissionsSchema = map[string]*schema.Schema{
	modifySidecarAndRepositoriesPermissionKey: {
		Description: "Allows modifying sidecars and repositories on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	modifyPoliciesPermissionKey: {
		Description: "Allows modifying policies on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	modifyIntegrationsPermissionKey: {
		Description: "Allows modifying integrations on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	modifyUsersPermissionKey: {
		Description: "Allows modifying users on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	modifyRolesPermissionKey: {
		Description: "Allows modifying roles on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	viewUsersPermissionKey: {
		Description: "Allows viewing users on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	viewAuditLogsPermissionKey: {
		Description: "Allows viewing audit logs on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	repoCrawlerPermissionKey: {
		Description: "Allows running the Cyral repo crawler data classifier and user discovery. " +
			"Defaults to `false`.",
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	viewDatamapsPermissionKey: {
		Description: "Allows viewing datamaps on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	viewRolesPermissionKey: {
		Description: "Allows viewing roles on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	viewPoliciesPermissionKey: {
		Description: "Allows viewing policies on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	approvalManagementPermissionKey: {
		Description: "Allows approving or denying approval requests on Cyral Control Plane. " +
			"Defaults to `false`.",
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	viewIntegrationsPermissionKey: {
		Description: "Allows viewing integrations on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
}
