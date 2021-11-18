# Active Directory Federation Services IdP Integration Resource

Provides [integration with Active Directory Federation Services](https://cyral.com/docs/sso/sso-adfs#add-your-adfs-as-an-idp-in-cyral) identity provider to 
allow single-sign on to Cyral.

## Example Usage

### Integration with Default Configuration

```hcl
resource "cyral_integration_idp_adfs" "some_resource_name" {
  samlp {
    config {
      single_sign_on_service_url = "some_sso_url"
    }
  }
}
```

### Integration using SAML Configuration Data Source

```hcl
locals {
  config = data.cyral_saml_configuration.some_data_source_name
}

data "cyral_saml_configuration" "some_data_source_name" {
  saml_metadata_url = "some_metadata_url"
}

resource "cyral_integration_idp_adfs" "some_resource_name" {
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
      single_sign_on_service_url = local.config.single_sign_on_service_url
      single_logout_service_url = local.config.single_logout_service_url == "" ? null : local.config.single_logout_service_url
      disable_using_jwks_url = local.config.disable_using_jwks_url
      sync_mode = local.config.sync_mode == "" ? null : local.config.sync_mode
      name_id_policy_format = local.config.name_id_policy_format == "" ? null : local.config.name_id_policy_format
      principal_type = local.config.principal_type == "" ? null : local.config.principal_type
      signature_type = local.config.signature_type == "" ? null : local.config.signature_type
      saml_xml_key_name_tranformer = local.config.saml_xml_key_name_tranformer == "" ? null : local.config.saml_xml_key_name_tranformer
      hide_on_login_page = local.config.hide_on_login_page
      back_channel_supported = local.config.back_channel_supported
      disable_post_binding_response = local.config.disable_post_binding_response
      disable_post_binding_authn_request = local.config.disable_post_binding_authn_request
      disable_post_binding_logout = local.config.disable_post_binding_logout
      want_assertions_encrypted = local.config.want_assertions_encrypted
      disable_force_authentication = local.config.disable_force_authentication
      gui_order = local.config.gui_order == "" ? null : local.config.gui_order
      xml_sig_key_info_key_name_transformer = local.config.xml_sig_key_info_key_name_transformer == "" ? null : local.config.xml_sig_key_info_key_name_transformer
      signing_certificate = local.config.signing_certificate == "" ? null : local.config.signing_certificate
      allowed_clock_skew = local.config.allowed_clock_skew
      saml_metadata_url = local.config.saml_metadata_url == "" ? null : local.config.saml_metadata_url
      base_64_saml_metadata_document = local.config.base_64_saml_metadata_document == "" ? null : local.config.base_64_saml_metadata_document
      ldap_group_attribute = local.config.ldap_group_attribute == "" ? null : local.config.ldap_group_attribute
    }
  }
}
```
-> When using the [SAML Configuration Data Source](../data-sources/saml_configuration.md) to configure this IdP Integration resource, consider verifying if the `string` attributes are `empty` like in the example above so that the resource arguments can be used with their default values, instead of setting them as empty.

## Argument Reference

* `samlp` - (Required) It contains the top-level configuration for an identity provider.

The `samlp` object supports the following:

* `config` - (Required) The SAML configuration for this IdP Integration.
* `provider_id` - (Optional) This is the provider ID of `saml`. Defaults to `saml`.
* `disabled` - (Optional) Disable maps to Keycloak's `enabled` field. Defaults to `false`.
* `first_broker_login_flow_alias` - (Optional) Alias of authentication flow, which is triggered after `First Login` with this identity provider. Term `First Login` means that no Keycloak account is currently linked to the authenticated identity provider account. Defaults to `SAML_First_Broker`.
* `post_broker_login_flow_alias` - (Optional) Alias of authentication flow, which is triggered after each login with this identity provider. Useful if you want additional verification of each user authenticated with this identity provider (for example OTP). Leave this empty if you need no any additional authenticators to be triggered after login with this identity provider. Defaults to `""`.
* `display_name` - (Optional) Name of the IdP Integration displayed in the UI. Defaults to `Active Directory`.
* `store_token` - (Optional) Enable if tokens must be stored after authenticating users. Defaults to `false`.
* `add_read_token_role_on_create` - (Optional) Adds read token role on creation. Defaults to `false`.
* `trust_email` - (Optional) If the identity provider supplies an email address this email address will be trusted. If the realm required email validation, users that log in from this identity provider will not have to go through the email verification process. Defaults to `false`.
* `link_only` - (Optional) If true, users cannot log in through this identity provider. They can only link to this identity provider. This is useful if you don't want to allow login from the identity provider, but want to integrate with an identity provider. Defaults to `false`.

The `config` object supports the following:

* `single_sign_on_service_url` - (Required) The URL that must be used to send authentication requests (SAML AuthnRequest).
* `disable_using_jwks_url` - (Optional) By default, the jwks URL is used for all SAML connections. Defaults to `false`.
* `sync_mode` - (Optional) Defaults to `FORCE` if unset.
* `name_id_policy_format` - (Optional) Defaults to `urn:oasis:names:tc:SAML:2.0:nameid-format:transient` if unset.
* `principal_type` - (Optional) Defaults to `SUBJECT` if unset.
* `signature_type` - (Optional) Defaults to `RSA_SHA256` if unset.
* `saml_xml_key_name_tranformer` - (Optional) Defaults to `CERT_SUBJECT` if unset.
* `hide_on_login_page` - (Optional) Defaults to `false` if unset.
* `back_channel_supported` - (Optional) Defaults to `false` if unset.
* `disable_post_binding_response` - (Optional) Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used. Defaults to `false`.
* `disable_post_binding_authn_request` - (Optional) Indicates whether the AuthnRequest must be sent using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used. Defaults to `false`.
* `disable_post_binding_logout` - (Optional) Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` binding will be used. Defaults to `true`.
* `want_assertions_encrypted` - (Optional) Indicates whether the service provider expects an encrypted Assertion. Defaults to `false`.
* `disable_force_authentication` - (Optional) Indicates whether the identity provider must authenticate the presenter directly rather than rely on a previous security context. Defaults to `false`.
* `gui_order` - GUI order. Defaults to `""`.
* `single_logout_service_url` - (Optional) The URL that must be used to send logout requests. Defaults to `""`.
* `xml_sig_key_info_key_name_transformer` - (Optional) Defaults to `CERT_SUBJECT` if unset.
* `signing_certificate` - (Optional) The signing certificate used to validate signatures. Required if signature validation is enabled. Defaults to `""`.
* `allowed_clock_skew` - (Optional) Clock skew in seconds that is tolerated when validating identity provider tokens. Defaults to `0`.
* `saml_metadata_url` - (Optional) This is the full SAML metadata URL that was used to import the SAML configuration. Defaults to `""`.
* `base_64_saml_metadata_document` - (Optional) This is the full SAML metadata document that was used to import the SAML configuration, Base64 encoded. Defaults to `""`.
* `ldap_group_attribute` - (Optional) The type of `LDAP Group RDN` that identifies the name of a group within a DN. For example, if an LDAP DN sent in a SAML assertion is `cn=Everyone`, `ou=groups`, `dc=openam`, `dc=forgerock`, `dc=org` and the `LDAP Group RDN` Type is `cn` Cyral will interpret `Everyone` as the group name. Defaults to `""`.

## Attribute Reference

* `id` - The ID of this resource, which corresponds to the IdP Integration `alias`.
* `internal_id` - An ID that is auto-generated internally for this IdP Integration.
