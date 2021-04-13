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

type AwsSecretsResource struct {
	DatabaseName string `json:"databaseName"`
	RepoAccount  string `json:"repoAccount"`
	SecretArn    string `json:"secretARN"`
}

func (resource AwsSecretsResource) WriteToSchema(d *schema.ResourceData) {
	d.Set("aws_secrets_manager", []interface{}{
		map[string]interface{}{
			"database_name": resource.DatabaseName,
			"local_account": resource.RepoAccount,
			"secret_arn":    resource.SecretArn,
		},
	})
}

func (resource *AwsSecretsResource) ReadFromSchema(d *schema.ResourceData) {
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

type EnviromentVariableResource struct {
	DatabaseName string `json:"databaseName"`
	RepoAccount  string `json:"repoAccount"`
	VariableName string `json:"variableName"`
}

func (resource EnviromentVariableResource) WriteToSchema(d *schema.ResourceData) {
	d.Set("enviroment_variable", []interface{}{
		map[string]interface{}{
			"database_name": resource.DatabaseName,
			"local_account": resource.RepoAccount,
			"variable_name": resource.VariableName,
		},
	})
}

func (resource *EnviromentVariableResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("enviroment_variable").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["local_account"].(string)
		resource.VariableName = idMap["variable_name"].(string)
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
	RepoID             *string                     `json:"-"`
	AwsIAM             *AwsIAMResource             `json:"awsIAM,omitempty"`
	AwsSecretsManager  *AwsSecretsResource         `json:"awsSecretsManager,omitempty"`
	CyralStorage       *CyralStorageResource       `json:"cyralStorage,omitempty"`
	HashicorpVault     *HashicorpVaultResource     `json:"hashicorpVault,omitempty"`
	EnviromentVariable *EnviromentVariableResource `json:"environmentVariable,omitempty"`
}

func (repoAccount RepositoryLocalAccountResource) WriteToSchema(d *schema.ResourceData) {
	log.Printf("[DEBUG] RepositoryLocalAccountResource - WriteToSchema START")

	if repoAccount.RepoID != nil {
		d.Set("repository_id", repoAccount.RepoID)
	}

	if repoAccount.AwsIAM != nil {
		repoAccount.AwsIAM.WriteToSchema(d)

	} else if repoAccount.AwsSecretsManager != nil {
		repoAccount.AwsSecretsManager.WriteToSchema(d)

	} else if repoAccount.CyralStorage != nil {
		repoAccount.CyralStorage.WriteToSchema(d)

	} else if repoAccount.HashicorpVault != nil {
		repoAccount.HashicorpVault.WriteToSchema(d)
	} else if repoAccount.EnviromentVariable != nil {
		repoAccount.EnviromentVariable.WriteToSchema(d)
	}

	log.Printf("[DEBUG] RepositoryLocalAccountResource - WriteToSchema END")
}

func (repoAccount *RepositoryLocalAccountResource) ReadFromSchema(d *schema.ResourceData) {
	log.Printf("[DEBUG] RepositoryLocalAccountResource - ReadFromSchema START")

	if _, hasAwsIam := d.GetOk("aws_iam"); hasAwsIam {
		repoAccount.AwsIAM = &AwsIAMResource{}
		repoAccount.AwsIAM.ReadFromSchema(d)
	}

	if _, hasAwsSecretsManager := d.GetOk("aws_secrets_manager"); hasAwsSecretsManager {
		repoAccount.AwsSecretsManager = &AwsSecretsResource{}
		repoAccount.AwsSecretsManager.ReadFromSchema(d)
	}

	if _, hasCyralStorage := d.GetOk("cyral_storage"); hasCyralStorage {
		repoAccount.CyralStorage = &CyralStorageResource{}
		repoAccount.CyralStorage.ReadFromSchema(d)
	}

	if _, hasHashicorpVault := d.GetOk("hashicorp_vault"); hasHashicorpVault {
		repoAccount.HashicorpVault = &HashicorpVaultResource{}
		repoAccount.HashicorpVault.ReadFromSchema(d)
	}

	if _, hasEnviromentVariable := d.GetOk("enviroment_variable"); hasEnviromentVariable {
		repoAccount.EnviromentVariable = &EnviromentVariableResource{}
		repoAccount.EnviromentVariable.ReadFromSchema(d)
	}

	if data, hasRepoId := d.GetOk("repository_id"); hasRepoId {
		repoId := data.(string)
		repoAccount.RepoID = &repoId
	}

	log.Printf("[DEBUG] RepositoryLocalAccountResource - ReadFromSchema END")
}

var ReadRepositoryLocalAccountConfig = ResourceOperationConfig{
	Name:       "RepositoryLocalAccountResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		repository_id := d.Get("repository_id")
		return fmt.Sprintf("https://%s/v1/repos/%s/repoAccounts/%s", c.ControlPlane, repository_id, d.Id())
	},
	ResponseData: &RepositoryLocalAccountResource{},
}

func resourceRepositoryLocalAccount() *schema.Resource {
	awsIAMSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ConflictsWith: []string{
			"aws_secrets_manager",
			"cyral_storage",
			"hashicorp_vault",
			"enviroment_variable",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"local_account": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"role_arn": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
			},
		},
	}

	awsSecretsSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ConflictsWith: []string{
			"aws_iam",
			"cyral_storage",
			"hashicorp_vault",
			"enviroment_variable",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"local_account": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"secret_arn": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
			},
		},
	}

	cyralStorageSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ConflictsWith: []string{
			"aws_iam",
			"aws_secrets_manager",
			"hashicorp_vault",
			"enviroment_variable",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"local_account": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
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
		Type:     schema.TypeSet,
		Optional: true,
		ConflictsWith: []string{
			"aws_iam",
			"aws_secrets_manager",
			"cyral_storage",
			"enviroment_variable",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"local_account": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"path": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
			},
		},
	}

	enviromentVariableSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ConflictsWith: []string{
			"aws_iam",
			"aws_secrets_manager",
			"hashicorp_vault",
			"cyral_storage",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"local_account": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"variable_name": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
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
					return fmt.Sprintf("https://%s/v1/repos/%s/repoAccounts", c.ControlPlane, repository_id)
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
					return fmt.Sprintf("https://%s/v1/repos/%s/repoAccounts/%s", c.ControlPlane, repository_id, d.Id())
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
					return fmt.Sprintf("https://%s/v1/repos/%s/repoAccounts/%s", c.ControlPlane, repository_id, d.Id())
				},
			},
		),
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_iam":             awsIAMSchema,
			"aws_secrets_manager": awsSecretsSchema,
			"cyral_storage":       cyralStorageSchema,
			"hashicorp_vault":     hashicorpVaultSchema,
			"enviroment_variable": enviromentVariableSchema,
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
