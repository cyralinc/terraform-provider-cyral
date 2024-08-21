package confauth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO This resource is more complex than it should be due to the fact that a call to
// repo creation automatically creates the conf/auth and also the conf/analysis configurations.
// Our API should be refactored so these operations should happen separately.

var urlFactory = func(d *schema.ResourceData, c *client.Client) string {
	return fmt.Sprintf("https://%s/v1/repos/%s/conf/auth",
		c.ControlPlane,
		d.Get("repository_id"),
	)
}

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:                 resourceName,
	ResourceType:                 resourcetype.Resource,
	SchemaReaderFactory:          func() core.SchemaReader { return &RepositoryConfAuthData{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &ReadRepositoryConfAuthResponse{} },
	BaseURLFactory:               urlFactory,
	ReadUpdateDeleteURLFactory:   urlFactory,
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages Repository Analysis Configuration. This resource allows configuring both " +
			"[Log Settings](https://cyral.com/docs/manage-repositories/repo-log-volume) " +
			"and [Advanced settings](https://cyral.com/docs/manage-repositories/repo-advanced-settings) " +
			"(Logs, Alerts, Analysis and Enforcement) configurations for Data Repositories.",
		CreateContext: resourceRepositoryConfAuthCreate,
		ReadContext: resourceContextHandler.ReadContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "Failed to read repo",
			OperationType:  operationtype.Read,
		}),
		UpdateContext: resourceContextHandler.UpdateContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "Failed to read repo",
			OperationType:  operationtype.Update,
		}, nil),
		DeleteContext: resourceContextHandler.DeleteContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "Failed to read repo",
			OperationType:  operationtype.Delete,
		}),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: repositoryConfAuthResourceSchemaV0().
					CoreConfigSchema().ImpliedType(),
				Upgrade: UpgradeRepositoryConfAuthV0,
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
				Description: fmt.Sprintf(
					"Specifies whether the sidecar will require TLS communication with clients."+
						" Defaults to `%q`. List of supported values: %q", TLSDisable, utils.SupportedValuesAsMarkdown(ClientTLSTypesAsString())),
				Type:         schema.TypeString,
				Optional:     true,
				Default:      TLSDisable,
				ValidateFunc: validation.StringInSlice(append(ClientTLSTypesAsString(), ""), false),
			},
			"identity_provider": {
				Description: fmt.Sprintf(
					"The semantics of this field changed in control planes `v4.13` and later. See how "+
						"it should be configured depending on your control plane version:\n"+
						"\t- `v4.12` and below:\n\t\t- Provide the ID (Alias) of the identity provider "+
						"integration to allow user authentication using an IdP.\n"+
						"\t- `v4.13` and later:\n\t\t- If not supplied, then end-user "+
						"authentication is disabled.\n\t\t- If end-user authentication "+
						"with Cyral Access Token is desired, then set to `ACCESS_TOKEN` or any "+
						"other non-empty string.\n\t\t- If end-user authentication with "+
						"AWS IAM is desired, then this must be the ID of an AWS IAM integration, "+
						"and the `auth_type` attribute must be set to `%s`.",
					AwsIAMAuthType,
				),
				Type:     schema.TypeString,
				Optional: true,
			},
			"repo_tls": {
				Description: fmt.Sprintf(
					"Specifies whether the sidecar will communicate with the repository using TLS."+
						" Defaults to `%q`. List of supported values: %q", TLSDisable, utils.SupportedValuesAsMarkdown(RepoTLSTypesAsString())),
				Type:         schema.TypeString,
				Optional:     true,
				Default:      TLSDisable,
				ValidateFunc: validation.StringInSlice(append(RepoTLSTypesAsString(), ""), false),
			},
			"auth_type": {
				Description: fmt.Sprintf("Authentication type for this repository. **Note**: `%s` is currently "+
					"only supported by `%s` repo type. List of supported values: %s",
					AwsIAMAuthType, repository.MongoDB, utils.SupportedValuesAsMarkdown(authTypes)),
				Type:         schema.TypeString,
				Optional:     true,
				Default:      DefaultAuthType,
				ValidateFunc: validation.StringInSlice(authTypes, false),
			},
		},
	}
}

// Previously, the id of the resource `cyral_repository_conf_auth` was hardcoded
// to `repo-conf`, which doesn't make sense. The goal here is to set it to be
// the repository ID.
func UpgradeRepositoryConfAuthV0(
	_ context.Context,
	rawState map[string]interface{},
	_ interface{},
) (map[string]interface{}, error) {
	rawState["id"] = rawState["repository_id"]
	return rawState, nil
}

func resourceRepositoryConfAuthCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRepositoryConfAuthCreate")
	c := m.(*client.Client)
	httpMethod := http.MethodPost
	if confAuthAlreadyExists(ctx, c, d) {
		httpMethod = http.MethodPut
	}
	tflog.Debug(ctx, "End resourceRepositoryConfAuthCreate")
	return core.CreateResource(
		core.ResourceOperationConfig{
			ResourceName:        resourceName,
			Type:                operationtype.Create,
			HttpMethod:          httpMethod,
			URLFactory:          urlFactory,
			SchemaReaderFactory: func() core.SchemaReader { return &RepositoryConfAuthData{} },
			SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &CreateRepositoryConfAuthResponse{} },
		},
		core.ResourceOperationConfig{
			ResourceName:        resourceName,
			ResourceType:        resourcetype.Resource,
			Type:                operationtype.Read,
			HttpMethod:          http.MethodGet,
			URLFactory:          urlFactory,
			SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &ReadRepositoryConfAuthResponse{} },
			RequestErrorHandler: &core.IgnoreNotFoundByMessage{
				ResName:        resourceName,
				MessageMatches: "Failed to read repo",
				OperationType:  operationtype.Read,
			},
		},
	)(ctx, d, m)
}

func confAuthAlreadyExists(ctx context.Context, c *client.Client, d *schema.ResourceData) bool {
	_, err := c.DoRequest(ctx, urlFactory(d, c), http.MethodGet, nil)
	// See TODO on the top of this file
	if err != nil {

		tflog.Debug(ctx, fmt.Sprintf("Unable to read Conf Auth resource for repository %s: %v",
			d.Get("repository_id").(string), err))
		return false
	}
	return true
}
