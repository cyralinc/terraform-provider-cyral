# All information related to the database that will be mapped
# to a cyral repository is defined here for clarity, but you
# can define somewhere else.
locals {
    database_name = "someDatabaseName"
    database_credentials = {
        username = "someUserName"
        password = "somePassword"
    }
}

# See the Cyral provider documentation for more
# information on how to initialize it correctly.
provider "cyral" {
    control_plane = "mycontrolplane.cyral.com:8000"
}

resource "cyral_repository" "mongodb_repo" {
    type = "mongodb"
    host = "mongodb.cyral.com"
    port = 27017
    name = "mymongodb"
}

resource "cyral_repository_local_account" "my_repo_account" {
    repository_id = cyral_repository.mongodb_repo.id
    aws_secrets_manager {
        database_name = local.database_name
        local_account = local.database_credentials.username
        secret_arn = aws_secretsmanager_secret.my_repository_secret.arn
    }
}

# See the AWS provider documentation for more
# information on how to initialize it correctly.
provider "aws" {
    # By deploying the secret to the same account and region of your
    # sidecar and using the name suggested in my_repository_secret,
    # the sidecar will gain access to the secret automatically.
    region = "us-east-1"
}

resource "aws_secretsmanager_secret" "my_repository_secret" {
    name = join("", [
      "/cyral/dbsecrets/",
      cyral_repository.mongodb_repo.id
    ])
}

resource "aws_secretsmanager_secret_version" "my_repository_secret_version" {
    secret_id     = aws_secretsmanager_secret.my_repository_secret.id
    secret_string = jsonencode(local.database_credentials)
}
