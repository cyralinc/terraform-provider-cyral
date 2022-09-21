package cyral

type IdPIntegrations struct {
	Connections *Connections `json:"connections,omitempty"`
}

type Connections struct {
	Connections []*Connection `json:"connections,omitempty"`
}

type Connection struct {
	DisplayName            string `json:"displayName,omitempty"`
	Alias                  string `json:"alias,omitempty"`
	SingleSignOnServiceURL string `json:"singleSignOnServiceURL,omitempty"`
	Enabled                bool   `json:"enabled,omitempty"`
}

type ParseSAMLMetadataRequest struct {
	// This is the full SAML metadata URL we should use to parse to a SAML config.
	SamlMetadataURL string `json:"samlMetadataURL,omitempty"`
	// This is the full SAML metadata document that we should use to parse to a SAML config,
	// base64 encoded.
	Base64SamlMetadataDocument string `json:"base64SamlMetadataDocument,omitempty"`
}

type SAMLConfiguration struct {
	Config *Config `json:"config,omitempty"`
}

type Config struct {
	// By default, we use the jwks URL for all SAML connections.
	DisableUsingJWKSUrl bool `json:"disableUsingJWKSUrl,omitempty"`
	// Defaults to "FORCE" if unset.
	SyncMode string `json:"syncMode,omitempty"`
	// Defaults to "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified" if unset.
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
	XmlSigKeyInfoKeyNameTransformer string `json:"xmlSigKeyInfoKeyNameTransformer,omitempty"`
	// The signing certificate used to validate signatures. Required if signature validation is enabled.
	SigningCertificate string `json:"signingCertificate,omitempty"`
	// Clock skew in seconds that is tolerated when validating identity provider tokens.
	// Default value is zero.
	AllowedClockSkew int `json:"allowedClockSkew,string,omitempty"`
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
	// auto-generated
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
	Samlp *IdentityProviderConfig `json:"samlp,omitempty"`
}

type SAMLIntegrationData struct {
	SAMLSetting *SAMLSetting `json:"samlSetting,omitempty"`
}

type SAMLIntegrationConnection struct {
	Alias                  string `json:"alias,omitempty"`
	DisplayName            string `json:"displayName,omitempty"`
	Enabled                bool   `json:"enabled"`
	SingleSignOnServiceURL string `json:"singleSignOnServiceURL"`
}

type SAMLIntegrationConnections struct {
	Connections []SAMLIntegrationConnection `json:"connections,omitempty`
}

type SAMLIntegrationListResponse struct {
	Connections *SAMLIntegrationConnections `json:"connections,omitempty"`
}
