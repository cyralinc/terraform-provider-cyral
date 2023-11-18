package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	repositoryNetworkAccessPolicyURLFormat = "https://%s/v1/repos/%s/networkAccessPolicy"

	defaultNetworkAccessPolicyEnabled    = true
	defaultNetworkAccessRulesBlockAccess = false
)

func repositoryTypesNetworkShield() []string {
	return []string{
		repository.SQLServer,
		repository.Oracle,
	}
}

type NetworkAccessPolicyUpsertResp struct {
	Policy NetworkAccessPolicy `json:"policy"`
}

func (resp NetworkAccessPolicyUpsertResp) WriteToSchema(d *schema.ResourceData) error {
	return resp.Policy.WriteToSchema(d)
}

type NetworkAccessPolicy struct {
	Enabled            bool `json:"enabled"`
	NetworkAccessRules `json:"networkAccessRules,omitempty"`
}

type NetworkAccessRules struct {
	RulesBlockAccess bool                `json:"rulesBlockAccess"`
	Rules            []NetworkAccessRule `json:"rules"`
}

type NetworkAccessRule struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	DBAccounts  []string `json:"dbAccounts,omitempty"`
	SourceIPs   []string `json:"sourceIPs,omitempty"`
}

func (nap *NetworkAccessPolicy) ReadFromSchema(d *schema.ResourceData) error {
	nap.Enabled = d.Get("enabled").(bool)

	var networkAccessRulesIfaces []interface{}
	if set, ok := d.GetOk("network_access_rule"); ok {
		networkAccessRulesIfaces = set.(*schema.Set).List()
	} else {
		return nil
	}

	nap.NetworkAccessRules = NetworkAccessRules{
		RulesBlockAccess: d.Get("network_access_rules_block_access").(bool),
		Rules:            []NetworkAccessRule{},
	}
	for _, networkAccessRuleIface := range networkAccessRulesIfaces {
		networkAccessRuleMap := networkAccessRuleIface.(map[string]interface{})
		nap.NetworkAccessRules.Rules = append(nap.NetworkAccessRules.Rules,
			NetworkAccessRule{
				Name:        networkAccessRuleMap["name"].(string),
				Description: networkAccessRuleMap["description"].(string),
				DBAccounts:  utils.GetStrList(networkAccessRuleMap, "db_accounts"),
				SourceIPs:   utils.GetStrList(networkAccessRuleMap, "source_ips"),
			})
	}

	return nil
}

func (nap *NetworkAccessPolicy) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get("repository_id").(string))
	d.Set("enabled", nap.Enabled)
	d.Set("network_access_rules_block_access", nap.NetworkAccessRules.RulesBlockAccess)

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

func createRepositoryNetworkAccessPolicy() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "RepositoryNetworkAccessPolicyCreate",
		HttpMethod:   http.MethodPost,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryNetworkAccessPolicyURLFormat,
				c.ControlPlane, d.Get("repository_id"))
		},
		SchemaReaderFactory: func() core.SchemaReader { return &NetworkAccessPolicy{} },
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &NetworkAccessPolicyUpsertResp{} },
	}
}

func readRepositoryNetworkAccessPolicy() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "RepositoryNetworkAccessPolicyRead",
		HttpMethod:   http.MethodGet,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryNetworkAccessPolicyURLFormat,
				c.ControlPlane, d.Get("repository_id"))
		},
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &NetworkAccessPolicy{} },
		RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Repository network access policy"},
	}
}

func updateRepositoryNetworkAccessPolicy() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "RepositoryNetworkAccessPolicyUpdate",
		HttpMethod:   http.MethodPut,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryNetworkAccessPolicyURLFormat,
				c.ControlPlane, d.Get("repository_id"))
		},
		SchemaReaderFactory: func() core.SchemaReader { return &NetworkAccessPolicy{} },
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &NetworkAccessPolicyUpsertResp{} },
	}
}

func deleteRepositoryNetworkAccessPolicy() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "RepositoryNetworkAccessPolicyDelete",
		HttpMethod:   http.MethodDelete,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryNetworkAccessPolicyURLFormat,
				c.ControlPlane, d.Get("repository_id"))
		},
		RequestErrorHandler: &core.DeleteIgnoreHttpNotFound{ResName: "Network Access Policy"},
	}
}

func ResourceRepositoryNetworkAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the network access policy of a repository. Network access policies are" +
			" also known as the [Network Shield](https://cyral.com/docs/manage-repositories/network-shield/)." +
			" This feature is supported for the following repository types:" +
			utils.SupportedValuesAsMarkdown(repositoryTypesNetworkShield()) +
			"\n\n-> **Note** If you also use the resource `cyral_repository_conf_auth` for the same repository," +
			" create a `depends_on` relationship from this resource to the `cyral_repository_conf_auth` to" +
			" avoid errors when running `terraform destroy`.",
		CreateContext: core.CreateResource(createRepositoryNetworkAccessPolicy(), readRepositoryNetworkAccessPolicy()),
		ReadContext:   core.ReadResource(readRepositoryNetworkAccessPolicy()),
		UpdateContext: core.UpdateResource(updateRepositoryNetworkAccessPolicy(), readRepositoryNetworkAccessPolicy()),
		DeleteContext: core.DeleteResource(deleteRepositoryNetworkAccessPolicy()),

		Schema: map[string]*schema.Schema{
			"repository_id": {
				Description: "ID of the repository for which to configure a network access policy.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			// This parameter also exists in the
			// v1/repos/{repoID}/conf/auth API, but putting it under
			// the `cyral_repository_conf_auth` resource was causing
			// a lot of trouble: the resources would get out of sync
			// and behave like crazy.
			"enabled": {
				Description: fmt.Sprintf("Is the network access policy enabled? Default is %t.", defaultNetworkAccessPolicyEnabled),
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     defaultNetworkAccessPolicyEnabled,
			},

			"network_access_rules_block_access": {
				Description: fmt.Sprintf("Determines what happens if an incoming connection matches one of the rules in `network_access_rule`. If set to true, the connection is blocked if it matches some rule (and allowed otherwise). Otherwise set to false, the connection is allowed only if it matches some rule. Default is %t.", defaultNetworkAccessRulesBlockAccess),
				Type:        schema.TypeBool,
				Optional:    true,
			},

			"network_access_rule": {
				Description: "Network access policy that decides whether access should be granted based on a set of rules.",
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
							Description: "Specify which accounts this rule applies to. The account name must match an existing account in your database.",
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
