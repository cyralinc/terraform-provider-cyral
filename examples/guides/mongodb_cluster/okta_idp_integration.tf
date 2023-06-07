locals {
  # Set the name to be displayed for this integration in Okta's UI.
  okta_app_name = "Cyral"
  # Set the name to be displayed for this integration in Cyral's UI.
  okta_integration_name = "my-okta-integration"
}

module "cyral_idp_okta" {
  source  = "cyralinc/idp/okta"
  version = "~> 4.0"

  okta_app_name        = local.okta_app_name
  idp_integration_name = local.okta_integration_name
}

resource "cyral_repository_access_rules" "access_rules" {
  repository_id   = cyral_repository.mongodb_repo.id
  user_account_id = cyral_repository_user_account.mongodb_user_account.user_account_id
  rule {
    identity {
      type = "username"
      name = "me@mycompany.com"
    }
  }
}
