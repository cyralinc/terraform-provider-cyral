package ssogroups

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RoleSSOGroupsCreateRequest struct {
	SSOGroupsMappings []*SSOGroup `json:"mappings,omitempty"`
}

type RoleSSOGroupsReadResponse struct {
	SSOGroupsMappings []*SSOGroup `json:"mappings,omitempty"`
}

type SSOGroup struct {
	Id        string `json:"id,omitempty"`
	GroupName string `json:"groupName,omitempty"`
	// IdentityProviderId corresponds to ConnectionAlias in API.
	IdentityProviderId string `json:"connectionAlias,omitempty"`
	// IdentityProviderName corresponds to ConnectionName in API.
	IdentityProviderName string `json:"connectionName,omitempty"`
}

type RoleSSOGroupsDeleteRequest struct {
	SSOGroupsMappingsIds []string `json:"mappings,omitempty"`
}

func (data RoleSSOGroupsCreateRequest) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get("role_id").(string))
	return nil
}

func (data *RoleSSOGroupsCreateRequest) ReadFromSchema(d *schema.ResourceData) error {
	var SSOGroupsMappings []*SSOGroup

	if ssoGroups, ok := d.GetOk("sso_group"); ok {
		ssoGroups := ssoGroups.(*schema.Set).List()

		for _, ssoGroup := range ssoGroups {
			ssoGroup := ssoGroup.(map[string]interface{})

			SSOGroupsMappings = append(SSOGroupsMappings, &SSOGroup{
				GroupName:          ssoGroup["group_name"].(string),
				IdentityProviderId: ssoGroup["idp_id"].(string),
			})
		}
	}

	data.SSOGroupsMappings = SSOGroupsMappings

	return nil
}

func (data RoleSSOGroupsCreateRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(data.SSOGroupsMappings)
}

func (data RoleSSOGroupsReadResponse) WriteToSchema(d *schema.ResourceData) error {
	var flatSSOGroupsMappings []interface{}

	for _, ssoGroup := range data.SSOGroupsMappings {
		ssoGroupMap := make(map[string]interface{})
		ssoGroupMap["id"] = ssoGroup.Id
		ssoGroupMap["group_name"] = ssoGroup.GroupName
		ssoGroupMap["idp_id"] = ssoGroup.IdentityProviderId
		ssoGroupMap["idp_name"] = ssoGroup.IdentityProviderName

		flatSSOGroupsMappings = append(flatSSOGroupsMappings, ssoGroupMap)
	}

	d.Set("sso_group", flatSSOGroupsMappings)

	return nil
}

func (data *RoleSSOGroupsDeleteRequest) ReadFromSchema(d *schema.ResourceData) error {
	var SSOGroupsMappingsIds []string

	if ssoGroups, ok := d.GetOk("sso_group"); ok {
		ssoGroups := ssoGroups.(*schema.Set).List()

		for _, ssoGroup := range ssoGroups {
			ssoGroup := ssoGroup.(map[string]interface{})

			SSOGroupsMappingsIds = append(SSOGroupsMappingsIds, ssoGroup["id"].(string))
		}
	}

	data.SSOGroupsMappingsIds = SSOGroupsMappingsIds

	return nil
}

func (data RoleSSOGroupsDeleteRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(data.SSOGroupsMappingsIds)
}
