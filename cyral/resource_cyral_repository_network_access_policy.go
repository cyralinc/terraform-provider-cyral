package cyral

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	repositoryNetworkAccessPolicyURLFormat = "https://%s/v1/repos/%s/networkAccessPolicy"
)

type NetworkAccessPolicyUpsertResp struct {
	Policy NetworkAccessPolicy `json:"policy"`
}

func (resp NetworkAccessPolicyUpsertResp) WriteToSchema(d *schema.ResourceData) error {
	return resp.Policy.WriteToSchema(d)
}

type NetworkAccessPolicy struct {
	NetworkAccessRules `json:"networkAccessRules,omitempty"`
}

type NetworkAccessRules struct {
	Rules []NetworkAccessRule `json:"rules"`
}

type NetworkAccessRule struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	Enabled          bool     `json:"enabled,omitempty"`
	DBAccounts       []string `json:"dbAccounts,omitempty"`
	SourceIPs        []string `json:"sourceIPs,omitempty"`
	RulesBlockAccess bool     `json:"rulesBlockAccess,omitempty"`
}

func (nap *NetworkAccessPolicy) ReadFromSchema(d *schema.ResourceData) error {
	var networkAccessRulesIfaces []interface{}
	if set, ok := d.GetOk("network_access_rule"); ok {
		log.Printf("[DEBUG] If1")
		networkAccessRulesIfaces = set.(*schema.Set).List()
	} else {
		log.Printf("[DEBUG] If2.")
		return nil
	}

	nap.NetworkAccessRules = NetworkAccessRules{Rules: []NetworkAccessRule{}}
	for _, networkAccessRuleIface := range networkAccessRulesIfaces {
		networkAccessRuleMap := networkAccessRuleIface.(map[string]interface{})
		nap.NetworkAccessRules.Rules = append(nap.NetworkAccessRules.Rules,
			NetworkAccessRule{
				Name:        networkAccessRuleMap["name"].(string),
				Description: networkAccessRuleMap["description"].(string),
				DBAccounts:  getStrList(networkAccessRuleMap, "db_accounts"),
				SourceIPs:   getStrList(networkAccessRuleMap, "source_ips"),
			})
	}

	log.Printf("[DEBUG] Rules: %#v", nap.NetworkAccessRules.Rules)

	return nil
}

func (nap *NetworkAccessPolicy) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get("repository_id").(string))

	var networkAccessRules []interface{}
	for _, rule := range nap.NetworkAccessRules.Rules {
		rulesMap := map[string]interface{}{
			"name":        rule.Name,
			"description": rule.Description,
			"db_accounts": rule.DBAccounts,
			"source_ips":  rule.SourceIPs,
		}
		networkAccessRules = append(networkAccessRules, rulesMap)
	}

	return d.Set("network_access_rule", networkAccessRules)
}

func createRepositoryNetworkAccessPolicy() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "RepositoryNetworkAccessPolicyCreate",
		HttpMethod: http.MethodPost,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryNetworkAccessPolicyURLFormat,
				c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &NetworkAccessPolicy{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &NetworkAccessPolicyUpsertResp{} },
	}
}

func readRepositoryNetworkAccessPolicy() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "RepositoryNetworkAccessPolicyRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryNetworkAccessPolicyURLFormat,
				c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &NetworkAccessPolicy{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &NetworkAccessPolicy{} },
	}
}

func updateRepositoryNetworkAccessPolicy() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "RepositoryNetworkAccessPolicyUpdate",
		HttpMethod: http.MethodPut,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryNetworkAccessPolicyURLFormat,
				c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &NetworkAccessPolicy{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &NetworkAccessPolicyUpsertResp{} },
	}
}

func deleteRepositoryNetworkAccessPolicy() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "RepositoryNetworkAccessPolicyDelete",
		HttpMethod: http.MethodDelete,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryNetworkAccessPolicyURLFormat,
				c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &NetworkAccessPolicy{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &IDBasedResponse{} },
	}
}

func resourceRepositoryNetworkAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages the [Network Shield](https://cyral.com/docs/manage-repositories/network-shield/) of a repository.",
		CreateContext: CreateResource(createRepositoryNetworkAccessPolicy(), readRepositoryNetworkAccessPolicy()),
		ReadContext:   ReadResource(readRepositoryNetworkAccessPolicy()),
		UpdateContext: UpdateResource(updateRepositoryNetworkAccessPolicy(), readRepositoryNetworkAccessPolicy()),
		DeleteContext: DeleteResource(deleteRepositoryNetworkAccessPolicy()),

		Schema: map[string]*schema.Schema{
			"repository_id": {
				Description: "ID of the repository for which to configure a network access policy.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"network_access_rule": {
				Description: "",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the rule.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"description": {
							Description: "Description of the network access policy.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"db_accounts": {
							Description: "Specify which accounts this rule applies to. The account name must match an existing account in your database. See also [cyral_repository_local_account](./repository_local_account.md).",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"source_ips": {
							Description: "Specify IPs to restrict the range of allowed IP addresses for this rule.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"id": {
				Description: "ID of this resource in the Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set("repository_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
