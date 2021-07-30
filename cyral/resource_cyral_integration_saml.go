package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (data SAMLIntegration) WriteToSchema(d *schema.ResourceData) {
	d.Set("identity_provider", data.SAMLSetting.IdentityProvider)
	d.Set("base_64_metadata_document", data.SAMLSetting)
	d.Set("sign_in_url", data.SAMLSetting.Samlp.Config.SingleSignOnServiceURL)
	d.Set("sign_out_url", data.SAMLSetting.Samlp.Config.SingleLogoutServiceURL)
	d.Set("disable_signature_validation", data.SAMLSetting.Samlp.Config.DisableValidateSignature)
	d.Set("disable_assertion_signature", data.SAMLSetting.Samlp.Config)
	d.Set("enable_assertion_encryption", data.SAMLSetting.Samlp.Config.WantAssertionsEncrypted)
	d.Set("disable_force_authentication", data.SAMLSetting.Samlp.Config.DisableForceAuthentication)
}

func (data *SAMLIntegration) ReadFromSchema(d *schema.ResourceData) {
	if data.OverrideProvider == "" {
		data.SAMLSetting.IdentityProvider = d.Get("identity_provider").(string)
	} else {

		data.SAMLSetting.IdentityProvider = data.OverrideProvider
		data.SAMLSetting.IdentityProvider = data.OverrideProvider
	}
	data.SAMLSetting.Samlp.Config.SingleSignOnServiceURL = d.Get("sign_in_url").(string)
	data.SAMLSetting.Samlp.Config.SingleLogoutServiceURL = d.Get("sign_out_url").(string)
	data.SAMLSetting.Samlp.Config.DisableValidateSignature = d.Get("disable_signature_validation").(bool)
	data.SAMLSetting.Samlp.Config.DisableWantAssertionsSigned = d.Get("disable_assertion_signature").(bool)
	data.SAMLSetting.Samlp.Config.WantAssertionsEncrypted = d.Get("enable_assertion_encryption").(bool)
	data.SAMLSetting.Samlp.Config.DisableForceAuthentication = d.Get("disable_force_authentication").(bool)
	data.SAMLSetting.Samlp.Config.XmlSigKeyInfoKeyNameTransformer = "KEY_ID"
	data.SAMLSetting.Samlp.Config.SamlXmlKeyNameTranformer = "KEY_ID"
	data.SAMLSetting.Samlp.Config.NameIDPolicyFormat = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	data.SAMLSetting.Samlp.FirstBrokerLoginFlowAlias = "SAML_First_Broker"
	data.SAMLSetting.Samlp.ProviderID = "saml"
	data.SAMLSetting.Samlp.Config.SyncMode = "FORCE"
	data.SAMLSetting.Samlp.Config.PrincipalType = "SUBJECT"
	data.SAMLSetting.Samlp.Config.SignatureType = "RSA_SHA256"
}

var ReadSAMLIntegrationConfig = ResourceOperationConfig{
	Name:       "SAMLIntegrationResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/saml/%s", c.ControlPlane, d.Id())
	},
	ResponseData: &SAMLIntegration{},
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

func resourceIntegrationSAMLIntegration(provider string) *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "SAMLIntegrationResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {

					return fmt.Sprintf("https://%s/v1/integrations/saml", c.ControlPlane)
				},
				ResourceData: &SAMLIntegration{OverrideProvider: provider},
				ResponseData: &AliasBasedResponse{},
			}, ReadSAMLIntegrationConfig,
		),
		ReadContext: ReadResource(ReadSAMLIntegrationConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "SAMLIntegrationResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/saml/%s", c.ControlPlane, d.Id())
				},
				ResourceData: &SAMLIntegration{OverrideProvider: provider},
			}, ReadSAMLIntegrationConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "SAMLIntegrationResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/saml/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: genericMetadataSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

var genericMetadataSchema = map[string]*schema.Schema{
	"identity_provider": {
		Type:     schema.TypeString,
		Required: true,
	},
	"display_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "name of the integration on cyral's UI",
	},
	"sign_in_url": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"sign_out_url": {
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
	"disable_authn_requests_signature": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"disable_assertion_signature": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"enable_assertion_encryption": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"disable_force_authentication": {
		Type:     schema.TypeString,
		Optional: true,
	},
}
