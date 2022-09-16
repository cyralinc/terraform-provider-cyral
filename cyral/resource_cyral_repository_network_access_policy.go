package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NetworkAccessPolicy struct {
	NetworkAccessRules `json:"networkAccessRules"`
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
	return nil
}

func (nap *NetworkAccessPolicy) WriteToSchema(d *schema.ResourceData) error {

	return nil
}

func createRepositoryNetworkAccessPolicy() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "CreateRepositoryNetworkAccessPolicy",
		HttpMethod: http.MethodPost,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/repos/%s/networkAccessPolicy",
				c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &NetworkAccessPolicy{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &NetworkAccessPolicy{} },
	}
}

func readRepositoryNetworkAccessPolicy() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "ReadRepositoryNetworkAccessPolicy",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/repos/%s/networkAccessPolicy",
				c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &NetworkAccessPolicy{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &NetworkAccessPolicy{} },
	}
}

func updateRepositoryNetworkAccessPolicy() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "UpdateRepositoryNetworkAccessPolicy",
		HttpMethod: http.MethodPut,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/repos/%s/networkAccessPolicy",
				c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &NetworkAccessPolicy{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &NetworkAccessPolicy{} },
	}
}

func deleteRepositoryNetworkAccessPolicy() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "DeleteRepositoryNetworkAccessPolicy",
		HttpMethod: http.MethodDelete,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/repos/%s/networkAccessPolicy",
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
				// The API is open to the addition new types of
				// network access policies. If new network
				// access policies are added in the future, this
				// field should probably be moved to `Optional:
				// true`, and we should add a flag saying that
				// the two network policy type fields conflict
				// with each other.
				Required: true,
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
