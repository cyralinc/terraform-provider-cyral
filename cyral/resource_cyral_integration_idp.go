package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIntegrationIdP(identityProvider string, deprecationMessage string) *schema.Resource {
	return &schema.Resource{
		Description:        fmt.Sprintf("%v", idpDefaultValues(identityProvider, "resource_description")),
		CreateContext:      resourceIntegrationIdPCreate(identityProvider),
		ReadContext:        resourceIntegrationIdPRead,
		UpdateContext:      resourceIntegrationIdPUpdate(identityProvider),
		DeleteContext:      resourceIntegrationIdPDelete,
		DeprecationMessage: deprecationMessage,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource, which corresponds to the IdP Integration `alias`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"draft_alias": {
				Description: "An `alias` that uniquely identifies a IdP Integration draft. If set, will delete any " +
					"correspondent draft and create a new IdP Integration with the same `alias`. Defaults to `\"\"`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"samlp": {
				Description: "It contains the top-level configuration for an identity provider.",
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_id": {
							Description: "This is the provider ID of `saml`. Defaults to `saml`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "saml",
						},
						"disabled": {
							Description: "Disable maps to Keycloak's `enabled` field. Defaults to `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"first_broker_login_flow_alias": {
							Description: "Alias of authentication flow, which is triggered after `First Login` with this identity provider. Term `First Login` means that no Keycloak account is currently linked to the authenticated identity provider account. Defaults to `SAML_First_Broker`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "SAML_First_Broker",
						},
						"post_broker_login_flow_alias": {
							Description: "Alias of authentication flow, which is triggered after each login with this identity provider. Useful if you want additional verification of each user authenticated with this identity provider (for example OTP). Leave this empty if you need no any additional authenticators to be triggered after login with this identity provider. Defaults to `\"\"`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
						},
						"display_name": {
							Description: "Name of the IdP Integration displayed in the control plane. Defaults to `" +
								fmt.Sprintf("%v", idpDefaultValues(identityProvider, "display_name")) + "`",
							Type:     schema.TypeString,
							Optional: true,
							Default:  idpDefaultValues(identityProvider, "display_name"),
						},
						"store_token": {
							Description: "Enable if tokens must be stored after authenticating users. Defaults to `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"add_read_token_role_on_create": {
							Description: "Adds read token role on creation. Defaults to `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"trust_email": {
							Description: "If the identity provider supplies an email address this email address will be trusted. If the realm required email validation, users that log in from this identity provider will not have to go through the email verification process. Defaults to `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"link_only": {
							Description: "If true, users cannot log in through this identity provider. They can only link to this identity provider. This is useful if you don't want to allow login from the identity provider, but want to integrate with an identity provider. Defaults to `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"internal_id": {
							Description: "An ID that is auto-generated internally for this IdP Integration.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"config": {
							Description: "SAML configuration for this IdP Integration.",
							Type:        schema.TypeSet,
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disable_using_jwks_url": {
										Description: "By default, the jwks URL is used for all SAML connections. Defaults to `false`.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"sync_mode": {
										Description: "Defaults to `FORCE` if unset.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "FORCE",
									},
									"name_id_policy_format": {
										Description: "Defaults to `urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified` if unset.",
										Type:        schema.TypeString,
										Optional:    true,
										Default: idpDefaultValues(identityProvider,
											"name_id_policy_format"),
									},
									"principal_type": {
										Description: "Defaults to `SUBJECT` if unset.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "SUBJECT",
									},
									"signature_type": {
										Description: "Defaults to `RSA_SHA256` if unset.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "RSA_SHA256",
									},
									"saml_xml_key_name_tranformer": {
										Description: "Defaults to `KEY_ID` if unset.",
										Type:        schema.TypeString,
										Optional:    true,
										Default: idpDefaultValues(identityProvider,
											"saml_xml_key_name_tranformer"),
									},
									"hide_on_login_page": {
										Description: "Defaults to `false` if unset.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"back_channel_supported": {
										Description: "Defaults to `false` if unset.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"disable_post_binding_response": {
										Description: "Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used. Defaults to `false`.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"disable_post_binding_authn_request": {
										Description: "Indicates whether the AuthnRequest must be sent using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used. Defaults to `false`.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"disable_post_binding_logout": {
										Description: "Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used. Defaults to `false`.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default: idpDefaultValues(identityProvider,
											"disable_post_binding_logout"),
									},
									"want_assertions_encrypted": {
										Description: "Indicates whether the service provider expects an encrypted Assertion. Defaults to `false`.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"disable_force_authentication": {
										Description: "Indicates whether the identity provider must authenticate the presenter directly rather than rely on a previous security context. Defaults to `false`",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"gui_order": {
										Description: "GUI order. Defaults to `\"\"`.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
									},
									"single_sign_on_service_url": {
										Description: "URL that must be used to send authentication requests (SAML AuthnRequest).",
										Type:        schema.TypeString,
										Required:    true,
									},
									"single_logout_service_url": {
										Description: "URL that must be used to send logout requests. Defaults to `\"\"`.",
										Type:        schema.TypeString,
										Optional:    true,
										Default: idpDefaultValues(identityProvider,
											"single_logout_service_url"),
									},
									"xml_sig_key_info_key_name_transformer": {
										Description: "Defaults to `KEY_ID` if unset.",
										Type:        schema.TypeString,
										Optional:    true,
										Default: idpDefaultValues(identityProvider,
											"xml_sig_key_info_key_name_transformer"),
									},
									"signing_certificate": {
										Description: "Signing certificate used to validate signatures. Required if signature validation is enabled. Defaults to `\"\"`.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
									},
									"allowed_clock_skew": {
										Description: "Clock skew in seconds that is tolerated when validating identity provider tokens. Defaults to `0`.",
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
									},
									"saml_metadata_url": {
										Description: "This is the full SAML metadata URL that was used to import the SAML configuration. Defaults to `\"\"`.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
									},
									"base_64_saml_metadata_document": {
										Description: "Full SAML metadata document that was used to import the SAML configuration, Base64 encoded. Defaults to `\"\"`.",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
									},
									"ldap_group_attribute": {
										Description: "Type of `LDAP Group RDN` that identifies the name of a group within a DN. For example, if an LDAP DN sent in a SAML assertion is `cn=Everyone`, `ou=groups`, `dc=openam`, `dc=forgerock`, `dc=org` and the `LDAP Group RDN` Type is `cn` Cyral will interpret `Everyone` as the group name.",
										Type:        schema.TypeString,
										Optional:    true,
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

func resourceIntegrationIdPCreate(identityProvider string) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		diag := CreateResource(
			ResourceOperationConfig{
				Name:       "resourceIntegrationIdPCreate - Integration",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/saml", c.ControlPlane)
				},
				NewResourceData: func() ResourceData {
					return &SAMLIntegrationData{
						SAMLSetting: &SAMLSetting{
							IdentityProvider: identityProvider,
						},
					}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData { return &AliasBasedResponse{} },
			}, readIntegrationIdPConfig,
		)(ctx, d, m)

		if !diag.HasError() {
			diag = CreateResource(
				ResourceOperationConfig{
					Name:       "resourceIntegrationIdPCreate - IdentityProvider",
					HttpMethod: http.MethodPost,
					CreateURL: func(d *schema.ResourceData, c *client.Client) string {
						return fmt.Sprintf("https://%s/v1/conf/identityProviders/%s", c.ControlPlane, d.Id())
					},
					NewResourceData: func() ResourceData { return &IdentityProviderData{} },
					NewResponseData: func(_ *schema.ResourceData) ResponseData { return &IdentityProviderData{} },
				}, readIdentityProviderConfig,
			)(ctx, d, m)

			if diag.HasError() {
				// Clean Up Integration IdP
				DeleteResource(deleteIntegrationIdPConfig)(ctx, d, m)
			}
		}

		return diag
	}
}

func resourceIntegrationIdPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := ReadResource(readIntegrationIdPConfig)(ctx, d, m)

	if !diag.HasError() {
		diag = ReadResource(readIdentityProviderConfig)(ctx, d, m)
	}

	return diag
}

