package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIntegrationSAML(provider string) *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "resourceIntegrationSAMLCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/saml", c.ControlPlane)
				},
				ResourceData: &SAMLSetting{IdentityProvider: provider},
				ResponseData: &AliasBasedResponse{},
			}, readSAMLIntegrationConfig,
		),
		ReadContext: ReadResource(readSAMLIntegrationConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "resourceIntegrationSAMLUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/saml/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &SAMLSetting{IdentityProvider: provider},
			}, readSAMLIntegrationConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "resourceIntegrationSAMLDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/saml/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"identity_provider": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Integration name displayed on cyral's UI.",
			},
			"single_sign_on_service_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"single_logout_service_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signing_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disable_validate_signature": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"back_channel_supported": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_want_authn_requests_signed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_want_assertions_signed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"want_assertions_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_force_authentication": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"allowed_clock_skew": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ldap_group_attribute": { //Forgerock only
				Type:     schema.TypeString,
				Optional: true,
				Default:  "cn",
			},
			"disable_post_binding_response": { //Forgerock only
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_post_binding_authn_request": { //Forgerock only
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_post_binding_logout": { //Forgerock only
				Type:     schema.TypeBool,
				Optional: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

var readSAMLIntegrationConfig = ResourceOperationConfig{
	Name:       "resourceIntegrationSAMLRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/saml/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &SAMLSetting{},
}

func (data SAMLSetting) WriteToSchema(d *schema.ResourceData) {
	d.Set("identity_provider", data.IdentityProvider)
	d.Set("display_name", data.Samlp.DisplayName)
	d.Set("single_sign_on_service_url", data.Samlp.Config.SingleSignOnServiceURL)
	d.Set("single_logout_service_url", data.Samlp.Config.SingleLogoutServiceURL)
	d.Set("signing_certificate", data.Samlp.Config.SigningCertificate)
	d.Set("disable_validate_signature", data.Samlp.Config.DisableValidateSignature)
	d.Set("back_channel_supported", data.Samlp.Config.BackChannelSupported)
	d.Set("disable_want_authn_requests_signed", data.Samlp.Config.DisableWantAuthnRequestsSigned)
	d.Set("disable_want_assertions_signed", data.Samlp.Config.DisableWantAssertionsSigned)
	d.Set("want_assertions_encrypted", data.Samlp.Config.WantAssertionsEncrypted)
	d.Set("disable_force_authentication", data.Samlp.Config.DisableForceAuthentication)
	d.Set("allowed_clock_skew", data.Samlp.Config.AllowedClockSkew)
}

func (data *SAMLSetting) ReadFromSchema(d *schema.ResourceData) {
	if data.IdentityProvider == "" {
		data.IdentityProvider = d.Get("identity_provider").(string)
		//TODO: validate if exists, since idp is required
	}
	data.Samlp.DisplayName = d.Get("display_name").(string)
	data.Samlp.Config.SingleSignOnServiceURL = d.Get("single_sign_on_service_url").(string)
	data.Samlp.Config.SingleLogoutServiceURL = d.Get("single_logout_service_url").(string)
	data.Samlp.Config.SigningCertificate = d.Get("signing_certificate").(string)
	data.Samlp.Config.DisableValidateSignature = d.Get("disable_validate_signature").(bool)
	data.Samlp.Config.BackChannelSupported = d.Get("back_channel_supported").(bool)
	data.Samlp.Config.DisableWantAuthnRequestsSigned = d.Get("disable_want_authn_requests_signed").(bool)
	data.Samlp.Config.DisableWantAssertionsSigned = d.Get("disable_want_assertions_signed").(bool)
	data.Samlp.Config.WantAssertionsEncrypted = d.Get("want_assertions_encrypted").(bool)
	data.Samlp.Config.DisableForceAuthentication = d.Get("disable_force_authentication").(bool)
	data.Samlp.Config.AllowedClockSkew = d.Get("allowed_clock_skew").(int)
}

type AliasBasedResponse struct {
	Alias string `json:"alias"`
}

func (response AliasBasedResponse) WriteToSchema(d *schema.ResourceData) {
	d.SetId(response.Alias)
}

func (response *AliasBasedResponse) ReadFromSchema(d *schema.ResourceData) {
	response.Alias = d.Id()
}
