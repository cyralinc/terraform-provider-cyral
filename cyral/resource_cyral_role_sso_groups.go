package cyral

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
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
	// IdentityProviderId corresponds to ConnectionName in API.
	IdentityProviderId string `json:"connectionName,omitempty"`
}

type RoleSSOGroupsDeleteRequest struct {
	SSOGroupsMappingsIds []string `json:"mappings,omitempty"`
}

func resourceRoleSSOGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(createRoleSSOGroupsConfig, readRoleSSOGroupsConfig),
		ReadContext:   ReadResource(readRoleSSOGroupsConfig),
		DeleteContext: DeleteResource(deleteRoleSSOGroupsConfig),

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sso_group": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"idp_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

var createRoleSSOGroupsConfig = ResourceOperationConfig{
	Name:       "resourceRoleSSOGroupsCreate",
	HttpMethod: http.MethodPatch,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane,
			d.Get("role_id").(string))
	},
	ResourceData: &RoleSSOGroupsCreateRequest{},
	ResponseData: &RoleSSOGroupsCreateRequest{},
}

var readRoleSSOGroupsConfig = ResourceOperationConfig{
	Name:       "resourceRoleSSOGroupsRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane,
			d.Get("role_id").(string))
	},
	ResponseData: &RoleSSOGroupsReadResponse{},
}

var deleteRoleSSOGroupsConfig = ResourceOperationConfig{
	Name:       "resourceRoleSSOGroupsDelete",
	HttpMethod: http.MethodDelete,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane,
			d.Get("role_id").(string))
	},
	ResourceData: &RoleSSOGroupsDeleteRequest{},
}

func (data RoleSSOGroupsCreateRequest) WriteToSchema(d *schema.ResourceData) {
	d.SetId(fmt.Sprintf("%s/SSOGroups", d.Get("role_id")))
}

func (data *RoleSSOGroupsCreateRequest) ReadFromSchema(d *schema.ResourceData) {
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
}

func (data RoleSSOGroupsCreateRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(data.SSOGroupsMappings)
}

func (data RoleSSOGroupsReadResponse) WriteToSchema(d *schema.ResourceData) {
	var flatSSOGroupsMappings []interface{}

	for _, ssoGroup := range data.SSOGroupsMappings {
		ssoGroupMap := make(map[string]interface{})
		ssoGroupMap["id"] = ssoGroup.Id
		ssoGroupMap["group_name"] = ssoGroup.GroupName
		ssoGroupMap["idp_id"] = ssoGroup.IdentityProviderId

		flatSSOGroupsMappings = append(flatSSOGroupsMappings, ssoGroupMap)
	}

	d.Set("sso_group", flatSSOGroupsMappings)
}

func (data *RoleSSOGroupsReadResponse) ReadFromSchema(d *schema.ResourceData) {}

func (data RoleSSOGroupsDeleteRequest) WriteToSchema(d *schema.ResourceData) {}

func (data *RoleSSOGroupsDeleteRequest) ReadFromSchema(d *schema.ResourceData) {
	var SSOGroupsMappingsIds []string

	if ssoGroups, ok := d.GetOk("sso_group"); ok {
		ssoGroups := ssoGroups.(*schema.Set).List()

		for _, ssoGroup := range ssoGroups {
			ssoGroup := ssoGroup.(map[string]interface{})

			SSOGroupsMappingsIds = append(SSOGroupsMappingsIds, ssoGroup["id"].(string))
		}
	}

	data.SSOGroupsMappingsIds = SSOGroupsMappingsIds
}

func (data RoleSSOGroupsDeleteRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(data.SSOGroupsMappingsIds)
}
