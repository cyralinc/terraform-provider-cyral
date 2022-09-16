package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"

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
}

type AuthScheme struct {
	AWSIAM              *AuthSchemeAWSIAM              `json:"awsIAM"`
	AWSSecretsManager   *AuthSchemeAWSSecretsManager   `json:"awsSecretsManager"`
	CyralStorage        *AuthSchemeCyralStorage        `json:"cyralStorage"`
	HashicorpVault      *AuthSchemeHashicorpVault      `json:"hashicorpVault"`
	EnvironmentVariable *AuthSchemeEnvironmentVariable `json:"environmentVariable"`
	KubernetesSecret    *AuthSchemeKubernetesSecret    `json:"kubernetesSecret"`
	GCPSecretManager    *AuthSchemeGCPSecretManager    `json:"gcpSecretManager"`
}

type AuthSchemeAWSIAM struct {
	RoleARN string `json:"roleARN,omitempty"`
}

type AuthSchemeAWSSecretsManager struct {
	SecretArn string `json:"secretARN,omitempty"`
}

type AuthSchemeCyralStorage struct {
	Password string `json:"password,omitempty"`
}

type AuthSchemeHashicorpVault struct {
	Path                 string `json:"path,omitempty"`
	IsDynamicUserAccount bool   `json:"isDynamicUserAccount,omitempty"`
}

type AuthSchemeEnvironmentVariable struct {
	VariableName string `json:"variableName,omitempty"`
}

type AuthSchemeKubernetesSecret struct {
	SecretName string `json:"secretName,omitempty"`
	SecretKey  string `json:"secretKey,omitempty"`
}

type AuthSchemeGCPSecretManager struct {
	SecretName string `json:"secretName,omitempty"`
}

type ApprovalConfig struct {
	AutomaticGrant            bool   `json:"automaticGrant,omitempty"`
	MaxAutomaticGrantDuration string `json:"maxAutomaticGrantDuration,omitempty"`
}

type UserAccountConfig struct {
	Approval *ApprovalConfig `json:"approvalConfig,omitempty"`
}

type UserAccountResource struct {
	UserAccountID    string             `json:"userAccountID,omitempty"`
	Name             string             `json:"name,omitempty"`
	AuthDatabaseName string             `json:"authDatabaseName,omitempty"`
	AuthScheme       *AuthScheme        `json:"authScheme,omitempty"`
	Config           *UserAccountConfig `json:"config,omitempty"`
}

type CreateUserAccountResponse struct {
	UserAccountID string `json:"userAccountID,omitempty"`
}

func (resp *CreateUserAccountResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(resp.UserAccountID)
	return nil
}

