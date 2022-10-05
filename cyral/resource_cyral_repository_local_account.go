package cyral

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AwsIAMResource struct {
	DatabaseName string `json:"databaseName"`
	RepoAccount  string `json:"repoAccount"`
	RoleArn      string `json:"roleARN"`
}

func (resource AwsIAMResource) WriteToSchema(d *schema.ResourceData) {
	d.Set("aws_iam", []interface{}{
		map[string]interface{}{
			"database_name": resource.DatabaseName,
			"local_account": resource.RepoAccount,
			"role_arn":      resource.RoleArn,
		},
	})
}

func (resource *AwsIAMResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("aws_iam").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["local_account"].(string)
		resource.RoleArn = idMap["role_arn"].(string)
	}
}

type AwsSecretsManagerResource struct {
	DatabaseName string `json:"databaseName"`
	RepoAccount  string `json:"repoAccount"`
	SecretArn    string `json:"secretARN"`
}

func (resource AwsSecretsManagerResource) WriteToSchema(d *schema.ResourceData) {
	d.Set("aws_secrets_manager", []interface{}{
		map[string]interface{}{
			"database_name": resource.DatabaseName,
			"local_account": resource.RepoAccount,
			"secret_arn":    resource.SecretArn,
		},
	})
}

func (resource *AwsSecretsManagerResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("aws_secrets_manager").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["local_account"].(string)
		resource.SecretArn = idMap["secret_arn"].(string)
	}
}

type CyralStorageResource struct {
	DatabaseName string `json:"databaseName"`
	RepoAccount  string `json:"repoAccount"`
	Password     string `json:"password"`
}

func (resource CyralStorageResource) WriteToSchema(d *schema.ResourceData) {
	d.Set("cyral_storage", []interface{}{
		map[string]interface{}{
			"database_name": resource.DatabaseName,
			"local_account": resource.RepoAccount,
			"password":      resource.Password,
		},
	})
}

func (resource *CyralStorageResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("cyral_storage").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["local_account"].(string)
		resource.Password = idMap["password"].(string)
	}
}

type HashicorpVaultResource struct {
	DatabaseName string `json:"databaseName"`
	RepoAccount  string `json:"repoAccount"`
	Path         string `json:"path"`
}

func (resource HashicorpVaultResource) WriteToSchema(d *schema.ResourceData) {
	d.Set("hashicorp_vault", []interface{}{
		map[string]interface{}{
			"database_name": resource.DatabaseName,
			"local_account": resource.RepoAccount,
			"path":          resource.Path,
		},
	})
}

func (resource *HashicorpVaultResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("hashicorp_vault").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["local_account"].(string)
		resource.Path = idMap["path"].(string)
	}
}

type EnvironmentVariableResource struct {
	DatabaseName string `json:"databaseName"`
	RepoAccount  string `json:"repoAccount"`
	VariableName string `json:"variableName"`
}

func (resource EnvironmentVariableResource) WriteToSchema(d *schema.ResourceData) {
	_, useDeprecatedArgument := d.GetOk("enviroment_variable")

	argument := "environment_variable"
	if useDeprecatedArgument {
		// Deprecated: should be removed in the next MAJOR release
		argument = "enviroment_variable"
	}

	d.Set(argument, []interface{}{
		map[string]interface{}{
			"database_name": resource.DatabaseName,
			"local_account": resource.RepoAccount,
			"variable_name": resource.VariableName,
		},
	})
}

func (resource *EnvironmentVariableResource) ReadFromSchema(d *schema.ResourceData) {
	_, useDeprecatedArgument := d.GetOk("enviroment_variable")

	argument := "environment_variable"
	if useDeprecatedArgument {
		// Deprecated: should be removed in the next MAJOR release
		argument = "enviroment_variable"
	}

	data := d.Get(argument).(*schema.Set)
	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["local_account"].(string)
		resource.VariableName = idMap["variable_name"].(string)
	}
}

type KubernetesSecretResource struct {
	DatabaseName string `json:"databaseName"`
	RepoAccount  string `json:"repoAccount"`
	SecretName   string `json:"secretName"`
	SecretKey    string `json:"secretKey"`
}

func (resource KubernetesSecretResource) WriteToSchema(d *schema.ResourceData) {
	d.Set("kubernetes_secret", []interface{}{
		map[string]interface{}{
			"database_name": resource.DatabaseName,
			"local_account": resource.RepoAccount,
			"secret_name":   resource.SecretName,
			"secret_key":    resource.SecretKey,
		},
	})
}

func (resource *KubernetesSecretResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("kubernetes_secret").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["local_account"].(string)
		resource.SecretName = idMap["secret_name"].(string)
		resource.SecretKey = idMap["secret_key"].(string)
	}
}

