locals {
  database_credentials = {
    # Native database credentials.
    username = ""
    password = ""
  }
}

resource "aws_secretsmanager_secret" "mongodb_creds" {
  # The sidecar deployed using our AWS sidecar module has access to
  # all secrets with the prefix '/cyral/' in the region it is
  # deployed.
  name = join("", [
    "/cyral/dbsecrets/",
    cyral_repository.mongodb_repo.id
  ])
}

resource "aws_secretsmanager_secret_version" "mongodb_creds_version" {
  secret_id     = aws_secretsmanager_secret.mongodb_creds.id
  secret_string = jsonencode(local.database_credentials)
}

resource "cyral_repository_user_account" "mongodb_user_account" {
  repository_id = cyral_repository.mongodb_repo.id
  # Set the name of the user account. This can be chosen freely.
  name = ""
  # Set the name of the target MongoDB database.
  auth_database_name = ""
  auth_scheme {
    aws_secrets_manager {
      secret_arn = aws_secretsmanager_secret.mongodb_creds.arn
    }
  }
}
