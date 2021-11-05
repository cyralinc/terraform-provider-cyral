package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RoleSSOGroupsData struct {
	Id             string `json:"id"`
	ConnectionName string `json:"connectionName"`
	GroupName      string `json:"groupName"`
}

type RoleSSOGroupsResponse struct {
	Mappings []RoleSSOGroupsData `json:"mappings"`
	Code     string              `json:"code"`
	Error    string              `json:"error"`
	Message  string              `json:"message"`
}

func resourceRoleSSOGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleSSOGroupsCreate,
		ReadContext:   resourceRoleSSOGroupsRead,
		DeleteContext: resourceRoleSSOGroupsDelete,

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sso_group": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"idp_id": {
							Type:     schema.TypeString,
							Required: true,
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

func resourceRoleSSOGroupsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceRoleSSOGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRoleSSOGroupsRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/users/groups/%s/mappings", c.ControlPlane, d.Get("groupId").(string))

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read Role SSO Groups. groupId: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := RoleSSOGroupsResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	log.Printf("[DEBUG] End resourceSidecarRead")

	return diag.Diagnostics{}
}

func resourceRoleSSOGroupsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}
