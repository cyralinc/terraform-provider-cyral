package role

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/permission"
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
