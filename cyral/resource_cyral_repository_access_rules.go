package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AccessRulesIdentity struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type AccessRulesConfig struct {
	AuthorizationPolicyInstanceIDs []string `json:"authorizationPolicyInstanceIDs"`
}

type AccessRule struct {
	Identity   *AccessRulesIdentity `json:"identity"`
	ValidFrom  *string              `json:"validFrom"`
	ValidUntil *string              `json:"validUntil"`
	Config     *AccessRulesConfig   `json:"config"`
}

type AccessRulesResource struct {
	AccessRules []*AccessRule `json:"accessRules"`
}

type AccessRulesResponse struct {
	AccessRules []*AccessRule `json:"accessRules"`
}

// WriteToSchema is used when reading a resource. It takes whatever the API
// read call returned and translates it into the Terraform schema.
func (arr *AccessRulesResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(
		marshalComposedID(
			[]string{
				d.Get("repository_id").(string),
				d.Get("user_account_id").(string),
			},
			"/",
		),
	)
	// We'll have to build the access rule set in the format expected by Terraform,
	// which boils down to doing a bunch of type casts
	rules := make([]interface{}, 0, len(arr.AccessRules))
	for _, rule := range arr.AccessRules {
		m := make(map[string]interface{})

		m["identity"] = []interface{}{
			map[string]interface{}{
				"type": rule.Identity.Type,
				"name": rule.Identity.Name,
			},
		}

		m["valid_from"] = rule.ValidFrom
		m["valid_until"] = rule.ValidUntil

		if rule.Config != nil && len(rule.Config.AuthorizationPolicyInstanceIDs) > 0 {
			m["config"] = []interface{}{
				map[string]interface{}{
					"policy_ids": rule.Config.AuthorizationPolicyInstanceIDs,
				},
			}
		}

		rules = append(rules, m)
	}
	return d.Set("rule", rules)
}

// ReadFromSchema is called when *creating* or *updating* a resource.
// Essentially, it translates the stuff from the .tf file into whatever the
// API expects. The `AccessRulesResource` will be marshalled verbatim, so
// make sure that it matches *exactly* what the API needs.
func (arr *AccessRulesResource) ReadFromSchema(d *schema.ResourceData) error {
	rules := d.Get("rule").([]interface{})
	var accessRules []*AccessRule

	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})

		accessRule := &AccessRule{}

		identity := ruleMap["identity"].(*schema.Set).List()[0].(map[string]interface{})
		accessRule.Identity = &AccessRulesIdentity{
			Type: identity["type"].(string),
			Name: identity["name"].(string),
		}

		validFrom := ruleMap["valid_from"].(string)
		if validFrom != "" {
			accessRule.ValidFrom = &validFrom
		}

		validUntil := ruleMap["valid_until"].(string)
		if validUntil != "" {
			accessRule.ValidUntil = &validUntil
		}

		conf, ok := ruleMap["config"]
		if ok {
			confList := conf.(*schema.Set).List()
			if len(confList) > 0 {
				config := confList[0].(map[string]interface{})
				policyIDs := config["policy_ids"].([]interface{})
				ids := make([]string, 0, len(policyIDs))
				for _, policyID := range policyIDs {
					ids = append(ids, policyID.(string))
				}
				accessRule.Config = &AccessRulesConfig{
					AuthorizationPolicyInstanceIDs: ids,
				}
			}
		}

		accessRules = append(accessRules, accessRule)
	}

	arr.AccessRules = accessRules
	return nil
}

var ReadRepositoryAccessRulesConfig = ResourceOperationConfig{
	Name:       "RepositoryAccessRulesRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/repos/%s/userAccounts/%s/accessRules",
			c.ControlPlane,
			d.Get("repository_id").(string),
			d.Get("user_account_id").(string),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData {
		return &AccessRulesResponse{}
	},
}

func resourceRepositoryAccessRules() *schema.Resource {
	return &schema.Resource{
		Description: "Manage access rules",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RepositoryAccessRulesCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					repoID := d.Get("repository_id").(string)
					userAccountID := d.Get("user_account_id").(string)
					return fmt.Sprintf("https://%s/v1/repos/%s/userAccounts/%s/accessRules",
						c.ControlPlane,
						repoID,
						userAccountID,
					)
				},
				NewResourceData: func() ResourceData {
					return &AccessRulesResource{}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData {
					return &AccessRulesResponse{}
				},
			},
			ReadRepositoryAccessRulesConfig,
		),
		ReadContext: ReadResource(ReadRepositoryAccessRulesConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RepositoryAccessRulesUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/repos/%s/userAccounts/%s/accessRules",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Get("user_account_id").(string),
					)
				},
				NewResourceData: func() ResourceData {
					return &AccessRulesResource{}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData {
					return &AccessRulesResponse{}
				},
			},
			ReadRepositoryAccessRulesConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RepositoryAccessRulesDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {

					idPieces, err := unmarshalComposedID(d.Id(), "/", 2)
					if err != nil {
						panic(fmt.Sprintf("Failed to unmarshal access rules ID: %v", err))
					}
					repoID := idPieces[0]
					userAccountID := idPieces[1]

					return fmt.Sprintf("https://%s/v1/repos/%s/userAccounts/%s/accessRules",
						c.ControlPlane,
						repoID,
						userAccountID,
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment.",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"repository_id": {
				Description: "ID of the repository.",
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
			},

			"user_account_id": {
				Description: "ID of the database account. This should be the attribute `user_account_id` " +
					"of the resource `cyral_repository_user_account`.",
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},

			"rule": {
				Description: "An ordered list of access rules.",
				Required:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"identity": {
							Description: "The identity of the person/group getting access.",
							Required:    true,
							Type:        schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description: "Identity type. List of supported values: " +
											supportedTypesMarkdown([]string{
												"username",
												"email",
												"group",
											}),
										Required: true,
										Type:     schema.TypeString,
									},

									"name": {
										Description: "The name of the person/group getting access.",
										Required:    true,
										Type:        schema.TypeString,
									},
								},
							},
							MinItems: 1,
							MaxItems: 1,
						},

						"valid_from": {
							Description: "The start time for the grant. Format is: " +
								"`yyyy-mm-ddThh:mm:ssZ`. Eg. `2022-01-24T18:30:00Z`.",
							Optional: true,
							Type:     schema.TypeString,
						},

						"valid_until": {
							Description: "The end time for the grant. Format is: " +
								"`yyyy-mm-ddThh:mm:ssZ`. Eg. `2022-01-24T18:30:00Z`.",
							Optional: true,
							Type:     schema.TypeString,
						},

						"config": {
							Description: "Extra (optional) configuration parameters.",
							Optional:    true,
							MaxItems:    1,
							Type:        schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_ids": {
										Description: "Extra authorization policies, such as PagerDuty or DUO.",
										Required:    true,
										Type:        schema.TypeList,
										MinItems:    1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := unmarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to unmarshal ID: %v",
						err,
					)
				}
				err = d.Set("repository_id", ids[0])
				if err != nil {
					return nil, fmt.Errorf(
						"failed to set 'repository_id': %v",
						err,
					)
				}
				err = d.Set("user_account_id", ids[1])
				if err != nil {
					return nil, fmt.Errorf(
						"failed to set 'user_account_id: %v",
						err,
					)
				}
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
