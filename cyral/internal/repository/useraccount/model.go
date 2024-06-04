package useraccount

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AuthScheme struct {
	AWSIAM              *AuthSchemeAWSIAM              `json:"awsIAM"`
	AWSSecretsManager   *AuthSchemeAWSSecretsManager   `json:"awsSecretsManager"`
	CyralStorage        *AuthSchemeCyralStorage        `json:"cyralStorage"`
	HashicorpVault      *AuthSchemeHashicorpVault      `json:"hashicorpVault"`
	EnvironmentVariable *AuthSchemeEnvironmentVariable `json:"environmentVariable"`
	KubernetesSecret    *AuthSchemeKubernetesSecret    `json:"kubernetesSecret"`
	GCPSecretManager    *AuthSchemeGCPSecretManager    `json:"gcpSecretManager"`
	AzureKeyVault       *AuthSchemeAzureKeyVault       `json:"azureKeyVault"`
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

type AuthSchemeAzureKeyVault struct {
	SecretURL string `json:"secretUrl,omitempty"`
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
	d.SetId(utils.MarshalComposedID(
		[]string{
			d.Get("repository_id").(string),
			resp.UserAccountID,
		},
		"/",
	))
	return nil
}

func (resource *UserAccountResource) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("user_account_id", resource.UserAccountID); err != nil {
		return fmt.Errorf("error setting 'user_account_id': %w", err)
	}

	if err := d.Set("name", resource.Name); err != nil {
		return fmt.Errorf("error setting 'name': %w", err)
	}
	if err := d.Set("auth_database_name", resource.AuthDatabaseName); err != nil {
		return fmt.Errorf("error setting 'auth_database_name': %w", err)
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
	switch {
	case resource.AuthScheme.AWSIAM != nil:
		authScheme = []interface{}{
			map[string]interface{}{
				"aws_iam": []interface{}{
					map[string]interface{}{
						"role_arn": resource.AuthScheme.AWSIAM.RoleARN,
					},
				},
			},
		}
	case resource.AuthScheme.AWSSecretsManager != nil:
		authScheme = []interface{}{
			map[string]interface{}{
				"aws_secrets_manager": []interface{}{
					map[string]interface{}{
						"secret_arn": resource.AuthScheme.AWSSecretsManager.SecretArn,
					},
				},
			},
		}
	case resource.AuthScheme.CyralStorage != nil:
		authScheme = []interface{}{
			map[string]interface{}{
				"cyral_storage": []interface{}{
					map[string]interface{}{
						"password": resource.AuthScheme.CyralStorage.Password,
					},
				},
			},
		}
	case resource.AuthScheme.EnvironmentVariable != nil:
		authScheme = []interface{}{
			map[string]interface{}{
				"environment_variable": []interface{}{
					map[string]interface{}{
						"variable_name": resource.AuthScheme.EnvironmentVariable.VariableName,
					},
				},
			},
		}
	case resource.AuthScheme.GCPSecretManager != nil:
		authScheme = []interface{}{
			map[string]interface{}{
				"gcp_secrets_manager": []interface{}{
					map[string]interface{}{
						"secret_name": resource.AuthScheme.GCPSecretManager.SecretName,
					},
				},
			},
		}
	case resource.AuthScheme.AzureKeyVault != nil:
		authScheme = []interface{}{
			map[string]interface{}{
				"azure_key_vault": []interface{}{
					map[string]interface{}{
						"secret_url": resource.AuthScheme.AzureKeyVault.SecretURL,
					},
				},
			},
		}
	case resource.AuthScheme.HashicorpVault != nil:
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
	case resource.AuthScheme.KubernetesSecret != nil:
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
	default:
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
	authSchemeSet := d.Get("auth_scheme").([]interface{})
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
					Path:                 m["path"].(string),
					IsDynamicUserAccount: m["is_dynamic_user_account"].(bool),
				},
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
		case "azure_key_vault":
			userAccount.AuthScheme = &AuthScheme{
				AzureKeyVault: &AuthSchemeAzureKeyVault{
					SecretURL: m["secret_url"].(string),
				},
			}
		default:
			return fmt.Errorf("unexpected auth_scheme [%s]", k)
		}
	}
	return nil
}
