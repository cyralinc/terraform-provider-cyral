resource "cyral_integration_idp_saml_draft" "example_draft" {
  display_name = "example-okta-integration"
  idp_type     = "okta"
}

# Here goes IdP provider configuration to obtain SAML metadata. Use the
# attribute `sp_metadata` of the SAML draft above to obtain the Cyral SP SAML
# metadata you will need to provide your IdP with.
#
# ...
##

resource "cyral_integration_idp_saml" "example_integration" {
  saml_draft_id = cyral_integration_idp_saml_draft.example_draft.id
  # This is the IdP metadata URL. You may choose instead to provide the
  # base64-encoded metadata XML document using the argument
  # `idp_metadata_document`.
  idp_metadata_url = "https://dev-123456.okta.com/app/1234567890/sso/saml/metadata"
}
