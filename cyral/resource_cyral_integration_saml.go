package cyral

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIntegrationSAML() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "resourceIntegrationSAMLCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/saml", c.ControlPlane)
				},
				ResourceData: &SAMLIntegrationData{},
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
				ResourceData: &SAMLIntegrationData{},
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: client.ValidateIntegrationSAMLIdentityProvider(),
			},
			"ldap_group_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "cn",
			},
			"samlp": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"provider_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"first_broker_login_flow_alias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"post_broker_login_flow_alias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"store_token": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"add_read_token_role_on_create": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"trust_email": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"link_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"internal_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disable_using_jwks_url": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"sync_mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name_id_policy_format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"principal_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"signature_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"saml_xml_key_name_tranformer": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"hide_on_login_page": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"back_channel_supported": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"disable_post_binding_response": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"disable_post_binding_authn_request": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"disable_post_binding_logout": {
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
									"disable_validate_signature": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"gui_order": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"single_sign_on_service_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"single_logout_service_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"xml_sig_key_info_key_name_transformer": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"signing_certificate": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_clock_skew": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"saml_metadata_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"base_64_saml_metadata_document": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ldap_group_attribute": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
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
	ResponseData: &SAMLIntegrationData{},
}

func (data SAMLIntegrationData) WriteToSchema(d *schema.ResourceData) {
	samlSetting := data.SAMLSetting
	samlp := make([]interface{}, 0, 1)
	if samlSetting.Samlp != nil {
		samlpMap := make(map[string]interface{})
		samlpMap["alias"] = samlSetting.Samlp.Alias
		samlpMap["provider_id"] = samlSetting.Samlp.ProviderID
		samlpMap["disabled"] = samlSetting.Samlp.Disabled
		samlpMap["first_broker_login_flow_alias"] = samlSetting.Samlp.FirstBrokerLoginFlowAlias
		samlpMap["post_broker_login_flow_alias"] = samlSetting.Samlp.PostBrokerLoginFlowAlias
		samlpMap["display_name"] = samlSetting.Samlp.DisplayName
		samlpMap["store_token"] = samlSetting.Samlp.StoreToken
		samlpMap["add_read_token_role_on_create"] = samlSetting.Samlp.AddReadTokenRoleOnCreate
		samlpMap["trust_email"] = samlSetting.Samlp.TrustEmail
		samlpMap["link_only"] = samlSetting.Samlp.LinkOnly
		samlpMap["internal_id"] = samlSetting.Samlp.InternalID

		config := make([]interface{}, 0, 1)
		if samlSetting.Samlp.Config != nil {
			configMap := make(map[string]interface{})
			configMap["disable_using_jwks_url"] = samlSetting.Samlp.Config.DisableUsingJWKSUrl
			configMap["sync_mode"] = samlSetting.Samlp.Config.SyncMode
			configMap["name_id_policy_format"] = samlSetting.Samlp.Config.NameIDPolicyFormat
			configMap["principal_type"] = samlSetting.Samlp.Config.PrincipalType
			configMap["signature_type"] = samlSetting.Samlp.Config.SignatureType
			configMap["saml_xml_key_name_tranformer"] = samlSetting.Samlp.Config.SamlXmlKeyNameTranformer
			configMap["hide_on_login_page"] = samlSetting.Samlp.Config.HideOnLoginPage
			configMap["back_channel_supported"] = samlSetting.Samlp.Config.BackChannelSupported
			configMap["disable_post_binding_response"] = samlSetting.Samlp.Config.DisablePostBindingResponse
			configMap["disable_post_binding_authn_request"] = samlSetting.Samlp.Config.DisablePostBindingAuthnRequest
			configMap["disable_post_binding_logout"] = samlSetting.Samlp.Config.DisablePostBindingLogout
			configMap["disable_want_authn_requests_signed"] = samlSetting.Samlp.Config.DisableWantAuthnRequestsSigned
			configMap["disable_want_assertions_signed"] = samlSetting.Samlp.Config.DisableWantAssertionsSigned
			configMap["want_assertions_encrypted"] = samlSetting.Samlp.Config.WantAssertionsEncrypted
			configMap["disable_force_authentication"] = samlSetting.Samlp.Config.DisableForceAuthentication
			configMap["disable_validate_signature"] = samlSetting.Samlp.Config.DisableValidateSignature
			configMap["gui_order"] = samlSetting.Samlp.Config.GuiOrder
			configMap["single_sign_on_service_url"] = samlSetting.Samlp.Config.SingleSignOnServiceURL
			configMap["single_logout_service_url"] = samlSetting.Samlp.Config.SingleLogoutServiceURL
			configMap["xml_sig_key_info_key_name_transformer"] = samlSetting.Samlp.Config.XmlSigKeyInfoKeyNameTransformer
			configMap["signing_certificate"] = samlSetting.Samlp.Config.SigningCertificate
			configMap["allowed_clock_skew"] = samlSetting.Samlp.Config.AllowedClockSkew
			configMap["saml_metadata_url"] = samlSetting.Samlp.Config.SamlMetadataURL
			configMap["base_64_saml_metadata_document"] = samlSetting.Samlp.Config.Base64SamlMetadataDocument
			configMap["ldap_group_attribute"] = samlSetting.Samlp.Config.LdapGroupAttribute
			config = append(config, configMap)
		}

		samlp = append(samlp, samlpMap)
		samlpMap["config"] = config
	}
	d.Set("identity_provider", samlSetting.IdentityProvider)
	d.Set("ldap_group_attribute", samlSetting.LdapGroupAttribute)
	d.Set("samlp", samlp)
}

