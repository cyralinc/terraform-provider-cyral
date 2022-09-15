package cyral

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var allAuthSchemes = []string{
	"environment_variable",
}

type AuthSchemeEnvironmentVariable struct {
	VariableName string `json:"variableName"`
}

type AuthScheme struct {
	EnvironmentVariable *AuthSchemeEnvironmentVariable `json:"environmentVariable"`
}

type ApprovalConfig struct {
	AutomaticGrant            bool `json:"automaticGrant"`
	MaxAutomaticGrantDuration int  `json:"maxAutomaticGrantDuration"`
}

type UserAccountConfig struct {
	Approval *ApprovalConfig `json:"approvalConfig"`
}

type UserAccountResource struct {
	UserAccountID    string             `json:"userAccountID"`
	Name             string             `json:"name"`
	AuthDatabaseName string             `json:"authDatabaseName"`
	AuthScheme       *AuthScheme        `json:"authScheme"`
	Config           *UserAccountConfig `json:"config"`
}

type CreateUserAccountResponse struct {
	UserAccountID string `json:"userAccountID"`
}

func (resp *CreateUserAccountResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(resp.UserAccountID)
	return nil
}

func (response *UserAccountResource) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(response.UserAccountID)
	return nil
}

// ReadFromSchema is used to translate a .tf file into whatever the
// UserAccounts API expects.
func (userAccount *UserAccountResource) ReadFromSchema(d *schema.ResourceData) error {
	userAccount.Name = d.Get("name").(string)

	authScheme := d.Get("auth_scheme").(*schema.Set).List()
	if len(authScheme) < 0 {
		return fmt.Errorf("auth_scheme is required")
	}
	if len(authScheme) > 1 {
		return fmt.Errorf("only one auth_scheme must be supplied")
	}

	authSchemes := authScheme[0].(map[string]interface{})
	for k, v := range authSchemes {
		switch k {
		case "environment_variable":
			m := v.(*schema.Set).List()[0].(map[string]interface{})
			userAccount.AuthScheme = &AuthScheme{
				EnvironmentVariable: &AuthSchemeEnvironmentVariable{
					VariableName: m["variable_name"].(string),
				},
			}
			break

		default:
			return fmt.Errorf("unexpected auth_scheme [%s]", k)
		}
	}
	return nil
}

var ReadRepositoryUserAccountConfig = ResourceOperationConfig{
	Name:       "RepositoryReadDatabaseAccount",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		log.Printf("[DEBUG] RIREY the id is: %s", d.Id())
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

func resourceRepositoryDatabaseAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RepositoryDatabaseAccountCreate",
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
				Name:       "RepositoryDatabaseAccountUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/userAccounts/%s",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Get("user_account_id").(string),
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
				Name:       "RepositoryDatabaseAccountDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/userAccounts/%s",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Get("user_account_id").(string),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment.",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"repository_id": {
				Description: "ID of the repository.",
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
			},

			"name": {
				Description: "The name of the database account.",
				Required:    true,
				Type:        schema.TypeString,
			},

			"auth_scheme": {
				Description: "Credential option. List of supported types: " +
					supportedTypesMarkdown(allAuthSchemes),
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_variable": {
							Description: "Credential option to set the local account from " +
								"Environment Variable.",
							Optional: true,
							Type:     schema.TypeSet,
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
					},
				},
			},
		},
	}
}
