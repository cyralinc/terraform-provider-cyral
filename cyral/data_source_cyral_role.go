package cyral

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type GetUserGroupsResponse struct {
	Groups []*UserGroup `json:"groups,omitempty"`
}

func (resp *GetUserGroupsResponse) WriteToSchema(d *schema.ResourceData) error {
	roleList := []interface{}{}
	for _, group := range resp.Groups {
		if group != nil {
			argumentVals := map[string]interface{}{
				"id":          group.ID,
				"name":        group.Name,
				"description": group.Description,
				"roles":       group.Roles,
				"members":     group.Members,
			}
			ssoGroups := []interface{}{}
			for _, mapping := range group.Mappings {
				if mapping == nil {
					continue
				}
				ssoGroups = append(ssoGroups, map[string]interface{}{
					"id":         mapping.Id,
					"group_name": mapping.GroupName,
					"idp_id":     mapping.IdentityProviderId,
					"idp_name":   mapping.IdentityProviderName,
				})
			}
			argumentVals["sso_groups"] = ssoGroups
			roleList = append(roleList, argumentVals)
		}
	}
	if err := d.Set("role_list", roleList); err != nil {
		return err
	}
	d.SetId(uuid.New().String())
	return nil
}

type UserGroup struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Roles       []string    `json:"roles,omitempty"`
	Members     []string    `json:"members"`
	Mappings    []*SSOGroup `json:"mappings"`
}

func dataSourceRoleReadConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "RoleDataSourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/users/groups", c.ControlPlane)
		},
		NewResponseData: func() ResponseData { return &GetUserGroupsResponse{} },
	}
}

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter SSO groups.",
		ReadContext: ReadResource(dataSourceRoleReadConfig()),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Filter the results by a regular expression (regex) that matches names of existing user roles.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"role_list": {
				Description: "List of existing roles satisfying given filter criteria.",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "ID of the role in the Cyral environment.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Role name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "Role description.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"roles": {
							Description: "IDs of the permission roles this role is allowed to assume.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"members": {
							Description: "IDs of the users that belong to this role.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"sso_groups": {
							Description: `SSO groups mapped to this role. An SSO group mapping means that this role was automatically granted to a user because there's a rule such as "If a user is an 'Engineer' (SSO group) in a specific Identity Provider, make them a 'Super Admin' (role) in Cyral".`,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "The ID of the SSO group mapping.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"group_name": {
										Description: "The name of the SSO group mapping.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"idp_id": {
										Description: "ID of the identity provider integration.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"idp_name": {
										Description: "Display name of the identity provider integration.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
