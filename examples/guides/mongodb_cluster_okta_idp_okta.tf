locals {
  okta_app_name         = "Cyral"
  okta_integration_name = "my-okta-integration"
}

module "cyral_idp_okta" {
  source = "cyralinc/idp/okta"
  version = "~> 3.0"

  tenant = "default"

  control_plane = "${local.control_plane}:8000"

  okta_app_name        = local.okta_app_name
  idp_integration_name = local.okta_integration_name
}

resource "cyral_repository_identity_map" "okta" {
  repository_id               = cyral_repository.mongodb_repo.id
  repository_local_account_id = cyral_repository_local_account.mongodb_local_account.id
  identity_type               = "user"
  identity_name               = "me@myemail.com"
}
