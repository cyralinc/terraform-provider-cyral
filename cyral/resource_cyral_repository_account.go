package cyral

import (
	"errors"
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
			"repo_account":  resource.RepoAccount,
			"role_arn":      resource.RoleArn,
		},
	})
}

func (resource *AwsIAMResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("aws_iam").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["repo_account"].(string)
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
			"repo_account":  resource.RepoAccount,
			"secret_arn":    resource.SecretArn,
		},
	})
}

func (resource *AwsSecretsResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("aws_secrets_manager").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["repo_account"].(string)
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
			"repo_account":  resource.RepoAccount,
			"password":      resource.Password,
		},
	})
}

func (resource *CyralStorageResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("cyral_storage").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["repo_account"].(string)
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
			"repo_account":  resource.RepoAccount,
			"path":          resource.Path,
		},
	})
}

func (resource *HashicorpVaultResource) ReadFromSchema(d *schema.ResourceData) {
	data := d.Get("hashicorp_vault").(*schema.Set)

	for _, id := range data.List() {
		idMap := id.(map[string]interface{})

		resource.DatabaseName = idMap["database_name"].(string)
		resource.RepoAccount = idMap["repo_account"].(string)
		resource.Path = idMap["path"].(string)
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

type RepositoryAccountResource struct {
	RepoID            *string                 `json:"-"`
	AwsIAM            *AwsIAMResource         `json:"awsIAM,omitempty"`
	AwsSecretsManager *AwsSecretsResource     `json:"awsSecretsManager,omitempty"`
	CyralStorage      *CyralStorageResource   `json:"cyralStorage,omitempty"`
	HashicorpVault    *HashicorpVaultResource `json:"hashicorpVault,omitempty"`
}

func (repoAccount RepositoryAccountResource) isValid() error {
	numberOfAuthSchemas := 0

	log.Printf("[DEBUG] TESTE REPOS ACCOUNT: %v", repoAccount)

	if repoAccount.AwsIAM != nil {
		numberOfAuthSchemas += 1
	}

	if repoAccount.AwsSecretsManager != nil {
		numberOfAuthSchemas += 1
	}

	if repoAccount.CyralStorage != nil {
		numberOfAuthSchemas += 1
	}

	if repoAccount.HashicorpVault != nil {
		numberOfAuthSchemas += 1
	}

	if numberOfAuthSchemas != 1 {
		return errors.New("there are multiple auth schemas")
	}

	return nil
}

func (repoAccount RepositoryAccountResource) WriteToSchema(d *schema.ResourceData) {

	log.Printf("[DEBUG] RepositoryAccountResource - WriteToSchema START")

	if err := repoAccount.isValid(); err != nil {
		panic(err)
	}

	if repoAccount.RepoID != nil {
		d.Set("repo_id", repoAccount.RepoID)
	}

	if repoAccount.AwsIAM != nil {
		repoAccount.AwsIAM.WriteToSchema(d)

	} else if repoAccount.AwsSecretsManager != nil {
		repoAccount.AwsSecretsManager.WriteToSchema(d)

	} else if repoAccount.CyralStorage != nil {
		repoAccount.CyralStorage.WriteToSchema(d)

	} else if repoAccount.HashicorpVault != nil {
		repoAccount.HashicorpVault.WriteToSchema(d)
	}

	log.Printf("[DEBUG] RepositoryAccountResource - WriteToSchema END")
}

func (repoAccount *RepositoryAccountResource) ReadFromSchema(d *schema.ResourceData) {

	log.Printf("[DEBUG] RepositoryAccountResource - ReadFromSchema START")

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

	if data, hasRepoId := d.GetOk("repo_id"); hasRepoId {
		repoId := data.(string)
		repoAccount.RepoID = &repoId
	}
	log.Printf("[DEBUG] REPO ID: %s", *repoAccount.RepoID)

	if err := repoAccount.isValid(); err != nil {
		panic(err)
	}
	log.Printf("[DEBUG] RepositoryAccountResource - ReadFromSchema END")
}

var ReadRepositoryAccountConfig = ResourceOperationConfig{
	Name:       "RepositoryAccountResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		repo_id := d.Get("repo_id")
		return fmt.Sprintf("https://%s/v1/repos/%s/repoAccounts/%s", c.ControlPlane, repo_id, d.Id())
	},
	ResponseData: &RepositoryAccountResource{},
}

func ResourceRepositoryAccount() *schema.Resource {

	awsIAMSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"repo_account": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"role_arn": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
			},
		},
	}

	awsSecretsSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"repo_account": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"secret_arn": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
			},
		},
	}

	cyralStorageSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"repo_account": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"password": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
			},
		},
	}

	hashicorpVaultSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database_name": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"repo_account": {
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

	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RepositoryAccountResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					repo_id := d.Get("repo_id").(string)
					return fmt.Sprintf("https://%s/v1/repos/%s/repoAccounts", c.ControlPlane, repo_id)
				},
				ResourceData: &RepositoryAccountResource{},
				ResponseData: &CreateRepoAccountResponse{},
			}, ReadRepositoryAccountConfig,
		),
		ReadContext: ReadResource(ReadRepositoryAccountConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RepositoryAccountResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					repo_id := d.Get("repo_id").(string)
					return fmt.Sprintf("https://%s/v1/repos/%s/repoAccounts/%s", c.ControlPlane, repo_id, d.Id())
				},
				ResourceData: &RepositoryAccountResource{},
			}, ReadRepositoryAccountConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RepositoryAccountResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					repo_id := d.Get("repo_id").(string)
					return fmt.Sprintf("https://%s/v1/repos/%s/repoAccounts/%s", c.ControlPlane, repo_id, d.Id())
				},
			},
		),
		Schema: map[string]*schema.Schema{
			"repo_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_iam":             awsIAMSchema,
			"aws_secrets_manager": awsSecretsSchema,
			"cyral_storage":       cyralStorageSchema,
			"hashicorp_vault":     hashicorpVaultSchema,
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
