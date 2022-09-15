package cyral

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func repositoryTypesNetworkShield() []string {
	return []string{
		"sqlserver",
		"oracle",
	}
}

type RepositoryConfAuthData struct {
	RepoID                     *string `json:"-"`
	AllowNativeAuth            bool    `json:"allowNativeAuth"`
	ClientTLS                  string  `json:"clientTLS"`
	IdentityProvider           string  `json:"identityProvider"`
	RepoTLS                    string  `json:"repoTLS"`
	EnableNetworkAccessControl bool    `json:"enableNetworkAccessControl"`
}

func (data RepositoryConfAuthData) WriteToSchema(d *schema.ResourceData) error {
	if data.RepoID != nil {
		d.Set("repository_id", data.RepoID)
	}

	d.Set("allow_native_auth", data.AllowNativeAuth)

	if err := data.isClientTLSValid(); err != nil {
		panic(err)
	}

	d.Set("client_tls", data.ClientTLS)

	d.Set("identity_provider", data.IdentityProvider)

	if err := data.isRepoTLSValid(); err != nil {
		panic(err)
	}

	d.Set("repo_tls", data.RepoTLS)

	return nil
}

func (data *RepositoryConfAuthData) ReadFromSchema(d *schema.ResourceData) error {
	if repoIdData, hasRepoId := d.GetOk("repository_id"); hasRepoId {
		repoId := repoIdData.(string)
		data.RepoID = &repoId
	}

	data.AllowNativeAuth = d.Get("allow_native_auth").(bool)
	data.ClientTLS = d.Get("client_tls").(string)
	data.IdentityProvider = d.Get("identity_provider").(string)
	data.RepoTLS = d.Get("repo_tls").(string)

	return nil
}

func (data RepositoryConfAuthData) isClientTLSValid() error {
	if !(data.ClientTLS == "enable" || data.ClientTLS == "disable" || data.ClientTLS == "enabledAndVerifyCertificate") {
		return errors.New("invalid option to client_tls")
	}
	return nil
}

func (data RepositoryConfAuthData) isRepoTLSValid() error {
	if !(data.RepoTLS == "enable" || data.RepoTLS == "disable" || data.RepoTLS == "enabledAndVerifyCertificate") {
		return errors.New("invalid option to repo_tls")
	}
	return nil
}

type CreateRepositoryConfAuthResponse struct{}

func (data CreateRepositoryConfAuthResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId("repo-conf")
	return nil
}

type ReadRepositoryConfAuthResponse struct {
	AuthInfo RepositoryConfAuthData `json:"authInfo"`
}

func (data ReadRepositoryConfAuthResponse) WriteToSchema(d *schema.ResourceData) error {
	data.AuthInfo.WriteToSchema(d)
	return nil
}

func CreateConfAuthConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "ConfAuthResourceCreate",
		HttpMethod: http.MethodPut,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/repos/%s/conf/auth", c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &RepositoryConfAuthData{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &CreateRepositoryConfAuthResponse{} },
	}
}

func ReadConfAuthConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "ConfAuthResourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/repos/%s/conf/auth", c.ControlPlane, d.Get("repository_id"))
		},
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &ReadRepositoryConfAuthResponse{} },
	}
}

func UpdateConfAuthConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "ConfAuthResourceUpdate",
		HttpMethod: http.MethodPut,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/repos/%s/conf/auth", c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &RepositoryConfAuthData{} },
	}
}

func DeleteConfAuthConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "ConfAuthResourceDelete",
		HttpMethod: http.MethodDelete,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/repos/%s/conf/auth", c.ControlPlane, d.Get("repository_id"))
		},
	}
}

func resourceRepositoryConfAuth() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages the [Repository Authentication settings](https://cyral.com/docs/manage-repositories/repo-advanced-settings/#authentication) that is shown in the Advanced tab.",
		CreateContext: CreateResource(CreateConfAuthConfig(), ReadConfAuthConfig()),
		ReadContext:   ReadResource(ReadConfAuthConfig()),
		UpdateContext: UpdateResource(UpdateConfAuthConfig(), ReadConfAuthConfig()),
		DeleteContext: DeleteResource(DeleteConfAuthConfig()),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"repository_id": {
				Description: "The ID of the repository to be configured.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"allow_native_auth": {
				Description: "Should the communication allow native authentication?",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"client_tls": {
				Description: "Is the repo Client using TLS?",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"identity_provider": {
				Description: "The ID (Alias) of the identity provider integration.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"repo_tls": {
				Description: "Is TLS enabled for the repository?",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enable_network_access_control": {
				Description: "If set to true, enables the [Network Shield](https://cyral.com/docs/manage-repositories/network-shield/) feature for the repository. This feature is supported for the following repository types:" + supportedTypesMarkdown(repositoryTypesNetworkShield()),
				Type:        schema.TypeString,
				Optional:    true,
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
