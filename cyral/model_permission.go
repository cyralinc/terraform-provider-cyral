package cyral

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Permission struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var permissionNames = []string{
	"Approval Management",
	"Modify Policies",
	"Modify Roles",
	"Modify Sidecars and Repositories",
	"Modify Users",
	"Repo Crawler",
	"View Audit Logs",
	"View Datamaps",
	"View Integrations",
	"View Policies",
	"View Roles",
	"View Users",
	"Modify Integrations",
}

const (
	// Schema keys
	approvalManagementPermissionKey           = "approval_management"
	modifyPoliciesPermissionKey               = "modify_policies"
	modifyRolesPermissionKey                  = "modify_roles"
	modifySidecarAndRepositoriesPermissionKey = "modify_sidecars_and_repositories"
	modifyUsersPermissionKey                  = "modify_users"
	repoCrawlerPermissionKey                  = "repo_crawler"
	viewAuditLogsPermissionKey                = "view_audit_logs"
	viewDatamapsPermissionKey                 = "view_datamaps"
	viewIntegrationsPermissionKey             = "view_integrations"
	viewPoliciesPermissionKey                 = "view_policies"
	viewRolesPermissionKey                    = "view_roles"
	viewUsersPermissionKey                    = "view_users"
	modifyIntegrationsPermissionKey           = "modify_integrations"
)

var permissionsSchema = map[string]*schema.Schema{
	approvalManagementPermissionKey: {
		Description: "Allows approving or denying approval requests on Cyral Control Plane. " +
			"Defaults to `false`.",
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	modifyPoliciesPermissionKey: {
		Description: "Allows modifying policies on Cyral Control Plane. Defaults to `false`.",
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
	modifySidecarAndRepositoriesPermissionKey: {
		Description: "Allows modifying sidecars and repositories on Cyral Control Plane. Defaults to `false`.",
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
	repoCrawlerPermissionKey: {
		Description: "Allows running the Cyral repo crawler data classifier and user discovery. " +
			"Defaults to `false`.",
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	viewAuditLogsPermissionKey: {
		Description: "Allows viewing audit logs on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	viewDatamapsPermissionKey: {
		Description: "Allows viewing datamaps on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	viewIntegrationsPermissionKey: {
		Description: "Allows viewing integrations on Cyral Control Plane. Defaults to `false`.",
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
	viewRolesPermissionKey: {
		Description: "Allows viewing roles on Cyral Control Plane. Defaults to `false`.",
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
	modifyIntegrationsPermissionKey: {
		Description: "Allows modifying integrations on Cyral Control Plane. Defaults to `false`.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
}
