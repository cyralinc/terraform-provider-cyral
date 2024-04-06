package serviceaccount

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
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
	serviceAccount.DisplayName = d.Get(ServiceAccountResourceDisplayNameKey).(string)
	permissionIDs := utils.ConvertFromInterfaceList[string](
		d.Get(ServiceAccountResourcePermissionIDsKey).(*schema.Set).List(),
	)
	if len(permissionIDs) == 0 {
		return fmt.Errorf("at least one permission must be specified for the service account")
	}
	serviceAccount.PermissionIDs = permissionIDs
	return nil
}

func (serviceAccount *ServiceAccount) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(serviceAccount.ClientID)
	d.Set(ServiceAccountResourceDisplayNameKey, serviceAccount.DisplayName)
	d.Set(ServiceAccountResourceClientIDKey, serviceAccount.ClientID)
	isCreateResponse := serviceAccount.ClientSecret != ""
	if isCreateResponse {
		d.Set(ServiceAccountResourceClientSecretKey, serviceAccount.ClientSecret)
	}
	d.Set(ServiceAccountResourcePermissionIDsKey, utils.ConvertToInterfaceList(serviceAccount.PermissionIDs))
	return nil
}
