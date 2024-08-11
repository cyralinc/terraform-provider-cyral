package useraccount

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var allAuthSchemes = []string{
	"aws_iam",
	"aws_secrets_manager",
	"cyral_storage",
	"hashicorp_vault",
	"environment_variable",
	"kubernetes_secret",
	"gcp_secrets_manager",
	"azure_key_vault",
}

var urlFactory = func(d *schema.ResourceData, c *client.Client) string {
	ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
	if err != nil {
		panic(fmt.Errorf("Unable to unmarshal composed id: %w", err))
	}
	repositoryID := ids[0]
	userAccountID := ids[1]
	return fmt.Sprintf(
		"https://%s/v1/repos/%s/userAccounts/%s",
		c.ControlPlane,
		repositoryID,
		userAccountID,
	)
}

var readRepositoryUserAccountConfig = core.ResourceOperationConfig{
	ResourceName: resourceName,
	Type:         operationtype.Read,
	HttpMethod:   http.MethodGet,
	URLFactory:   urlFactory,
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
		return &UserAccountResource{}
	},
	RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: "User account"},
}

func resourceSchema() *schema.Resource {
	authSchemeTypesFullScopes := make([]string, 0, len(allAuthSchemes))
	for _, authType := range allAuthSchemes {
		authSchemeTypesFullScopes = append(authSchemeTypesFullScopes,
			fmt.Sprintf("auth_scheme.0.%s", authType))
	}
	return &schema.Resource{
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				ResourceName: resourceName,
				Type:         operationtype.Create,
				HttpMethod:   http.MethodPost,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/userAccounts",
						c.ControlPlane,
						d.Get("repository_id").(string),
					)
				},
				SchemaReaderFactory: func() core.SchemaReader { return &UserAccountResource{} },
				SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &CreateUserAccountResponse{} },
			},
			readRepositoryUserAccountConfig,
		),
		ReadContext: core.ReadResource(readRepositoryUserAccountConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				ResourceName:        resourceName,
				Type:                operationtype.Update,
				HttpMethod:          http.MethodPut,
				URLFactory:          urlFactory,
				SchemaReaderFactory: func() core.SchemaReader { return &UserAccountResource{} },
			},
			readRepositoryUserAccountConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				ResourceName: resourceName,
				HttpMethod:   http.MethodDelete,
				URLFactory:   urlFactory,
			},
		),
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				i interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to unmarshal ID: %v",
						err,
					)
				}
				repositoryID := ids[0]
				err = d.Set("repository_id", repositoryID)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to set 'repository_id': %v",
						err,
					)
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Terraform ID of this resource. Follows syntax `{repository_id}/{user_account_id}`",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"repository_id": {
				Description: "ID of the repository.",
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
			},

			"user_account_id": {
				Description: "ID of the user account.",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"name": {
				Description: "The name of the User Account.",
				Required:    true,
				Type:        schema.TypeString,
			},

			"auth_database_name": {
				Description: "The database name that this User Account is scoped to, for " +
					"`cyral_repository` types that support multiple databases.",
				Optional: true,
				Type:     schema.TypeString,
			},

			"approval_config": {
				Description: "Configurations related to Approvals.",
				Optional:    true,
				Type:        schema.TypeSet,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"automatic_grant": {
							Description: "If `true`, approvals can be automatically granted.",
							Required:    true,
							Type:        schema.TypeBool,
						},
						"max_auto_grant_duration": {
							Description: "The maximum duration in seconds for approvals can be " +
								"automatically granted. E.g.: `\"2000s\"`, `\"3000.5s\"",
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},

			"auth_scheme": {
				Description: "Credential option. List of supported types: " +
					utils.SupportedValuesAsMarkdown(allAuthSchemes),
				Required: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_variable": {
							Description: "Credential option to set the repository user account from " +
								"Environment Variable.",
							Optional:     true,
							Type:         schema.TypeSet,
							ExactlyOneOf: authSchemeTypesFullScopes,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"variable_name": {
										Description: "Name of the environment variable that will store credentials.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},

						"aws_iam": {
							Description: "Credential option to set the repository user account from " +
								"AWS IAM.",
							Optional:     true,
							Type:         schema.TypeSet,
							ExactlyOneOf: authSchemeTypesFullScopes,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_arn": {
										Description: "The AWS IAM roleARN to gain access to the database.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"authenticate_as_iam_role": {
										Description: "Indicates whether to access as an AWS IAM role (`true`)" +
											"or a native database user (`false`). Defaults to `false`.",
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},

						"aws_secrets_manager": {
							Description: "Credential option to set the repository user account from " +
								"AWS Secrets Manager.",
							Optional:     true,
							Type:         schema.TypeSet,
							ExactlyOneOf: authSchemeTypesFullScopes,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_arn": {
										Description: "The AWS Secrets Manager secretARN to gain access to the database.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},

						"cyral_storage": {
							Description: "Credential option to set the repository user account from " +
								"Cyral Storage.",
							Optional:     true,
							Type:         schema.TypeSet,
							ExactlyOneOf: authSchemeTypesFullScopes,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"password": {
										Description: "The Cyral Storage password to gain access to the database.",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
								},
							},
						},

						"hashicorp_vault": {
							Description: "Credential option to set the repository user account from " +
								"Hashicorp Vault.",
							Optional:     true,
							Type:         schema.TypeSet,
							ExactlyOneOf: authSchemeTypesFullScopes,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Description: "The location in the Vault where the database username and" +
											" password may be retrieved.",
										Type:     schema.TypeString,
										Required: true,
									},
									"is_dynamic_user_account": {
										Description: "Some Vault engines allow the dynamic creation of user accounts," +
											" meaning the username used to log in to the database may change from time to time.",
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},

						"kubernetes_secret": {
							Description: "Credential option to set the repository user account from " +
								"a Kubernetes secret.",
							Optional:     true,
							Type:         schema.TypeSet,
							ExactlyOneOf: authSchemeTypesFullScopes,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_name": {
										Description: "The unique identifier of the secret in Kubernetes.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"secret_key": {
										Description: "The key of the credentials JSON blob within the secret.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},

						"gcp_secrets_manager": {
							Description: "Credential option to set the repository user account from " +
								"GCP Secrets Manager.",
							Optional:     true,
							Type:         schema.TypeSet,
							ExactlyOneOf: authSchemeTypesFullScopes,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_name": {
										Description: "The unique identifier of the secret in GCP Secrets Manager.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},

						"azure_key_vault": {
							Description: "Credential option to set the repository user account from " +
								"Azure Key Vault.",
							Optional:     true,
							Type:         schema.TypeSet,
							ExactlyOneOf: authSchemeTypesFullScopes,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_url": {
										Description: "The URL of the secret in the Azure Key Vault.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
