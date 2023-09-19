package cyral

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	repositoryConfAuthURLFormat = "https://%s/v1/repos/%s/conf/auth"

	defaultClientTLS = "disable"
	defaultRepoTLS   = "disable"

	accessTokenAuthType = "ACCESS_TOKEN"
	awsIAMAuthType      = "AWS_IAM"
	defaultAuthType     = accessTokenAuthType
)

var authTypes = []string{
	accessTokenAuthType,
	awsIAMAuthType,
}

type RepositoryConfAuthData struct {
	RepoID           *string `json:"-"`
	AllowNativeAuth  bool    `json:"allowNativeAuth"`
	ClientTLS        string  `json:"clientTLS"`
	IdentityProvider string  `json:"identityProvider"`
	RepoTLS          string  `json:"repoTLS"`
	AuthType         string  `json:"authType"`
}

func (data RepositoryConfAuthData) WriteToSchema(d *schema.ResourceData, c *client.Client) error {
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

	d.Set("auth_type", data.AuthType)

	return nil
}

func (data *RepositoryConfAuthData) ReadFromSchema(d *schema.ResourceData, c *client.Client) error {
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

func (data CreateRepositoryConfAuthResponse) WriteToSchema(d *schema.ResourceData, c *client.Client) error {
	d.SetId(d.Get("repository_id").(string))
	return nil
}

type ReadRepositoryConfAuthResponse struct {
	AuthInfo RepositoryConfAuthData `json:"authInfo"`
}

func (data ReadRepositoryConfAuthResponse) WriteToSchema(d *schema.ResourceData, c *client.Client) error {
	data.AuthInfo.WriteToSchema(d, c)
	return nil
}

func resourceRepositoryConfAuthCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryConfAuthCreate")
	c := m.(*client.Client)
	httpMethod := http.MethodPost
	if confAuthAlreadyExists(c, d.Get("repository_id").(string)) {
		httpMethod = http.MethodPut
	}
	defer log.Printf("[DEBUG] End resourceRepositoryConfAuthCreate")
	return CreateResource(CreateConfAuthConfig(httpMethod), ReadConfAuthConfig())(ctx, d, m)
}

func confAuthAlreadyExists(c *client.Client, repositoryID string) bool {
	url := fmt.Sprintf(repositoryConfAuthURLFormat, c.ControlPlane, repositoryID)
	_, err := c.DoRequest(url, http.MethodGet, nil)
	// The GET /v1/repos/{repoID}/conf/auth API currently returns 500 status code for every type
	// of error, so its not possible to distinguish if the error is due to a 404 Not Found or not.
	// Once the status code returned by this API is fixed we should return false only if it returns
	// a 404 Not Found, otherwise, if a different error occurs, this function should return an error.
	if err != nil {
		log.Printf("[DEBUG] Unable to read Conf Auth resource for repository %s: %v", repositoryID, err)
		return false
	}
	return true
}

func CreateConfAuthConfig(httpMethod string) ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "ConfAuthResourceCreate",
		HttpMethod: httpMethod,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryConfAuthURLFormat, c.ControlPlane, d.Get("repository_id"))
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
			return fmt.Sprintf(repositoryConfAuthURLFormat, c.ControlPlane, d.Get("repository_id"))
		},
		NewResponseData:     func(_ *schema.ResourceData) ResponseData { return &ReadRepositoryConfAuthResponse{} },
		RequestErrorHandler: &ReadIgnoreHttpNotFound{resName: "Repository conf auth"},
	}
}

func UpdateConfAuthConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "ConfAuthResourceUpdate",
		HttpMethod: http.MethodPut,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryConfAuthURLFormat, c.ControlPlane, d.Get("repository_id"))
		},
		NewResourceData: func() ResourceData { return &RepositoryConfAuthData{} },
	}
}

func DeleteConfAuthConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "ConfAuthResourceDelete",
		HttpMethod: http.MethodDelete,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(repositoryConfAuthURLFormat, c.ControlPlane, d.Get("repository_id"))
		},
	}
}

func repositoryConfAuthResourceSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of this resource is set to `repository_id`.",
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
				Description: fmt.Sprintf("Is the repo Client using TLS? Default is %q.", defaultClientTLS),
				Type:        schema.TypeString,
				Optional:    true,
				Default:     defaultClientTLS,
			},
			"identity_provider": {
				Description: "The ID (Alias) of the identity provider integration.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"repo_tls": {
				Description: fmt.Sprintf("Is TLS enabled for the repository? Default is %q.", defaultRepoTLS),
				Type:        schema.TypeString,
				Optional:    true,
				Default:     defaultRepoTLS,
			},
			"auth_type": {
				Description: fmt.Sprintf("Authentication type for this repository. **Note**: `%s` is currently "+
					"only supported by `%s` repo type. List of supported values: %s",
					awsIAMAuthType, MongoDB, supportedTypesMarkdown(authTypes)),
				Type:         schema.TypeString,
				Optional:     true,
				Default:      defaultAuthType,
				ValidateFunc: validation.StringInSlice(authTypes, false),
			},
		},
	}
}

// Previously, the id of the resource `cyral_repository_conf_auth` was hardcoded
// to `repo-conf`, which doesn't make sense. The goal here is to set it to be
// the repository ID.
func upgradeRepositoryConfAuthV0(
	_ context.Context,
	rawState map[string]interface{},
	_ interface{},
) (map[string]interface{}, error) {
	rawState["id"] = rawState["repository_id"]
	return rawState, nil
}

func resourceRepositoryConfAuth() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages the [Repository Authentication settings](https://cyral.com/docs/manage-repositories/repo-advanced-settings/#authentication) that is shown in the Advanced tab.",
		CreateContext: resourceRepositoryConfAuthCreate,
		ReadContext:   ReadResource(ReadConfAuthConfig()),
		UpdateContext: UpdateResource(UpdateConfAuthConfig(), ReadConfAuthConfig()),
		DeleteContext: DeleteResource(DeleteConfAuthConfig()),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: repositoryConfAuthResourceSchemaV0().
					CoreConfigSchema().ImpliedType(),
				Upgrade: upgradeRepositoryConfAuthV0,
			},
		},

		Schema: repositoryConfAuthResourceSchemaV0().Schema,

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