type GcpSecretManagerResource struct {
	DatabaseName string `json:"databaseName"`
	RepoAccount  string `json:"repoAccount"`
	SecretName   string `json:"secretName"`
}

func (resource GcpSecretManagerResource) WriteToSchema(d *schema.ResourceData) {
	d.Set("gcp_secret_manager", []interface{}{
		map[string]interface{}{
			"database_name": resource.DatabaseName,
			"local_account": resource.RepoAccount,
			"secret_name":   resource.SecretName,
		},
	})
}

func (resource *GcpSecretManagerResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("gcp_secret_manager").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["local_account"].(string)
		resource.SecretName = idMap["secret_name"].(string)
	}
}

type CreateRepoAccountResponse struct {
	UUID string `json:"uuid"`
}

func (resource CreateRepoAccountResponse) WriteToSchema(d *schema.ResourceData) error {
	// TODO (next MAJOR): set ID to be of the format
	// {repositoryID}/{localAccountID}, to facilitate importing. -aholmquist 2022-08-01
	d.SetId(resource.UUID)
	return nil
}

func (resource *CreateRepoAccountResponse) ReadFromSchema(d *schema.ResourceData) error {
	resource.UUID = d.Id()
	return nil
}

type RepoAccountConfig struct {
	AutoApproveAccess      bool   `json:"autoApproveAccess,omitempty"`
	MaxAutoApproveDuration string `json:"maxAutoApproveDuration,omitempty"`
}

type RepositoryLocalAccountResource struct {
	Config              *RepoAccountConfig           `json:"config,omitempty"`
	AwsIAM              *AwsIAMResource              `json:"awsIAM,omitempty"`
	AwsSecretsManager   *AwsSecretsManagerResource   `json:"awsSecretsManager,omitempty"`
	CyralStorage        *CyralStorageResource        `json:"cyralStorage,omitempty"`
	HashicorpVault      *HashicorpVaultResource      `json:"hashicorpVault,omitempty"`
	EnvironmentVariable *EnvironmentVariableResource `json:"environmentVariable,omitempty"`
	KubernetesSecret    *KubernetesSecretResource    `json:"kubernetesSecret,omitempty"`
	GcpSecretManager    *GcpSecretManagerResource    `json:"gcpSecretManager,omitempty"`
}

func (repoAccount RepositoryLocalAccountResource) WriteToSchema(d *schema.ResourceData) error {
	log.Printf("[DEBUG] RepositoryLocalAccountResource - WriteToSchema START")

	// Optional `config` field
	if repoAccount.Config != nil {
		config := repoAccount.Config
		configMap := make(map[string]interface{})
		configMap["auto_approve_access"] = config.AutoApproveAccess
		configMap["max_auto_approve_duration"] = config.MaxAutoApproveDuration
		d.Set("config", []interface{}{configMap})
	}

	if repoAccount.AwsIAM != nil {
		repoAccount.AwsIAM.WriteToSchema(d)
	} else if repoAccount.AwsSecretsManager != nil {
		repoAccount.AwsSecretsManager.WriteToSchema(d)
	} else if repoAccount.CyralStorage != nil {
		repoAccount.CyralStorage.WriteToSchema(d)
	} else if repoAccount.HashicorpVault != nil {
		repoAccount.HashicorpVault.WriteToSchema(d)
	} else if repoAccount.EnvironmentVariable != nil {
		repoAccount.EnvironmentVariable.WriteToSchema(d)
	} else if repoAccount.KubernetesSecret != nil {
		repoAccount.KubernetesSecret.WriteToSchema(d)
	} else if repoAccount.GcpSecretManager != nil {
		repoAccount.GcpSecretManager.WriteToSchema(d)
	}

	log.Printf("[DEBUG] RepositoryLocalAccountResource - WriteToSchema END")

	return nil
}

