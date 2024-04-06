package ssogroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSchema() *schema.Resource {
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

var createRoleSSOGroupsConfig = core.ResourceOperationConfig{
	ResourceName: "resourceRoleSSOGroupsCreate",
	Type:         operationtype.Create,
	HttpMethod:   http.MethodPatch,
	URLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane,
			d.Get("role_id").(string))
	},
	SchemaReaderFactory: func() core.SchemaReader { return &RoleSSOGroupsCreateRequest{} },
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &RoleSSOGroupsCreateRequest{} },
}

var readRoleSSOGroupsConfig = core.ResourceOperationConfig{
	ResourceName: "resourceRoleSSOGroupsRead",
	Type:         operationtype.Read,
	HttpMethod:   http.MethodGet,
	URLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane,
			d.Get("role_id").(string))
	},
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &RoleSSOGroupsReadResponse{} },
	RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: "Role SSO groups"},
}

var deleteRoleSSOGroupsConfig = core.ResourceOperationConfig{
	ResourceName: "resourceRoleSSOGroupsDelete",
	Type:         operationtype.Delete,
	HttpMethod:   http.MethodDelete,
	URLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane,
			d.Get("role_id").(string))
	},
	SchemaReaderFactory: func() core.SchemaReader { return &RoleSSOGroupsDeleteRequest{} },
}
