package cyral

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RepositoryConfAuthData struct {
	RepoID           *string `json:"-"`
	AllowNativeAuth  bool    `json:"allowNativeAuth"`
	ClientTLS        string  `json:"clientTLS"`
	IdentityProvider string  `json:"identity_provider"`
	RepoTLS          string  `json:"repoTLS"`
}

func (data RepositoryConfAuthData) WriteToSchema(d *schema.ResourceData) {
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
}

func (data *RepositoryConfAuthData) ReadFromSchema(d *schema.ResourceData) {
	if repoIdData, hasRepoId := d.GetOk("repository_id"); hasRepoId {
		repoId := repoIdData.(string)
		data.RepoID = &repoId
	}

	data.AllowNativeAuth = d.Get("allow_native_auth").(bool)
	data.ClientTLS = d.Get("client_tls").(string)
	data.IdentityProvider = d.Get("identity_provider").(string)
	data.RepoTLS = d.Get("repo_tls").(string)

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

func (data CreateRepositoryConfAuthResponse) WriteToSchema(d *schema.ResourceData) {
	d.SetId("repo-conf")
}

func (data *CreateRepositoryConfAuthResponse) ReadFromSchema(d *schema.ResourceData) {}