func resourceIntegrationIdPUpdate(identityProvider string) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		diag := UpdateResource(
			ResourceOperationConfig{
				Name:       "resourceIntegrationIdPUpdate - Integration",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/saml/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() ResourceData {
					return &SAMLIntegrationData{
						SAMLSetting: &SAMLSetting{
							IdentityProvider: identityProvider,
						},
					}
				},
			}, readIntegrationIdPConfig,
		)(ctx, d, m)

		return diag
	}
}

func resourceIntegrationIdPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := DeleteResource(deleteIntegrationIdPConfig)(ctx, d, m)

	if !diag.HasError() {
		diag = DeleteResource(
			ResourceOperationConfig{
				Name:       "resourceIntegrationIdPDelete - IdentityProvider",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/conf/identityProviders/%s", c.ControlPlane, d.Id())
				},
			},
		)(ctx, d, m)
	}

	return diag
}

var readIntegrationIdPConfig = ResourceOperationConfig{
	Name:       "resourceIntegrationIdPRead - Integration",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/saml/%s", c.ControlPlane, d.Id())
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData { return &SAMLIntegrationData{} },
}

var readIdentityProviderConfig = ResourceOperationConfig{
	Name:       "resourceIntegrationIdPRead - IdentityProvider",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/conf/identityProviders/%s", c.ControlPlane, d.Id())
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData { return &IdentityProviderData{} },
}

