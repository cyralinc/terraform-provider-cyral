package role

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
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

func roleSSOGroupsResourceSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_id": {
				Description: "The ID of the role resource that will be configured.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"sso_group": {
				Description: "A block responsible for mapping an SSO group to a role.",
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of an SSO group mapping.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"group_name": {
							Description: "The name of the SSO group to be mapped.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
						"idp_id": {
							Description: "The ID of the identity provider integration to be mapped.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
						"idp_name": {
							Description: "The name of the identity provider integration of an SSO group mapping.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Previously, the ID for cyral_role_sso_groups had the format
// {role_id}/SSOGroups. The goal of this state upgrade is to remove the suffix
// `SSOGroups`.
func UpgradeRoleSSOGroupsV0(
	_ context.Context,
	rawState map[string]interface{},
	_ interface{},
) (map[string]interface{}, error) {
	rawState["id"] = rawState["role_id"]
	return rawState, nil
}

func ResourceRoleSSOGroups() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [mapping SSO groups to specific roles](https://cyral.com/docs/account-administration/acct-manage-cyral-roles/#map-an-sso-group-to-a-cyral-administrator-role) on Cyral control plane. See also: [Role](./role.md).",
		CreateContext: core.CreateResource(createRoleSSOGroupsConfig, readRoleSSOGroupsConfig),
		ReadContext:   core.ReadResource(readRoleSSOGroupsConfig),
		DeleteContext: core.DeleteResource(deleteRoleSSOGroupsConfig),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: roleSSOGroupsResourceSchemaV0().
					CoreConfigSchema().ImpliedType(),
				Upgrade: UpgradeRoleSSOGroupsV0,
			},
		},

		Schema: roleSSOGroupsResourceSchemaV0().Schema,

		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set("role_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

var createRoleSSOGroupsConfig = core.ResourceOperationConfig{
	Name:       "resourceRoleSSOGroupsCreate",
	HttpMethod: http.MethodPatch,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane,
			d.Get("role_id").(string))
	},
	NewResourceData: func() core.SchemaReader { return &RoleSSOGroupsCreateRequest{} },
	NewResponseData: func(_ *schema.ResourceData) core.SchemaWriter { return &RoleSSOGroupsCreateRequest{} },
}

var readRoleSSOGroupsConfig = core.ResourceOperationConfig{
	Name:       "resourceRoleSSOGroupsRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane,
			d.Get("role_id").(string))
	},
	NewResponseData:     func(_ *schema.ResourceData) core.SchemaWriter { return &RoleSSOGroupsReadResponse{} },
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Role SSO groups"},
}

var deleteRoleSSOGroupsConfig = core.ResourceOperationConfig{
	Name:       "resourceRoleSSOGroupsDelete",
	HttpMethod: http.MethodDelete,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane,
			d.Get("role_id").(string))
	},
	NewResourceData: func() core.SchemaReader { return &RoleSSOGroupsDeleteRequest{} },
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
