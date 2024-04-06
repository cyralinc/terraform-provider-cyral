package network

import (
	"context"
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var urlFactory = func(d *schema.ResourceData, c *client.Client) string {
	return fmt.Sprintf("https://%s/v1/repos/%s/networkAccessPolicy",
		c.ControlPlane,
		d.Get("repository_id"),
	)
}

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:                  resourceName,
	ResourceType:                  resourcetype.Resource,
	SchemaReaderFactory:           func() core.SchemaReader { return &NetworkAccessPolicy{} },
	SchemaWriterFactoryGetMethod:  func(_ *schema.ResourceData) core.SchemaWriter { return &NetworkAccessPolicy{} },
	SchemaWriterFactoryPostMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &NetworkAccessPolicyUpsertResp{} },
	BaseURLFactory:                urlFactory,
	ReadUpdateDeleteURLFactory:    urlFactory,
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the network access policy of a repository. Network access policies are" +
			" also known as the [Network Shield](https://cyral.com/docs/manage-repositories/network-shield/)." +
			" This feature is supported for the following repository types:" +
			utils.SupportedValuesAsMarkdown(repositoryTypesNetworkShield()) +
			"\n\n-> **Note** If you also use the resource `cyral_repository_conf_auth` for the same repository," +
			" create a `depends_on` relationship from this resource to the `cyral_repository_conf_auth` to" +
			" avoid errors when running `terraform destroy`.",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),

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
