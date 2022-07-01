locals {
  okta_app_name         = "Cyral"
  okta_integration_name = "my-okta-integration"

  database_credentials = {
    # Native database credentials.
    username = ""
    password = ""
  }
}

module "cyral_idp_okta" {
  source = "cyralinc/idp/okta"
  version = ">= 3.0.2"

  tenant = "default"

  control_plane = "${local.control_plane}:8000"

  okta_app_name        = local.okta_app_name
  idp_integration_name = local.okta_integration_name
}

resource "aws_secretsmanager_secret" "mongodb_creds" {
  name = join("", [
    "/cyral/dbsecrets/",
    cyral_repository.mongodb_repo.id
  ])
}

resource "aws_secretsmanager_secret_version" "mongodb_creds_version" {
  secret_id     = aws_secretsmanager_secret.mongodb_creds.id
  secret_string = jsonencode(local.database_credentials)
}

resource "cyral_repository_local_account" "mongodb_local_account" {
  repository_id = cyral_repository.mongodb_repo.id
  aws_secrets_manager {
    # Set the name of the target MongoDB database.
    database_name = ""
	# Set the name of local account. This can be chosen freely.
    local_account = ""
    secret_arn    = aws_secretsmanager_secret.mongodb_creds.arn
  }
}

resource "cyral_repository_identity_map" "okta" {
  repository_id               = cyral_repository.mongodb_repo.id
  repository_local_account_id = cyral_repository_local_account.mongodb_local_account.id
  identity_type               = "user"
  identity_name               = "me@myemail.com"
}
