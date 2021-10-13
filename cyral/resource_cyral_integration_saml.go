package cyral

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIntegrationSAML(identityProvider string) *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "resourceIntegrationSAMLCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/saml", c.ControlPlane)
				},
				ResourceData: &SAMLIntegrationData{
					SAMLSetting: &SAMLSetting{
						IdentityProvider: identityProvider,
					},
				},
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
				ResourceData: &SAMLIntegrationData{
					SAMLSetting: &SAMLSetting{
						IdentityProvider: identityProvider,
					},
				},
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
			"draft_alias": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"samlp": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_id": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "saml",
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"first_broker_login_flow_alias": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "SAML_First_Broker",
						},
						"post_broker_login_flow_alias": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  idpDefaultValues(identityProvider, "display_name"),
						},
						"store_token": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"add_read_token_role_on_create": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"trust_email": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"link_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"internal_id": {
							Type:     schema.TypeString,
							Computed: true,
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
										Default:  false,
									},
									"sync_mode": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "FORCE",
									},
									"name_id_policy_format": {
										Type:     schema.TypeString,
										Optional: true,
										Default: idpDefaultValues(identityProvider,
											"name_id_policy_format"),
									},
									"principal_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "SUBJECT",
									},
									"signature_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "RSA_SHA256",
									},
									"saml_xml_key_name_tranformer": {
										Type:     schema.TypeString,
										Optional: true,
										Default: idpDefaultValues(identityProvider,
											"saml_xml_key_name_tranformer"),
									},
									"hide_on_login_page": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"back_channel_supported": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"disable_post_binding_response": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"disable_post_binding_authn_request": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"disable_post_binding_logout": {
										Type:     schema.TypeBool,
										Optional: true,
										Default: idpDefaultValues(identityProvider,
											"disable_post_binding_logout"),
									},
									"want_assertions_encrypted": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"disable_force_authentication": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"gui_order": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"single_sign_on_service_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"single_logout_service_url": {
										Type:     schema.TypeString,
										Optional: true,
										Default: idpDefaultValues(identityProvider,
											"single_logout_service_url"),
									},
									"xml_sig_key_info_key_name_transformer": {
										Type:     schema.TypeString,
										Optional: true,
										Default: idpDefaultValues(identityProvider,
											"xml_sig_key_info_key_name_transformer"),
									},
									"signing_certificate": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"allowed_clock_skew": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"saml_metadata_url": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"base_64_saml_metadata_document": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"ldap_group_attribute": {
										Type:     schema.TypeString,
										Optional: true,
										Default: idpDefaultValues(identityProvider,
											"ldap_group_attribute"),
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

var (
	defaultValuesMap = map[string]interface{}{
		"display_name":                          "",
		"disable_post_binding_logout":           false,
		"name_id_policy_format":                 "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
		"saml_xml_key_name_tranformer":          "KEY_ID",
		"single_logout_service_url":             "",
		"xml_sig_key_info_key_name_transformer": "KEY_ID",
		"ldap_group_attribute":                  "",
	}
	adfsDefaultValuesMap = map[string]interface{}{
		"display_name":                          "Active Directory",
		"disable_post_binding_logout":           true,
		"name_id_policy_format":                 "urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
		"saml_xml_key_name_tranformer":          "CERT_SUBJECT",
		"single_logout_service_url":             "",
		"xml_sig_key_info_key_name_transformer": "CERT_SUBJECT",
	}
	aadDefaultValuesMap = map[string]interface{}{
		"display_name": "Azure Active Directory",
	}
	forgerockDefaultValuesMap = map[string]interface{}{
		"display_name":         "Forgerock",
		"ldap_group_attribute": "cn",
	}
	gsuiteDefaultValuesMap = map[string]interface{}{
		"display_name": "GSuite",
	}
	oktaDefaultValuesMap = map[string]interface{}{
		"display_name": "Okta",
	}
	pingoneDefaultValuesMap = map[string]interface{}{
		"display_name": "Pingone",
	}
)

func idpDefaultValues(identityProvider, fieldName string) interface{} {
	switch identityProvider {
	case "adfs-2016":
		if value, ok := adfsDefaultValuesMap[fieldName]; ok {
			return value
		}
	case "aad":
		if value, ok := aadDefaultValuesMap[fieldName]; ok {
			return value
		}
	case "forgerock":
		if value, ok := forgerockDefaultValuesMap[fieldName]; ok {
			return value
		}
	case "gsuite":
		if value, ok := gsuiteDefaultValuesMap[fieldName]; ok {
			return value
		}
	case "okta":
		if value, ok := oktaDefaultValuesMap[fieldName]; ok {
			return value
		}
	case "pingone":
		if value, ok := pingoneDefaultValuesMap[fieldName]; ok {
			return value
		}
	}
	return defaultValuesMap[fieldName]
}

func (data SAMLIntegrationData) WriteToSchema(d *schema.ResourceData) {
	samlSetting := data.SAMLSetting
	samlp := make([]interface{}, 0, 1)
	if samlSetting.Samlp != nil {
		samlpMap := make(map[string]interface{})
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
			configMap["want_assertions_encrypted"] = samlSetting.Samlp.Config.WantAssertionsEncrypted
			configMap["disable_force_authentication"] = samlSetting.Samlp.Config.DisableForceAuthentication
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
		samlpMap["config"] = config
		samlp = append(samlp, samlpMap)
	}
	d.Set("samlp", samlp)
}

func (data *SAMLIntegrationData) ReadFromSchema(d *schema.ResourceData) {
	samlp := new(IdentityProviderConfig)
	for _, samlpMap := range d.Get("samlp").(*schema.Set).List() {
		samlpMap := samlpMap.(map[string]interface{})
		samlp.Alias = d.Get("draft_alias").(string)
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
			config.WantAssertionsEncrypted = configMap["want_assertions_encrypted"].(bool)
			config.DisableForceAuthentication = configMap["disable_force_authentication"].(bool)
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

	data.SAMLSetting.LdapGroupAttribute = samlp.Config.LdapGroupAttribute
	data.SAMLSetting.Samlp = samlp
}

func (resource SAMLIntegrationData) MarshalJSON() ([]byte, error) {
	return json.Marshal(resource.SAMLSetting)
}

type AliasBasedResponse struct {
	Alias string `json:"alias"`
}

func (response AliasBasedResponse) WriteToSchema(d *schema.ResourceData) {
	d.SetId(response.Alias)
}

func (response *AliasBasedResponse) ReadFromSchema(d *schema.ResourceData) {}
