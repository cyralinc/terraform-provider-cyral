package permission

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Permission struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func permissionsToInterfaceList(permissions []Permission) []any {
	permissionsInterfaceList := make([]any, len(permissions))
	for index, permission := range permissions {
		permissionsInterfaceList[index] = map[string]any{
			utils.IDKey:          permission.Id,
			utils.NameKey:        permission.Name,
			utils.DescriptionKey: permission.Description,
		}
	}
	return permissionsInterfaceList
}

var AllPermissionNames = []string{
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

type PermissionDataSourceResponse struct {
	// Permissions correspond to Roles in API.
	Permissions []Permission `json:"roles"`
}

func (response *PermissionDataSourceResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(uuid.New().String())
	d.Set(PermissionDataSourcePermissionListKey, permissionsToInterfaceList(response.Permissions))
	return nil
}
