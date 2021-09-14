package cyral

import (
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ClockSkew int

func (c *ClockSkew) UnmarshalJSON(data []byte) error {
	stew := ""
	err := json.Unmarshal(data, &stew)
	if err != nil {
		return err
	}
	var d int
	d, err = strconv.Atoi(stew)
	if err != nil {
		return err
	}
	*c = ClockSkew(d)
	return nil
}

type SAMLConfig struct {
	AllowedClockSkew                ClockSkew `json:"allowedClockSkew,omitempty"`
	BackChannelSupported            bool      `json:"backChannelSupported,omitempty"`
	DisableForceAuthentication      bool      `json:"disableForceAuthentication,omitempty"`
	DisablePostBindingAuthnRequest  bool      `json:"disablePostBindingAuthnRequest,omitempty"`
	DisablePostBindingResponse      bool      `json:"disablePostBindingResponse,omitempty"`
	DisablePostBindingLogout        bool      `json:"disablePostBindingLogout,omitempty"`
	DisableUsingJWKSURL             bool      `json:"disableUsingJWKSUrl,omitempty"`
	DisableValidateSignature        bool      `json:"disableValidateSignature,omitempty"`
	DisableWantAssertionSigned      bool      `json:"disableWantAssertionSigned,omitempty"`
	DisableWantAuthnRequestsSigned  bool      `json:"disableWantAuthnRequestsSigned,omitempty"`
	GUIOrder                        string    `json:"guiOrder,omitempty"`
	HideOnLoginPage                 bool      `json:"hideOnLoginPage,omitempty"`
	NameIDPolicyFormat              string    `json:"nameIDPolicyFormat,omitempty"`
	PrincipalType                   string    `json:"principalType,omitempty"`
	SAMLXMLKeyNameTransformer       string    `json:"samlXmlKeyNameTransformer,omitempty"`
	SignatureType                   string    `json:"signatureType,omitempty"`
	SigningCertificate              string    `json:"signingCertificate,omitempty"`
	SingleLogoutServiceURL          string    `json:"singleLogoutServiceURL,omitempty"`
	SingleSignOnServiceURL          string    `json:"singleSignOnServiceURL,omitempty"`
	SyncMode                        string    `json:"syncMode,omitempty"`
	WantAssertionsEncrypted         bool      `json:"wantAssertionsEncrypted,omitempty"`
	XMLSigKeyInfoKeyNameTransformer string    `json:"xmlSigKeyInfoKeyNameTransformer,omitempty"`
	SAMLMetadataURL                 string    `json:"samlMetadataURL"`
	Base64SamlMetadataDocument      string    `json:"base64SamlMetadataDocument"`
	LDAPGroupAttribute              string    `json:"ldapGroupAttribute"`
}

func (data *SAMLConfig) WriteToSchema(rd *schema.ResourceData) {
	rd.Set("allowed_clock_skew", data.AllowedClockSkew)
	rd.Set("back_channel_supported", data.BackChannelSupported)
	rd.Set("disable_force_authentication", data.DisableForceAuthentication)
	rd.Set("disable_post_binding_authn_request", data.DisablePostBindingAuthnRequest)
	rd.Set("disable_post_binding_response", data.DisablePostBindingResponse)
	rd.Set("disable_post_binding_logout", data.DisablePostBindingLogout)
	rd.Set("disable_using_jwks_url", data.DisableUsingJWKSURL)
	rd.Set("disable_validate_signature", data.DisableValidateSignature)
	rd.Set("disable_want_assertions_signed", data.DisableWantAssertionSigned)
	rd.Set("disable_want_authn_requests_signed", data.DisableWantAuthnRequestsSigned)
	rd.Set("gui_order", data.GUIOrder)
	rd.Set("hide_on_login_page", data.HideOnLoginPage)
	rd.Set("name_id_policy_format", data.NameIDPolicyFormat)
	rd.Set("principal_type", data.PrincipalType)
	rd.Set("saml_xml_key_name_transformer", data.SAMLXMLKeyNameTransformer)
	rd.Set("signature_type", data.SignatureType)
	rd.Set("signing_certificate", data.SigningCertificate)
	rd.Set("single_logout_service_url", data.SingleLogoutServiceURL)
	rd.Set("single_sign_on_service_url", data.SingleSignOnServiceURL)
	rd.Set("sync_mode", data.SyncMode)
	rd.Set("want_assertions_encrypted", data.WantAssertionsEncrypted)
	rd.Set("xml_sig_key_info_key_name_transformer", data.XMLSigKeyInfoKeyNameTransformer)
	rd.Set("saml_metadata_url", data.SAMLMetadataURL)
	rd.Set("base_64_saml_metadata_document", data.Base64SamlMetadataDocument)
	rd.Set("ldap_group_attribute", data.LDAPGroupAttribute)
}
func (data *SAMLConfig) ReadFromSchema(rd *schema.ResourceData) {
	data.AllowedClockSkew = ClockSkew(rd.Get("allowed_clock_skew").(int))
	data.BackChannelSupported = rd.Get("back_channel_supported").(bool)
	data.DisableForceAuthentication = rd.Get("disable_force_authentication").(bool)
	data.DisablePostBindingAuthnRequest = rd.Get("disable_post_binding_authn_request").(bool)
	data.DisablePostBindingResponse = rd.Get("disable_post_binding_response").(bool)
	data.DisablePostBindingLogout = rd.Get("disable_post_binding_logout").(bool)
	data.DisableUsingJWKSURL = rd.Get("disable_using_jwks_url").(bool)
	data.DisableValidateSignature = rd.Get("disable_validate_signature").(bool)
	data.DisableWantAssertionSigned = rd.Get("disable_want_assertions_signed").(bool)
	data.DisableWantAuthnRequestsSigned = rd.Get("disable_want_authn_requests_signed").(bool)
	data.GUIOrder = rd.Get("gui_order").(string)
	data.HideOnLoginPage = rd.Get("hide_on_login_page").(bool)
	data.NameIDPolicyFormat = rd.Get("name_id_policy_format").(string)
	data.PrincipalType = rd.Get("principal_type").(string)
	data.SAMLXMLKeyNameTransformer = rd.Get("saml_xml_key_name_transformer").(string)
	data.SignatureType = rd.Get("signature_type").(string)
	data.SigningCertificate = rd.Get("signing_certificate").(string)
	data.SingleLogoutServiceURL = rd.Get("single_logout_service_url").(string)
	data.SingleSignOnServiceURL = rd.Get("single_sign_on_service_url").(string)
	data.SyncMode = rd.Get("sync_mode").(string)
	data.WantAssertionsEncrypted = rd.Get("want_assertions_encrypted").(bool)
	data.XMLSigKeyInfoKeyNameTransformer = rd.Get("xml_sig_key_info_key_name_transformer").(string)
	data.SAMLMetadataURL = rd.Get("saml_metadata_url").(string)
	data.Base64SamlMetadataDocument = rd.Get("base_64_saml_metadata_document").(string)
	data.LDAPGroupAttribute = rd.Get("ldap_group_attribute").(string)
}

type ParseRequest struct {
	Metadata string `json:"metadata,omitempty"`
}

func (ParseRequest) ReadFromSchema(rd *schema.ResourceData) {}
func (ParseRequest) WriteToSchema(rd *schema.ResourceData)  {}

type SAMLPayload struct {
	Config                    SAMLConfig `json:"config,omitempty"`
	AddReadTokenRoleOnCreate  bool       `json:"addReadTokenRoleOnCreate,omitempty"`
	Alias                     string     `json:"alias,omitempty"`
	Disabled                  bool       `json:"disabled,omitempty"`
	DisplayName               string     `json:"displayName,omitempty"`
	FirstBrokerLoginFlowAlias string     `json:"firstBrokerLoginFlowAlias,omitempty"`
	InternalID                string     `json:"internalID,omitempty"`
	LinkOnly                  bool       `json:"linkOnly,omitempty"`
	ProviderID                string     `json:"providerID,omitempty"`
	StoreToken                bool       `json:"storeToken,omitempty"`
	TrustEmail                bool       `json:"trustEmail,omitempty"`
}

type Config struct {
	// By default, we use the jwks URL for all SAML connections.
	DisableUsingJWKSUrl bool `json:"disableUsingJWKSUrl,omitempty"`
	// Defaults to "FORCE" if unset.
	SyncMode string `json:"syncMode,omitempty"`
	// Defaults to "urn:oasis:names:tc:SAML:1.1:nameidFormat:unspecified" if unset.
	NameIDPolicyFormat string `json:"nameIDPolicyFormat,omitempty"`
	// Defaults to "SUBJECT" if unset.
	PrincipalType string `json:"principalType,omitempty"`
	// Defaults to "RSA_SHA256" if unset.
	SignatureType string `json:"signatureType,omitempty"`
	// Defaults to "KEY_ID" if unset.
	SamlXmlKeyNameTranformer string `json:"samlXmlKeyNameTranformer,omitempty"`
	// Defaults to false.
	HideOnLoginPage bool `json:"hideOnLoginPage,omitempty"`
	// Defaults to false.
	BackChannelSupported bool `json:"backChannelSupported,omitempty"`
	// Indicates whether to respond to requests using HTTP-POST binding.
	// If true, HTTP-REDIRECT binding will be used.
	// Defaults to false.
	DisablePostBindingResponse bool `json:"disablePostBindingResponse,omitempty"`
	// Indicates whether the AuthnRequest must be sent using HTTP-POST binding.
	// If true, HTTP-REDIRECT binding will be used.
	// Defaults to false.
	DisablePostBindingAuthnRequest bool `json:"disablePostBindingAuthnRequest,omitempty"`
	// Indicates whether to respond to requests using HTTP-POST binding.
	// If true, HTTP-REDIRECT binding will be used.
	// Defaults to false.
	DisablePostBindingLogout bool `json:"disablePostBindingLogout,omitempty"`
	// Indicates whether the identity provider expects a signed AuthnRequest.
	// Defaults to false.
	DisableWantAuthnRequestsSigned bool `json:"disableWantAuthnRequestsSigned,omitempty"`
	// Indicates whether this service provider expects a signed Assertion.
	// Defaults to false.
	DisableWantAssertionsSigned bool `json:"disableWantAssertionsSigned,omitempty"`
	// Indicates whether this service provider expects an encrypted Assertion.
	// Defaults to false.
	WantAssertionsEncrypted bool `json:"wantAssertionsEncrypted,omitempty"`
	// Indicates whether the identity provider must authenticate the presenter directly
	// rather than rely on a previous security context.
	// Defaults to false.
	DisableForceAuthentication bool `json:"disableForceAuthentication,omitempty"`
	// Enable/Disable signature validation of SAML responses. Highly recommended for minimum security.
	// Defaults to false.
	DisableValidateSignature bool   `json:"disableValidateSignature,omitempty"`
	GuiOrder                 string `json:"guiOrder,omitempty"`
	// The Url that must be used to send authentication requests (SAML AuthnRequest).
	SingleSignOnServiceURL string `json:"singleSignOnServiceURL,omitempty"`
	// The Url that must be used to send logout requests.
	SingleLogoutServiceURL string `json:"singleLogoutServiceURL,omitempty"`
	// Defaults to "KEY_ID" if unset.
	XmlSigKeyInfoKeyNameTransformer string `json:"xmlSigKeInfoKeyNameTransformer,omitempty"`
	// The signing certificate used to validate signatures. Required if signature validation is enabled.
	SigningCertificate string `json:"signingCertificate,omitempty"`
	// Clock skew in seconds that is tolerated when validating identity provider tokens.
	// Default value is zero.
	AllowedClockSkew uint64 `json:"allowedClockSkew,omitempty"`
	// The SAML metadata URL that was used to create the integration, if any.
	// This field is added on Cyral's side. The Cyral Terraform provider needs
	// it to tell if a new SAML metadata being used is the same as before, or
	// if it's different.
	SamlMetadataURL string `json:"samlMetadataURL,omitempty"`
	// This is the full SAML metadata document that was used to import the SAML config,
	// base64 encoded.
	Base64SamlMetadataDocument string `json:"base64SamlMetadataDocument,omitempty"`
	// The type of LDAP RDN that identifies the name of a group within a DN. For example, if an LDAP DN sent
	// in a SAML assertion is "cn=Everyone,ou=groups,dc=openam,dc=forgerock,dc=org and the LDAP Group RDN
	// Type is "cn," Cyral will interpret Everyone as the group name.
	LdapGroupAttribute string `json:"ldapGroupAttribute,omitempty"`
}

type IdentityProviderConfig struct {
	// The alias uniquely identifies an identity provider and it is also used to build the redirect uri. Must be uri friendly.
	Alias string `json:"alias,omitempty"`
	// This is the provider ID of "saml".
	ProviderID string `json:"providerID,omitempty"`
	// Disable maps to Keycloak's "enabled" field.
	// Defaults to false.
	Disabled bool `json:"disabled,omitempty"`
	// Alias of authentication flow, which is triggered after first login with this identity provider.
	//Term 'First Login' means that no Keycloak account is currently linked to the authenticated identity provider account.
	FirstBrokerLoginFlowAlias string `json:"firstBrokerLoginFlowAlias,omitempty"`
	// Alias of authentication flow, which is triggered after each login with this identity provider.
	// Useful if you want additional verification of each user authenticated with this identity provider (for example OTP).
	// Leave this empty if you need no any additional authenticators to be triggered after login with this identity provider.
	// Also note that authenticator implementations must assume that user is already set in ClientSession as identity provider already set it.
	PostBrokerLoginFlowAlias string `json:"postBrokerLoginFlowAlias,omitempty"`
	// Friendly name for Identity Providers.
	// SHOW IN UI
	DisplayName string `json:"displayName,omitempty"`
	// Enable/disable if tokens must be stored after authenticating users.
	// Default to False -- no tokens in SAML, just assertions
	StoreToken               bool `json:"storeToken,omitempty"`
	AddReadTokenRoleOnCreate bool `json:"addReadTokenRoleOnCreate,omitempty"`
	// If the identity provider supplies an email address this email address will be trusted. If the realm required email validation, users that log in from this IDP will not have to go through the email verification process.
	TrustEmail bool `json:"trustEmail,omitempty"`
	// If true, users cannot log in through this provider. They can only link to this provider.
	// This is useful if you don't want to allow login from the provider, but want to integrate with a provider
	LinkOnly bool `json:"linkOnly,omitempty"`
	// autoGenerated
	InternalID string  `json:"internalID,omitempty"`
	Config     *Config `json:"config,omitempty"`
}

type SAMLSetting struct {
	// IdentityProvider corresponds to the IdP provider type (currently supported: generic or forgerock)
	IdentityProvider string `json:"identityProvider,omitempty"`
	// LdapGroupAttribute corresponds to the LDAP Groups Search Attribute in ForgeRock's Identity Store
	LdapGroupAttribute string `json:"ldapGroupAttribute,omitempty"`
	// samlp is for providing SAML configuration directly rather
	// than through a metadata document.
	Samlp IdentityProviderConfig `json:"samlp,omitempty"`
	// This is the full SAML metadata URL we should use to import the SAML config.
	// It will only be populated if a SAML metadata URL was used to create the
	// integration.
	SamlMetadataURL string `json:"samlMetadataURL,omitempty"`
	// This is the full SAML metadata document that we should use to import the SAML config,
	// base64 encoded.
}