func (repoAccount *RepositoryLocalAccountResource) ReadFromSchema(d *schema.ResourceData) error {
	log.Printf("[DEBUG] RepositoryLocalAccountResource - ReadFromSchema START")

	// `config` optional field
	configList := d.Get("config").(*schema.Set).List()
	if len(configList) > 0 {
		configMap := configList[0].(map[string]interface{})
		repoAccount.Config = &RepoAccountConfig{
			AutoApproveAccess:      configMap["auto_approve_access"].(bool),
			MaxAutoApproveDuration: configMap["max_auto_approve_duration"].(string),
		}
	}

	if _, hasAwsIam := d.GetOk("aws_iam"); hasAwsIam {
		repoAccount.AwsIAM = &AwsIAMResource{}
		repoAccount.AwsIAM.ReadFromSchema(d)
	} else if _, hasAwsSecretsManager := d.GetOk("aws_secrets_manager"); hasAwsSecretsManager {
		repoAccount.AwsSecretsManager = &AwsSecretsManagerResource{}
		repoAccount.AwsSecretsManager.ReadFromSchema(d)
	} else if _, hasCyralStorage := d.GetOk("cyral_storage"); hasCyralStorage {
		repoAccount.CyralStorage = &CyralStorageResource{}
		repoAccount.CyralStorage.ReadFromSchema(d)
	} else if _, hasHashicorpVault := d.GetOk("hashicorp_vault"); hasHashicorpVault {
		repoAccount.HashicorpVault = &HashicorpVaultResource{}
		repoAccount.HashicorpVault.ReadFromSchema(d)
	} else if _, hasDeprecatedEnvVar := d.GetOk("enviroment_variable"); hasDeprecatedEnvVar {
		// Deprecated: should be removed in the next MAJOR version
		repoAccount.EnvironmentVariable = &EnvironmentVariableResource{}
		repoAccount.EnvironmentVariable.ReadFromSchema(d)
	} else if _, hasEnvironmentVariable := d.GetOk("environment_variable"); hasEnvironmentVariable {
		repoAccount.EnvironmentVariable = &EnvironmentVariableResource{}
		repoAccount.EnvironmentVariable.ReadFromSchema(d)
	} else if _, hasKubernetesSecret := d.GetOk("kubernetes_secret"); hasKubernetesSecret {
		repoAccount.KubernetesSecret = &KubernetesSecretResource{}
		repoAccount.KubernetesSecret.ReadFromSchema(d)
	} else if _, hasGcpSecretManager := d.GetOk("gcp_secret_manager"); hasGcpSecretManager {
		repoAccount.GcpSecretManager = &GcpSecretManagerResource{}
		repoAccount.GcpSecretManager.ReadFromSchema(d)
	}

	log.Printf("[DEBUG] RepositoryLocalAccountResource - ReadFromSchema END")

	return nil
}

