---
page_title: "Deploy Native Repository Credentials to AWS Secrets Manager"
---

Use the following code as an example to deploy native repository credentials to
AWS Secrets Manager and also to associate it to a Cyral Repository that can be
later bound to a sidecar.

```terraform
# All information related to the database that will be mapped to a
# cyral repository is defined here for clarity, but you can define
# somewhere else.
locals {
  # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
  control_plane_host = "[TENANT].app.cyral.com"
  control_plane_port = 443

  database_name = ""
  database_credentials = {
    username = ""
    password = ""
  }
}

terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = "~> 3.0"
    }
  }
}

# Follow the instructions in the Cyral Terraform Provider page to set
# up the Control Plane credentials:
#
# * https://registry.terraform.io/providers/cyralinc/cyral/latest/docs
provider "cyral" {
  client_id     = ""
  client_secret = ""
  control_plane = "${local.control_plane_host}:${local.control_plane_port}"
}

# See the AWS provider documentation for more information on how to
# initialize it correctly.
provider "aws" {
  # By deploying the secret to the same account and region of your
  # sidecar and using the name suggested in my_repository_secret, the
  # sidecar will gain access to the secret automatically.
  region = "us-east-1"
}

resource "cyral_repository" "mongodb_repo" {
  type = "mongodb"
  name = "mymongodb"
  repo_node {
    host = "mongodb.mycompany.com"
    port = 27017
  }
  mongodb_settings {
    server_type = "standalone"
  }
}

resource "cyral_repository_user_account" "my_user_account" {
  repository_id      = cyral_repository.mongodb_repo.id
  name               = local.database_credentials.username
  auth_database_name = local.database_name
  auth_scheme {
    aws_secrets_manager {
      secret_arn = aws_secretsmanager_secret.my_repository_secret.arn
    }
  }
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
```
