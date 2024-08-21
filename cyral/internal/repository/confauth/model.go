package confauth

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var authTypes = []string{
	AccessTokenAuthType,
	AwsIAMAuthType,
}

type RepositoryConfAuthData struct {
	RepoID           *string `json:"-"`
	AllowNativeAuth  bool    `json:"allowNativeAuth"`
	ClientTLS        string  `json:"clientTLS"`
	IdentityProvider string  `json:"identityProvider"`
	RepoTLS          string  `json:"repoTLS"`
	AuthType         string  `json:"authType"`
}

func (data RepositoryConfAuthData) WriteToSchema(d *schema.ResourceData) error {
	if data.RepoID != nil {
		d.Set("repository_id", data.RepoID)
	}

	d.Set("allow_native_auth", data.AllowNativeAuth)
	d.Set("client_tls", data.ClientTLS)
	d.Set("identity_provider", data.IdentityProvider)
	d.Set("repo_tls", data.RepoTLS)
	d.Set("auth_type", data.AuthType)

	return nil
}

func (data *RepositoryConfAuthData) ReadFromSchema(d *schema.ResourceData) error {
	if repoIdData, hasRepoId := d.GetOk("repository_id"); hasRepoId {
		repoId := repoIdData.(string)
		data.RepoID = &repoId
	}

	data.AllowNativeAuth = d.Get("allow_native_auth").(bool)
	data.AuthType = d.Get("auth_type").(string)
	data.ClientTLS = d.Get("client_tls").(string)
	data.IdentityProvider = d.Get("identity_provider").(string)
	data.RepoTLS = d.Get("repo_tls").(string)

	return nil
}

type CreateRepositoryConfAuthResponse struct{}

func (data CreateRepositoryConfAuthResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get("repository_id").(string))
	return nil
}

type ReadRepositoryConfAuthResponse struct {
	AuthInfo RepositoryConfAuthData `json:"authInfo"`
}

func (data ReadRepositoryConfAuthResponse) WriteToSchema(d *schema.ResourceData) error {
	data.AuthInfo.WriteToSchema(d)
	return nil
}
