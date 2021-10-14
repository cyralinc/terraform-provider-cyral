# SAML Configuration Data Source

Parses a SAML metadata URL or a Base64 document into a SAML configuration.

## Example Usage

```hcl
# Parses metadata Base64 document to SAML configuration.
data "cyral_saml_configuration" "some_data_source_name" {
  base_64_saml_metadata_document = "some_metadata_base_64_document"
}

# Parses metadata URL to SAML configuration.
data "cyral_saml_configuration" "some_data_source_name" {
  saml_metadata_url = "some_metadata_url"
}
```

## Argument Reference

* `base_64_saml_metadata_document` - (Required, unless using `saml_metadata_url`) This is the full SAML metadata document that should be used to parse a SAML configuration, Base64 encoded.
* `saml_metadata_url` - (Required, unless using `base_64_saml_metadata_document`) This is the full SAML metadata URL we should use to parse to a SAML configuration.

## Attribute Reference

* `id` - The ID of this data source.
* `disable_using_jwks_url` - By default, the jwks URL is used for all SAML connections.
* `sync_mode` - Defaults to `FORCE` if unset.
* `name_id_policy_format` - Defaults to `urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified` if unset.
* `principal_type` - Defaults to `SUBJECT` if unset.
* `signature_type` - Defaults to `RSA_SHA256` if unset.
* `saml_xml_key_name_tranformer` - Defaults to `KEY_ID` if unset.
* `hide_on_login_page` - Defaults to `false` if unset.
* `back_channel_supported` - Defaults to `false` if unset.
* `disable_post_binding_response` - Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used.
* `disable_post_binding_authn_request` - Indicates whether the AuthnRequest must be sent using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used.
* `disable_post_binding_logout` - Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used.
* `disable_want_authn_requests_signed` - Indicates whether the identity provider expects a signed AuthnRequest.
* `disable_want_assertions_signed` - Indicates whether the service provider expects a signed Assertion.
* `want_assertions_encrypted` - Indicates whether the service provider expects an encrypted Assertion.
* `disable_force_authentication` - Indicates whether the identity provider must authenticate the presenter directly rather than rely on a previous security context.
* `disable_validate_signature` - Enable/Disable signature validation of SAML responses. Highly recommended for minimum security.
* `gui_order` - GUI order.
* `single_sign_on_service_url` - The URL that must be used to send authentication requests (SAML AuthnRequest).
* `single_logout_service_url` - The URL that must be used to send logout requests.
* `xml_sig_key_info_key_name_transformer` - Defaults to `KEY_ID` if unset.
* `signing_certificate` - The signing certificate used to validate signatures. Required if signature validation is enabled.
* `allowed_clock_skew` - Clock skew in seconds that is tolerated when validating identity provider tokens. Default value is `0`.
* `ldap_group_attribute` - The type of `LDAP Group RDN` that identifies the name of a group within a DN. For example, if an LDAP DN sent in a SAML assertion is `cn=Everyone`, `ou=groups`, `dc=openam`, `dc=forgerock`, `dc=org` and the `LDAP Group RDN` Type is `cn` Cyral will interpret `Everyone` as the group name.