func (resource *UserAccountResource) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(resource.UserAccountID)
	if err := d.Set("name", resource.Name); err != nil {
		return fmt.Errorf("error setting 'name': %w", err)
	}
	if resource.AuthDatabaseName != "" {
		if err := d.Set("auth_database_name", resource.AuthDatabaseName); err != nil {
			return fmt.Errorf("error setting 'auth_database_name': %w", err)
		}
	}

	if resource.Config != nil {
		if resource.Config.Approval != nil {
			err := d.Set("approval_config", []interface{}{
				map[string]interface{}{
					"automatic_grant":         resource.Config.Approval.AutomaticGrant,
					"max_auto_grant_duration": resource.Config.Approval.MaxAutomaticGrantDuration,
				},
			})
			if err != nil {
				return fmt.Errorf("error setting 'approval_config': %w", err)
			}
		}
	}

	var authScheme []interface{}
	if resource.AuthScheme.AWSIAM != nil {
		authScheme = []interface{}{
			map[string]interface{}{
				"aws_iam": []interface{}{
					map[string]interface{}{
						"role_arn": resource.AuthScheme.AWSIAM.RoleARN,
					},
				},
			},
		}
	} else if resource.AuthScheme.AWSSecretsManager != nil {
		authScheme = []interface{}{
			map[string]interface{}{
				"aws_secrets_manager": []interface{}{
					map[string]interface{}{
						"secret_arn": resource.AuthScheme.AWSSecretsManager.SecretArn,
					},
				},
			},
		}
	} else if resource.AuthScheme.CyralStorage != nil {
		authScheme = []interface{}{
			map[string]interface{}{
				"cyral_storage": []interface{}{
					map[string]interface{}{
						"password": resource.AuthScheme.CyralStorage.Password,
					},
				},
			},
		}
	} else if resource.AuthScheme.EnvironmentVariable != nil {
		authScheme = []interface{}{
			map[string]interface{}{
				"environment_variable": []interface{}{
					map[string]interface{}{
						"variable_name": resource.AuthScheme.EnvironmentVariable.VariableName,
					},
				},
			},
		}
	} else if resource.AuthScheme.GCPSecretManager != nil {
		authScheme = []interface{}{
			map[string]interface{}{
				"gcp_secrets_manager": []interface{}{
					map[string]interface{}{
						"secret_name": resource.AuthScheme.GCPSecretManager.SecretName,
					},
				},
			},
		}
	} else if resource.AuthScheme.HashicorpVault != nil {
		authScheme = []interface{}{
			map[string]interface{}{
				"hashicorp_vault": []interface{}{
					map[string]interface{}{
						"path":                    resource.AuthScheme.HashicorpVault.Path,
						"is_dynamic_user_account": resource.AuthScheme.HashicorpVault.IsDynamicUserAccount,
					},
				},
			},
		}
	} else if resource.AuthScheme.KubernetesSecret != nil {
		authScheme = []interface{}{
			map[string]interface{}{
				"kubernetes_secret": []interface{}{
					map[string]interface{}{
						"secret_name": resource.AuthScheme.KubernetesSecret.SecretName,
						"secret_key":  resource.AuthScheme.KubernetesSecret.SecretKey,
					},
				},
			},
		}
	} else {
		return fmt.Errorf("auth scheme is required, user account is corrupt: %v", resource)
	}

	if err := d.Set("auth_scheme", authScheme); err != nil {
		return fmt.Errorf("error setting 'auth_scheme': %w", err)
	}

	return nil
}

// ReadFromSchema is used to translate a .tf file into whatever the
// UserAccounts API expects.
func (userAccount *UserAccountResource) ReadFromSchema(d *schema.ResourceData) error {
	// Set basic values required fields (UserID is computed).
	userAccount.Name = d.Get("name").(string)
	userAccount.AuthDatabaseName = d.Get("auth_database_name").(string)

	// Handle approval config (optional field).
	approvalConfig := d.Get("approval_config").(*schema.Set).List()
	if len(approvalConfig) > 0 {
		m := approvalConfig[0].(map[string]interface{})
		userAccount.Config = &UserAccountConfig{
			Approval: &ApprovalConfig{},
		}
		autogrant, ok := m["automatic_grant"]
		if ok {
			userAccount.Config.Approval.AutomaticGrant =
				autogrant.(bool)
		}
		maxAutoGrant, ok := m["max_auto_grant_duration"]
		if ok {
			userAccount.Config.Approval.MaxAutomaticGrantDuration =
				maxAutoGrant.(string)
		}
	}

	// Handle Auth Scheme (required field).
	authSchemeSet := d.Get("auth_scheme").(*schema.Set).List()
	if len(authSchemeSet) != 1 {
		return fmt.Errorf(
			"exactly one auth_scheme attribute is required",
		)
	}

	authSchemes := authSchemeSet[0].(map[string]interface{})

	for k, v := range authSchemes {
		authSchemeDetails := v.(*schema.Set).List()
		if len(authSchemeDetails) == 0 {
			continue
		}
		if len(authSchemeDetails) > 1 {
			return fmt.Errorf(
				"exactly one of %v is required in auth_scheme",
				allAuthSchemes,
			)
		}
		m := authSchemeDetails[0].(map[string]interface{})

		switch k {
		case "environment_variable":
			userAccount.AuthScheme = &AuthScheme{
				EnvironmentVariable: &AuthSchemeEnvironmentVariable{
					VariableName: m["variable_name"].(string),
				},
			}
		case "aws_iam":
			userAccount.AuthScheme = &AuthScheme{
				AWSIAM: &AuthSchemeAWSIAM{
					RoleARN: m["role_arn"].(string),
				},
			}
		case "aws_secrets_manager":
			userAccount.AuthScheme = &AuthScheme{
				AWSSecretsManager: &AuthSchemeAWSSecretsManager{
					SecretArn: m["secret_arn"].(string),
				},
			}
		case "cyral_storage":
			userAccount.AuthScheme = &AuthScheme{
				CyralStorage: &AuthSchemeCyralStorage{
					Password: m["password"].(string),
				},
			}
		case "hashicorp_vault":
			userAccount.AuthScheme = &AuthScheme{
				HashicorpVault: &AuthSchemeHashicorpVault{
					Path: m["path"].(string),
				},
			}
			if m["is_dynamic_user_account"] != nil {
				userAccount.AuthScheme.HashicorpVault.IsDynamicUserAccount =
					m["is_dynamic_user_account"].(bool)
			}
		case "kubernetes_secret":
			userAccount.AuthScheme = &AuthScheme{
				KubernetesSecret: &AuthSchemeKubernetesSecret{
					SecretName: m["secret_name"].(string),
					SecretKey:  m["secret_key"].(string),
				},
			}
		case "gcp_secrets_manager":
			userAccount.AuthScheme = &AuthScheme{
				GCPSecretManager: &AuthSchemeGCPSecretManager{
					SecretName: m["secret_name"].(string),
				},
			}
		default:
			return fmt.Errorf("unexpected auth_scheme [%s]", k)
		}
	}
	return nil
}

