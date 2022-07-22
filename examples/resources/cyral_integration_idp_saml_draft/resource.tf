resource "cyral_integration_idp_saml_draft" "example_draft" {
  display_name = "example-okta-integration"
  idp_type     = "okta"
  attributes {
    first_name = "some-first-name"
    last_name  = "some-last-name"
    email      = "some-email"
    groups     = "some-group"
  }
  disable_idp_initiated_login = false
}
