# Active Directory Federation Services SAML Integration Resource

Provides an ADFS SAML integration resource.

## Example Usage

### Integration with Default Configuration

```hcl
resource "cyral_integration_saml_adfs" "some_resource_name" {
  samlp {
    config {
      single_sign_on_service_url = "some_sso_url"
    }
  }
}
```

### Integration with Custom Configuration

```hcl
data "cyral_saml_configuration" "some_data_source_name" {
  saml_metadata_url = "some_metadata_url"
}

resource "cyral_integration_saml_adfs" "some_resource_name" {
  draft_alias = "some_draft_alias"
  samlp {
    provider_id = "saml"
    disabled = false
    first_broker_login_flow_alias = "SAML_First_Broker"
    post_broker_login_flow_alias = ""
    display_name = "Custom-ADFS"
    store_token = false
    add_read_token_role_on_create = false
    trust_email = false
    link_only = false
    config {
      disable_using_jwks_url = false
      sync_mode = "FORCE"
      name_id_policy_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
      principal_type = "SUBJECT"
      signature_type = "RSA_SHA256"
      saml_xml_key_name_tranformer = "KEY_ID"
      hide_on_login_page = false
      back_channel_supported = false
      disable_post_binding_response = false
      disable_post_binding_authn_request = false
      disable_post_binding_logout = false
      want_assertions_encrypted = false
      disable_force_authentication = false
      gui_order = ""
      single_sign_on_service_url = "some_sso_url"
      single_logout_service_url = ""
      xml_sig_key_info_key_name_transformer = "KEY_ID"
      signing_certificate = ""
      allowed_clock_skew = 0
      saml_metadata_url = ""
      base_64_saml_metadata_document = ""
      ldap_group_attribute = ""
    }
  }
}
```

## Argument Reference

* `samlp` - (Required) It contains the top-level configuration for an identity provider.
* `draft_alias` - (Optional) An `alias` that uniquely identifies a SAML Integration draft. If set, will delete any correspondent draft and create a new SAML Integration with the same `alias`.

The `samlp` object supports the following:

* `config` - (Required) The SAML configuration for this integration.
* `provider_id` - (Optional) This is the provider ID of `saml`.
* `disabled` - (Optional) Disable maps to Keycloak's `enabled` field. Defaults to `false`.
* `first_broker_login_flow_alias` - (Optional) Alias of authentication flow, which is triggered after `First Login` with this identity provider. Term `First Login` means that no Keycloak account is currently linked to the authenticated identity provider account.
* `post_broker_login_flow_alias` - (Optional) Alias of authentication flow, which is triggered after each login with this identity provider. Useful if you want additional verification of each user authenticated with this identity provider (for example OTP). Leave this empty if you need no any additional authenticators to be triggered after login with this identity provider.
* `display_name` - (Optional) Name of the SAML Integration displayed in the UI.
* `store_token` - (Optional) Enable if tokens must be stored after authenticating users. Defaults to `false`.
* `add_read_token_role_on_create` - (Optional) Adds read token role on creation. Defaults to `false`.
* `trust_email` - (Optional) If the identity provider supplies an email address this email address will be trusted. If the realm required email validation, users that log in from this identity provider will not have to go through the email verification process.
* `link_only` - (Optional) If true, users cannot log in through this identity provider. They can only link to this identity provider. This is useful if you don't want to allow login from the identity provider, but want to integrate with an identity provider.

The `config` object supports the following:

* `single_sign_on_service_url` - (Required) The URL that must be used to send authentication requests (SAML AuthnRequest).
* `disable_using_jwks_url` - (Optional) By default, the jwks URL is used for all SAML connections.
* `sync_mode` - (Optional) Defaults to `FORCE` if unset.
* `name_id_policy_format` - (Optional) Defaults to `urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified` if unset.
* `principal_type` - (Optional) Defaults to `SUBJECT` if unset.
* `signature_type` - (Optional) Defaults to `RSA_SHA256` if unset.
* `saml_xml_key_name_tranformer` - (Optional) Defaults to `KEY_ID` if unset.
* `hide_on_login_page` - (Optional) Defaults to `false` if unset.
* `back_channel_supported` - (Optional) Defaults to `false` if unset.
* `disable_post_binding_response` - (Optional) Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used.
* `disable_post_binding_authn_request` - (Optional) Indicates whether the AuthnRequest must be sent using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used.
* `disable_post_binding_logout` - (Optional) Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used.
* `disable_want_authn_requests_signed` - (Optional) Indicates whether the identity provider expects a signed AuthnRequest.
* `disable_want_assertions_signed` - (Optional) Indicates whether the service provider expects a signed Assertion.
* `want_assertions_encrypted` - (Optional) Indicates whether the service provider expects an encrypted Assertion.
* `disable_force_authentication` - (Optional) Indicates whether the identity provider must authenticate the presenter directly rather than rely on a previous security context.
* `disable_validate_signature` - (Optional) Enable/Disable signature validation of SAML responses. Highly recommended for minimum security.
* `gui_order` - GUI order.
* `single_logout_service_url` - (Optional) The URL that must be used to send logout requests.
* `xml_sig_key_info_key_name_transformer` - (Optional) Defaults to `KEY_ID` if unset.
* `signing_certificate` - (Optional) The signing certificate used to validate signatures. Required if signature validation is enabled.
* `allowed_clock_skew` - (Optional) Clock skew in seconds that is tolerated when validating identity provider tokens. Default value is `0`.
* `ldap_group_attribute` - (Optional) The type of `LDAP Group RDN` that identifies the name of a group within a DN. For example, if an LDAP DN sent in a SAML assertion is `cn=Everyone`, `ou=groups`, `dc=openam`, `dc=forgerock`, `dc=org` and the `LDAP Group RDN` Type is `cn` Cyral will interpret `Everyone` as the group name.

## Attribute Reference

* `id` - The ID of this resource, which corresponds to the Integration `alias`.
* `internal_id` - An ID that is auto-generated internally for this Integration.