var ReadRepositoryUserAccountConfig = ResourceOperationConfig{
	Name:       "RepositoryReadUserAccount",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(
			"https://%s/v1/repos/%s/userAccounts/%s",
			c.ControlPlane,
			d.Get("repository_id").(string),
			d.Id(),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData {
		return &UserAccountResource{}
	},
}

func resourceRepositoryUserAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RepositoryUserAccountCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/userAccounts",
						c.ControlPlane,
						d.Get("repository_id").(string),
					)
				},
				NewResourceData: func() ResourceData {
					return &UserAccountResource{}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData {
					return &CreateUserAccountResponse{}
				},
			},
			ReadRepositoryUserAccountConfig,
		),
		ReadContext: ReadResource(ReadRepositoryUserAccountConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RepositoryUserAccountUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/userAccounts/%s",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Id(),
					)
				},
				NewResourceData: func() ResourceData {
					return &UserAccountResource{}
				},
			},
			ReadRepositoryUserAccountConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RepositoryUserAccountDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/userAccounts/%s",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Id(),
					)
				},
			},
		),
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				i interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := unmarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to unmarshal ID: %v",
						err,
					)
				}
				repositoryID := ids[0]
				userAccountID := ids[1]
				err = d.Set("repository_id", repositoryID)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to set 'repository_id': %v",
						err,
					)
				}
				d.SetId(userAccountID)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"repository_id": {
				Description: "ID of the repository.",
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
			},

			"name": {
				Description: "The name of the User Account.",
				Required:    true,
				Type:        schema.TypeString,
			},

			"auth_database_name": {
				Description: "For User Accounts associated with a MongoDB cyral_repository, " +
					"This value is only meaningful for User Accounts associated with " +
					"MongoDB.",
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
							Description: "If true, it is possible to create an approval with an automatic grant" +
								"associated with this User Account.",
							Optional: true,
							Type:     schema.TypeBool,
						},
						"max_auto_grant_duration": {
							Description: "The maximum duration in seconds for the lifespan of the approval, " +
								"if automatic_grant is true. E.g.: '2000s', '3000.5s'",
							Optional: true,
							Type:     schema.TypeString,
						},
					},
				},
			},

			"auth_scheme": {
				Description: "Credential option. List of supported types: " +
					supportedTypesMarkdown(allAuthSchemes),
				Required: true,
				Type:     schema.TypeSet,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_variable": {
							Description: "Credential option to set the repository user account from " +
								"Environment Variable.",
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							MinItems: 1,
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
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_arn": {
										Description: "The AWS IAM roleARN to gain access to the database.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
						"aws_secrets_manager": {
							Description: "Credential option to set the repository user account from " +
								"AWS Secrets Manager.",
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							MinItems: 1,
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
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							MinItems: 1,
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
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							MinItems: 1,
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
										Optional: true,
									},
								},
							},
						},
						"kubernetes_secret": {
							Description: "Credential option to set the repository user account from " +
								"a Kubernetes secret.",
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							MinItems: 1,
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
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							MinItems: 1,
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
					},
				},
			},
		},
	}
}