var deleteIntegrationIdPConfig = ResourceOperationConfig{
	Name:       "resourceIntegrationIdPDelete - Integration",
	HttpMethod: http.MethodDelete,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/saml/%s", c.ControlPlane, d.Id())
	},
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
		"name_id_policy_format":                 "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
		"saml_xml_key_name_tranformer":          "CERT_SUBJECT",
		"single_logout_service_url":             "",
		"xml_sig_key_info_key_name_transformer": "CERT_SUBJECT",
		"resource_description":                  "Manages [integration with Active Directory Federation Services](https://cyral.com/docs/sso/sso-adfs#add-your-adfs-as-an-idp-in-cyral) identity provider to allow single-sign on to Cyral.",
	}
	aadDefaultValuesMap = map[string]interface{}{
		"display_name":         "Azure Active Directory",
		"resource_description": "Manages [integration with Azure Active Directory](https://cyral.com/docs/sso/sso-azure-ad#add-azure-ad-integration-to-cyral) identity provider to allow single-sign on to Cyral.",
	}
	forgerockDefaultValuesMap = map[string]interface{}{
		"display_name":         "Forgerock",
		"ldap_group_attribute": "cn",
		"resource_description": "Manages [integration with Forgerock](https://cyral.com/docs/sso/sso-forgerock#add-forgerock-idp-to-cyral) identity provider to allow single-sign on to Cyral.",
	}
	gsuiteDefaultValuesMap = map[string]interface{}{
		"display_name":         "GSuite",
		"resource_description": "Manages integration with GSuite identity provider to allow single-sign on to Cyral.",
	}
	oktaDefaultValuesMap = map[string]interface{}{
		"display_name":         "Okta",
		"resource_description": "Manages [integration with Okta](https://cyral.com/docs/sso/sso-okta#in-cyral-management-console-create-okta-integration) identity provider to allow single-sign on to Cyral.",
	}
	pingoneDefaultValuesMap = map[string]interface{}{
		"display_name":         "PingOne",
		"resource_description": "Manages integration with PingOne identity provider to allow single-sign on to Cyral.",
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

func (data SAMLIntegrationData) WriteToSchema(d *schema.ResourceData, c *client.Client) error {
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

	return nil
}

func (data *SAMLIntegrationData) ReadFromSchema(d *schema.ResourceData, c *client.Client) error {
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

	return nil
}

func (resource SAMLIntegrationData) MarshalJSON() ([]byte, error) {
	return json.Marshal(resource.SAMLSetting)
}

type AliasBasedResponse struct {
	Alias string `json:"alias"`
}

func (response AliasBasedResponse) WriteToSchema(d *schema.ResourceData, c *client.Client) error {
	d.SetId(response.Alias)
	return nil
}

type KeycloakProvider struct{}

type IdentityProviderData struct {
	Keycloak KeycloakProvider `json:"keycloakProvider"`
}

func (data IdentityProviderData) WriteToSchema(d *schema.ResourceData, c *client.Client) error {
	return nil
}

func (data *IdentityProviderData) ReadFromSchema(d *schema.ResourceData, c *client.Client) error {
	return nil
}
