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
			"sign_in_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sign_out_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"x_509_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disable_signature_validation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"back_channel_logout": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disable_authn_requests_signed": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disable_assertions_signed": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_assertions_encrypted": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disable_force_authentication": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allowed_clock_skew": {
				Type:     schema.TypeInt,
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
	d.Set("sign_in_url", data.Samlp.Config.SingleSignOnServiceURL)
	d.Set("sign_out_url", data.Samlp.Config.SingleLogoutServiceURL)
	d.Set("x_509_certificate", data.Samlp.Config.SigningCertificate)
	d.Set("disable_signature_validation", data.Samlp.Config.DisableValidateSignature)
	d.Set("back_channel_logout", data.Samlp.Config.BackChannelSupported)
	d.Set("disable_authn_requests_signed", data.Samlp.Config.DisableWantAuthnRequestsSigned)
	d.Set("disable_assertions_signed", data.Samlp.Config.DisableWantAssertionsSigned)
	d.Set("enable_assertions_encrypted", data.Samlp.Config.WantAssertionsEncrypted)
	d.Set("disable_force_authentication", data.Samlp.Config.DisableForceAuthentication)
	d.Set("allowed_clock_skew", data.Samlp.Config.AllowedClockSkew)
}

func (data *SAMLSetting) ReadFromSchema(d *schema.ResourceData) {
	if data.IdentityProvider == "" {
		data.IdentityProvider = d.Get("identity_provider").(string)
	}
	data.Samlp.DisplayName = d.Get("display_name").(string)
	data.Samlp.Config.SingleSignOnServiceURL = d.Get("sign_in_url").(string)
	data.Samlp.Config.SingleLogoutServiceURL = d.Get("sign_out_url").(string)
	data.Samlp.Config.SigningCertificate = d.Get("x_509_certificate").(string)
	data.Samlp.Config.DisableValidateSignature = d.Get("disable_signature_validation").(bool)
	data.Samlp.Config.BackChannelSupported = d.Get("back_channel_logout").(bool)
	data.Samlp.Config.DisableWantAuthnRequestsSigned = d.Get("disable_authn_requests_signed").(bool)
	data.Samlp.Config.DisableWantAssertionsSigned = d.Get("disable_assertions_signed").(bool)
	data.Samlp.Config.WantAssertionsEncrypted = d.Get("enable_assertions_encrypted").(bool)
	data.Samlp.Config.DisableForceAuthentication = d.Get("disable_force_authentication").(bool)
	data.Samlp.Config.AllowedClockSkew = d.Get("allowed_clock_skew").(uint64)
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
