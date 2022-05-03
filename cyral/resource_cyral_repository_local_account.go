package cyral

import (
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

	if useDeprecatedArgument { // Deprecated: should be removed in the next MAJOR release
		d.Set("enviroment_variable", []interface{}{
			map[string]interface{}{
				"database_name": resource.DatabaseName,
				"local_account": resource.RepoAccount,
				"variable_name": resource.VariableName,
			},
		})
	} else {
		d.Set("environment_variable", []interface{}{
			map[string]interface{}{
				"database_name": resource.DatabaseName,
				"local_account": resource.RepoAccount,
				"variable_name": resource.VariableName,
			},
		})
	}
}

func (resource *EnvironmentVariableResource) ReadFromSchema(d *schema.ResourceData) {
	deprecatedData, useDeprecatedArgument := d.GetOk("enviroment_variable")

	if useDeprecatedArgument { // Deprecated: should be removed in the next MAJOR release
		for _, id := range deprecatedData.(*schema.Set).List() {
			idMap := id.(map[string]interface{})

			resource.DatabaseName = idMap["database_name"].(string)
			resource.RepoAccount = idMap["local_account"].(string)
			resource.VariableName = idMap["variable_name"].(string)
		}
	} else {
		data := d.Get("environment_variable").(*schema.Set)

		for _, id := range data.List() {
			idMap := id.(map[string]interface{})

			resource.DatabaseName = idMap["database_name"].(string)
			resource.RepoAccount = idMap["local_account"].(string)
			resource.VariableName = idMap["variable_name"].(string)
		}
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

type CreateRepoAccountResponse struct {
	UUID string `json:"uuid"`
}

func (resource CreateRepoAccountResponse) WriteToSchema(d *schema.ResourceData) {
	d.SetId(resource.UUID)
}

func (resource *CreateRepoAccountResponse) ReadFromSchema(d *schema.ResourceData) {
	resource.UUID = d.Id()
}

type RepositoryLocalAccountResource struct {
	AwsIAM              *AwsIAMResource              `json:"awsIAM,omitempty"`
	AwsSecretsManager   *AwsSecretsManagerResource   `json:"awsSecretsManager,omitempty"`
	CyralStorage        *CyralStorageResource        `json:"cyralStorage,omitempty"`
	HashicorpVault      *HashicorpVaultResource      `json:"hashicorpVault,omitempty"`
	EnvironmentVariable *EnvironmentVariableResource `json:"environmentVariable,omitempty"`
	KubernetesSecret    *KubernetesSecretResource    `json:"kubernetesSecret,omitempty"`
}

func (repoAccount RepositoryLocalAccountResource) WriteToSchema(d *schema.ResourceData) {
	log.Printf("[DEBUG] RepositoryLocalAccountResource - WriteToSchema START")

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
	}

	log.Printf("[DEBUG] RepositoryLocalAccountResource - WriteToSchema END")
}

func (repoAccount *RepositoryLocalAccountResource) ReadFromSchema(d *schema.ResourceData) {
	log.Printf("[DEBUG] RepositoryLocalAccountResource - ReadFromSchema START")

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
	}

	log.Printf("[DEBUG] RepositoryLocalAccountResource - ReadFromSchema END")
}

var ReadRepositoryLocalAccountConfig = ResourceOperationConfig{
	Name:       "RepositoryLocalAccountResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		repository_id := d.Get("repository_id")
		return fmt.Sprintf(
			"https://%s/v1/repos/%s/repoAccounts/%s",
			c.ControlPlane, repository_id, d.Id(),
		)
	},
	ResponseData: &RepositoryLocalAccountResource{},
}

func resourceRepositoryLocalAccount() *schema.Resource {
	secretManagersTypes := []string{
		"aws_iam",
		"aws_secrets_manager",
		"cyral_storage",
		"hashicorp_vault",
		"enviroment_variable", // Deprecated: should be removed in the next MAJOR release
		"environment_variable",
		"kubernetes_secret",
	}

	awsIAMSchema := &schema.Schema{
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"local_account": {
					Type:     schema.TypeString,
					Required: true,
				},
				"role_arn": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	awsSecretsManagerSchema := &schema.Schema{
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"local_account": {
					Type:     schema.TypeString,
					Required: true,
				},
				"secret_arn": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	cyralStorageSchema := &schema.Schema{
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"local_account": {
					Type:     schema.TypeString,
					Required: true,
				},
				"password": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
			},
		},
	}

	hashicorpVaultSchema := &schema.Schema{
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"local_account": {
					Type:     schema.TypeString,
					Required: true,
				},
				"path": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	// Deprecated: should be removed in the next MAJOR release
	environmentVariableSchemaDeprecated := &schema.Schema{
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Deprecated:   "This argument is deprecated, use 'environment_variable' instead.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"local_account": {
					Type:     schema.TypeString,
					Required: true,
				},
				"variable_name": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	environmentVariableSchema := &schema.Schema{
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"local_account": {
					Type:     schema.TypeString,
					Required: true,
				},
				"variable_name": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	kubernetesSecretSchema := &schema.Schema{
		Type:         schema.TypeSet,
		Optional:     true,
		ExactlyOneOf: secretManagersTypes,
		MaxItems:     1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"local_account": {
					Type:     schema.TypeString,
					Required: true,
				},
				"secret_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"secret_key": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	return &schema.Resource{
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
				ResourceData: &RepositoryLocalAccountResource{},
				ResponseData: &CreateRepoAccountResponse{},
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
				ResourceData: &RepositoryLocalAccountResource{},
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
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aws_iam":              awsIAMSchema,
			"aws_secrets_manager":  awsSecretsManagerSchema,
			"cyral_storage":        cyralStorageSchema,
			"hashicorp_vault":      hashicorpVaultSchema,
			"enviroment_variable":  environmentVariableSchemaDeprecated,
			"environment_variable": environmentVariableSchema,
			"kubernetes_secret":    kubernetesSecretSchema,
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