func resourceRepositoryLocalAccount() *schema.Resource {
	ReadRepositoryLocalAccountConfig := ResourceOperationConfig{
		Name:       "RepositoryLocalAccountResourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			repository_id := d.Get("repository_id")
			return fmt.Sprintf(
				"https://%s/v1/repos/%s/repoAccounts/%s",
				c.ControlPlane, repository_id, d.Id(),
			)
		},
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &RepositoryLocalAccountResource{} },
	}

	secretManagersTypes := []string{
		"aws_iam",
		"aws_secrets_manager",
		"cyral_storage",
		"hashicorp_vault",
		"enviroment_variable", // Deprecated: should be removed in the next MAJOR release
		"environment_variable",
		"kubernetes_secret",
		"gcp_secret_manager",
	}

	databaseNameDescription := "Database name that the local account corresponds to."
	localAccountDescription := "Local account name."

	awsIAMSchema := &schema.Schema{
		Description:  "Credential option to set the local account from AWS IAM.",
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Description: databaseNameDescription,
					Type:        schema.TypeString,
					Optional:    true,
				},
				"local_account": {
					Description: localAccountDescription,
					Type:        schema.TypeString,
					Required:    true,
				},
				"role_arn": {
					Description: "AWS IAM role ARN.",
					Type:        schema.TypeString,
					Required:    true,
				},
			},
		},
	}

	awsSecretsManagerSchema := &schema.Schema{
		Description:  "Credential option to set the local account from AWS Secrets Manager.",
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Description: databaseNameDescription,
					Type:        schema.TypeString,
					Optional:    true,
				},
				"local_account": {
					Description: localAccountDescription,
					Type:        schema.TypeString,
					Required:    true,
				},
				"secret_arn": {
					Description: "ARN of the AWS Secret Manager that stores the credential.",
					Type:        schema.TypeString,
					Required:    true,
				},
			},
		},
	}

	cyralStorageSchema := &schema.Schema{
		Description:  "Credential option to set the local account from Cyral Storage.",
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Description: databaseNameDescription,
					Type:        schema.TypeString,
					Optional:    true,
				},
				"local_account": {
					Description: localAccountDescription,
					Type:        schema.TypeString,
					Required:    true,
				},
				"password": {
					Description: "Local account password.",
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
				},
			},
		},
	}

	hashicorpVaultSchema := &schema.Schema{
		Description:  "Credential option to set the local account from Hashicorp Vault.",
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Description: databaseNameDescription,
					Type:        schema.TypeString,
					Optional:    true,
				},
				"local_account": {
					Description: localAccountDescription,
					Type:        schema.TypeString,
					Required:    true,
				},
				"path": {
					Description: "Hashicorp Vault path.",
					Type:        schema.TypeString,
					Required:    true,
				},
			},
		},
	}

	environmentVariableSchema := &schema.Schema{
		Description:  "Credential option to set the local account from Environment Variable.",
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Description: databaseNameDescription,
					Type:        schema.TypeString,
					Optional:    true,
				},
				"local_account": {
					Description: localAccountDescription,
					Type:        schema.TypeString,
					Required:    true,
				},
				"variable_name": {
					Description: "Name of the environment variable that will store credentials.",
					Type:        schema.TypeString,
					Required:    true,
				},
			},
		},
	}

	// Deprecated: should be removed in the next MAJOR release
	environmentVariableSchemaDeprecated := *environmentVariableSchema
	environmentVariableSchemaDeprecated.Deprecated = "This argument is deprecated, use " +
		"'environment_variable' instead."

	kubernetesSecretSchema := &schema.Schema{
		Description:  "Credential option to set the local account from Kubernetes Secret.",
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Description: databaseNameDescription,
					Type:        schema.TypeString,
					Optional:    true,
				},
				"local_account": {
					Description: localAccountDescription,
					Type:        schema.TypeString,
					Required:    true,
				},
				"secret_name": {
					Description: "Name of the secret in kubernetes.",
					Type:        schema.TypeString,
					Required:    true,
				},
				"secret_key": {
					Description: "Name of the key that stores the credentials within the secret.",
					Type:        schema.TypeString,
					Required:    true,
				},
			},
		},
	}

	gcpSecretManagerSchema := &schema.Schema{
		Description:  "Credential option to set the local account from GCP Secret Manager.",
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Description: databaseNameDescription,
					Type:        schema.TypeString,
					Optional:    true,
				},
				"local_account": {
					Description: localAccountDescription,
					Type:        schema.TypeString,
					Required:    true,
				},
				"secret_name": {
					Description: "The unique identifier of the secret in GCP Secret Manager. Should obey one of the following formats: `projects/<project-name>/secrets/<secret-name>` or `projects/<project-name>/secrets/<secret-name>/versions/<version>`.",
					Type:        schema.TypeString,
					Required:    true,
				},
			},
		},
	}

	return &schema.Resource{
		Description: "Manages [repository local accounts](https://cyral.com/docs/using-cyral/sso-auth-users#give-your-sidecar-access-to-the-local-account).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RepositoryLocalAccountResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					repository_id := d.Get("repository_id").(string)
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/repoAccounts",
						c.ControlPlane, repository_id,
					)
				},
				NewResourceData: func() ResourceData { return &RepositoryLocalAccountResource{} },
				NewResponseData: func(_ *schema.ResourceData) ResponseData { return &CreateRepoAccountResponse{} },
			}, ReadRepositoryLocalAccountConfig,
		),
		ReadContext: ReadResource(ReadRepositoryLocalAccountConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RepositoryLocalAccountResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					repository_id := d.Get("repository_id").(string)
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/repoAccounts/%s",
						c.ControlPlane, repository_id, d.Id(),
					)
				},
				NewResourceData: func() ResourceData { return &RepositoryLocalAccountResource{} },
			}, ReadRepositoryLocalAccountConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RepositoryLocalAccountResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					repository_id := d.Get("repository_id").(string)
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/repoAccounts/%s",
						c.ControlPlane, repository_id, d.Id(),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"repository_id": {
				Description: "ID of the repository that will be used by the local account.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"config": {
				Description: "Optional configuration for automatic approvals.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_approve_access": {
							Description: "If true, enables automatic approvals for this local account.",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"max_auto_approve_duration": {
							Description: `Maximum duration (in seconds) for which approvals can be automatically granted, following ISO 8601 time format. For example, "PT1H2M3S" indicates 1 hour, 2 minutes and 3 seconds. "PT4S" denotes 4 seconds.`,
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "P0D",
						},
					},
				},
			},

			"aws_iam":              awsIAMSchema,
			"aws_secrets_manager":  awsSecretsManagerSchema,
			"cyral_storage":        cyralStorageSchema,
			"hashicorp_vault":      hashicorpVaultSchema,
			"enviroment_variable":  &environmentVariableSchemaDeprecated,
			"environment_variable": environmentVariableSchema,
			"kubernetes_secret":    kubernetesSecretSchema,
			"gcp_secret_manager":   gcpSecretManagerSchema,
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := unmarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, err
				}
				repositoryID := ids[0]
				localAccountID := ids[1]
				d.Set("repository_id", repositoryID)
				d.SetId(localAccountID)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
