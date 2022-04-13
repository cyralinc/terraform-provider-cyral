### Role with Single SSO Group
resource "cyral_integration_idp_okta" "some_idp_okta" {
  samlp {
    config {
      single_sign_on_service_url = "https://some-sso-url.com"
    }
  }
}

resource "cyral_role" "some_role" {
  name="some-role-name"
}

resource "cyral_role_sso_groups" "some_role_sso_groups" {
  role_id=cyral_role.some_role.id
  sso_group {
    group_name="Everyone"
    idp_id=cyral_integration_idp_okta.some_idp_okta.id
  }
}

### Role with Multiple SSO Groups
resource "cyral_integration_idp_okta" "some_idp_okta" {
  samlp {
    config {
      single_sign_on_service_url = "https://some-sso-url.com"
    }
  }
}

resource "cyral_role" "some_role" {
  name="some-role-name"
}

resource "cyral_role_sso_groups" "some_role_sso_groups" {
  role_id=cyral_role.some_role.id
  sso_group {
    group_name="Admin"
    idp_id=cyral_integration_idp_okta.some_idp_okta.id
  }
  sso_group {
    group_name="Dev"
    idp_id=cyral_integration_idp_okta.some_idp_okta.id
  }
}