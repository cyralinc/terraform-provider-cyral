package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceIntegrationOkta struct {
	Name         string   `json:"name"`
	Certificate  string   `json:"cert"`
	EmailDomains []string `json:"idpDomains"`
	SignInUrl    string   `json:"signInEndpoint"`
	SignOutUrl   string   `json:"signOutEndpoint"`
}

func (data ResourceIntegrationOkta) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("certificate", data.Certificate)
	d.Set("email_domains", data.EmailDomains)
	d.Set("signin_url", data.SignInUrl)
	d.Set("signout_url", data.SignOutUrl)
}

func (data *ResourceIntegrationOkta) ReadFromSchema(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.Certificate = d.Get("certificate").(string)

	keyMap := d.Get("email_domains").([]interface{})
	emailsDomain := []string{}

	for _, i := range keyMap {
		emailsDomain = append(emailsDomain, i.(string))
	}

	data.EmailDomains = emailsDomain
	data.SignInUrl = d.Get("signin_url").(string)
	data.SignOutUrl = d.Get("signout_url").(string)
}

type ResourceIntegrationOktaPayload struct {
	Samlp ResourceIntegrationOkta `json:"samlp"`
}

func (data ResourceIntegrationOktaPayload) WriteToSchema(d *schema.ResourceData) {
	data.Samlp.WriteToSchema(d)
}

func (data *ResourceIntegrationOktaPayload) ReadFromSchema(d *schema.ResourceData) {
	data.Samlp.ReadFromSchema(d)
}

type CreateResourceIntegrationOktaResponse struct {
	ID string `json:"status"`
}

func (data CreateResourceIntegrationOktaResponse) WriteToSchema(d *schema.ResourceData) {
	d.SetId("okta")
}

func (data *CreateResourceIntegrationOktaResponse) ReadFromSchema(d *schema.ResourceData) {
	data.ID = "okta"
}

type KeycloakProvider struct{}

type ResourceIntegrationOktaIdentityProviderPayload struct {
	Keycloak KeycloakProvider `json:"keyCloakProvider"`
}

func (data ResourceIntegrationOktaIdentityProviderPayload) WriteToSchema(d *schema.ResourceData) {}

func (data *ResourceIntegrationOktaIdentityProviderPayload) ReadFromSchema(d *schema.ResourceData) {}

var ReadResourceIntegrationOktaConfig = ResourceOperationConfig{
	Name:       " OktaResourceRead - Integration ",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/okta/%s", c.ControlPlane, d.Get("name").(string))
	},
	ResponseData: &ResourceIntegrationOktaPayload{},
}

var ReadResourceIntegrationOktaIdentityProviderConfig = ResourceOperationConfig{
	Name:       " OktaResourceRead - IdentityProvider ",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/conf/identityProviders/%s", c.ControlPlane, d.Get("name").(string))
	},
	ResponseData: &ResourceIntegrationOktaIdentityProviderPayload{},
}

var cleanUpOktaIntegration = ResourceOperationConfig{
	Name:       " OktaResourceDelete - Integration ",
	HttpMethod: http.MethodDelete,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/okta/%s", c.ControlPlane, d.Get("name").(string))
	},
}

func resourceIntegrationOkta() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use `cyral_integration_sso_okta` instead.",
		CreateContext:      CreateOktaIntegration,
		ReadContext:        ReadOktaIntegration,
		UpdateContext:      UpdateOktaIntegration,
		DeleteContext:      DeleteOktaIntegration,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"email_domains": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"signin_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"signout_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateOktaIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := CreateResource(
		ResourceOperationConfig{
			Name:       " OktaResourceCreate - Integration ",
			HttpMethod: http.MethodPost,
			CreateURL: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf("https://%s/v1/integrations/okta", c.ControlPlane)
			},
			ResourceData: &ResourceIntegrationOktaPayload{},
			ResponseData: &CreateResourceIntegrationOktaResponse{},
		}, ReadResourceIntegrationOktaConfig,
	)(ctx, d, m)

	if diag.HasError() {
		// Silent clean up
		_ = DeleteResource(cleanUpOktaIntegration)(ctx, d, m)
	} else {
		diag = CreateResource(
			ResourceOperationConfig{
				Name:       " OktaResourceCreate - IdentityProvider ",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/conf/identityProviders/%s", c.ControlPlane, d.Get("name").(string))
				},
				ResourceData: &ResourceIntegrationOktaIdentityProviderPayload{},
				ResponseData: &CreateResourceIntegrationOktaResponse{},
			}, ReadResourceIntegrationOktaIdentityProviderConfig,
		)(ctx, d, m)
	}

	return diag
}

func ReadOktaIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := ReadResource(
		ReadResourceIntegrationOktaConfig,
	)(ctx, d, m)

	if !diag.HasError() {
		diag = ReadResource(
			ReadResourceIntegrationOktaIdentityProviderConfig,
		)(ctx, d, m)
	}

	return diag
}

func UpdateOktaIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := UpdateResource(ResourceOperationConfig{
		Name:       " OktaResourceUpdate - Integration ",
		HttpMethod: http.MethodPut,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/okta", c.ControlPlane)
		},
		ResourceData: &ResourceIntegrationOktaPayload{},
	}, ReadResourceIntegrationOktaConfig,
	)(ctx, d, m)

	if !diag.HasError() {
		diag = UpdateResource(
			ResourceOperationConfig{
				Name:       " OktaResourceUpdate - IdentityProvider ",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/conf/identityProviders/%s", c.ControlPlane, d.Get("name").(string))
				},
				ResourceData: &ResourceIntegrationOktaIdentityProviderPayload{},
			}, ReadResourceIntegrationOktaIdentityProviderConfig,
		)(ctx, d, m)
	}

	return diag
}

func DeleteOktaIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := DeleteResource(cleanUpOktaIntegration)(ctx, d, m)

	if !diag.HasError() {
		diag = DeleteResource(
			ResourceOperationConfig{
				Name:       " OktaResourceDelete - IdendityProvider",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/conf/identityProviders/%s", c.ControlPlane, d.Get("name").(string))
				},
			},
		)(ctx, d, m)
	}

	return diag
}