func (data *SAMLIntegrationData) ReadFromSchema(d *schema.ResourceData) {
	samlp := new(IdentityProviderConfig)
	for _, samlpMap := range d.Get("samlp").(*schema.Set).List() {
		samlpMap := samlpMap.(map[string]interface{})
		samlp.Alias = samlpMap["alias"].(string)
		samlp.ProviderID = samlpMap["provider_id"].(string)
		samlp.Disabled = samlpMap["disabled"].(bool)
		samlp.FirstBrokerLoginFlowAlias = samlpMap["first_broker_login_flow_alias"].(string)
		samlp.PostBrokerLoginFlowAlias = samlpMap["post_broker_login_flow_alias"].(string)
		samlp.DisplayName = samlpMap["display_name"].(string)
		samlp.StoreToken = samlpMap["store_token"].(bool)
		samlp.AddReadTokenRoleOnCreate = samlpMap["add_read_token_role_on_create"].(bool)
		samlp.TrustEmail = samlpMap["trust_email"].(bool)
		samlp.LinkOnly = samlpMap["link_only"].(bool)
		samlp.InternalID = samlpMap["internal_id"].(string)
		config := new(Config)
		for _, configMap := range samlpMap["config"].(*schema.Set).List() {
			configMap := configMap.(map[string]interface{})
			config.DisableUsingJWKSUrl = configMap["disable_using_jwks_url"].(bool)
			config.SyncMode = configMap["sync_mode"].(string)
			config.NameIDPolicyFormat = configMap["name_id_policy_format"].(string)
			config.PrincipalType = configMap["principal_type"].(string)
			config.SignatureType = configMap["signature_type"].(string)
			config.SamlXmlKeyNameTranformer = configMap["saml_xml_key_name_tranformer"].(string)
			config.HideOnLoginPage = configMap["hide_on_login_page"].(bool)
			config.BackChannelSupported = configMap["back_channel_supported"].(bool)
			config.DisablePostBindingResponse = configMap["disable_post_binding_response"].(bool)
			config.DisablePostBindingAuthnRequest = configMap["disable_post_binding_authn_request"].(bool)
			config.DisablePostBindingLogout = configMap["disable_post_binding_logout"].(bool)
			config.DisableWantAuthnRequestsSigned = configMap["disable_want_authn_requests_signed"].(bool)
			config.DisableWantAssertionsSigned = configMap["disable_want_assertions_signed"].(bool)
			config.WantAssertionsEncrypted = configMap["want_assertions_encrypted"].(bool)
			config.DisableForceAuthentication = configMap["disable_force_authentication"].(bool)
			config.DisableValidateSignature = configMap["disable_validate_signature"].(bool)
			config.GuiOrder = configMap["gui_order"].(string)
			config.SingleSignOnServiceURL = configMap["single_sign_on_service_url"].(string)
			config.SingleLogoutServiceURL = configMap["single_logout_service_url"].(string)
			config.XmlSigKeyInfoKeyNameTransformer = configMap["xml_sig_key_info_key_name_transformer"].(string)
			config.SigningCertificate = configMap["signing_certificate"].(string)
			config.AllowedClockSkew = configMap["allowed_clock_skew"].(int)
			config.SamlMetadataURL = configMap["saml_metadata_url"].(string)
			config.Base64SamlMetadataDocument = configMap["base_64_saml_metadata_document"].(string)
			config.LdapGroupAttribute = configMap["ldap_group_attribute"].(string)
		}
		samlp.Config = config
	}

	data.SAMLSetting = &SAMLSetting{
		IdentityProvider:   d.Get("identity_provider").(string),
		LdapGroupAttribute: d.Get("ldap_group_attribute").(string),
		Samlp:              samlp,
	}
}

func (resource *SAMLIntegrationData) MarshalJSON() ([]byte, error) {
	return json.Marshal(resource.SAMLSetting)
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
