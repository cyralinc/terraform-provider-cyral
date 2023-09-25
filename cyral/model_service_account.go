package cyral

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServiceAccount struct {
	DisplayName  string `json:"displayName"`
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	// Permissions correspond to Roles in Cyral APIs.
	PermissionIDs []string `json:"roleIds"`
}

func (serviceAccount *ServiceAccount) ReadFromSchema(d *schema.ResourceData) error {
	serviceAccount.DisplayName = d.Get(serviceAccountResourceDisplayNameKey).(string)
	permissionIDs := convertFromInterfaceList[string](
		d.Get(serviceAccountResourcePermissionIDsKey).(*schema.Set).List(),
	)
	if len(permissionIDs) == 0 {
		return fmt.Errorf("at least one permission must be specified for the service account")
	}
	serviceAccount.PermissionIDs = permissionIDs
	return nil
}

func (serviceAccount *ServiceAccount) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(serviceAccount.ClientID)
	d.Set(serviceAccountResourceDisplayNameKey, serviceAccount.DisplayName)
	d.Set(serviceAccountResourceClientIDKey, serviceAccount.ClientID)
	isCreateResponse := serviceAccount.ClientSecret != ""
	if isCreateResponse {
		d.Set(serviceAccountResourceClientSecretKey, serviceAccount.ClientSecret)
	}
	d.Set(serviceAccountResourcePermissionIDsKey, convertToInterfaceList(serviceAccount.PermissionIDs))
	return nil
}
