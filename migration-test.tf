terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = "~> 2.0"
    }
  }
}

provider "cyral" {
  # Follow the instructions in the Cyral Terraform Provider page to set up the
  # credentials: https://registry.terraform.io/providers/cyralinc/cyral/latest/docs
  client_id     = "sa/default/facd57b2-60db-407b-b7f2-d03cf5f828aa"
  client_secret = "FIoM9rLPBCMcD20l0J5_s_1RStmz7gyIsLss18bkZBL2Yjan"

  control_plane = "vg1221-235-a01.dev.cyral.com"

}

resource "cyral_sidecar" "pg_sidecar" {
  name              = "MainSidecar"
  deployment_method = "terraform"
}

resource "cyral_sidecar" "mongodb_sidecar" {
  name              = "MongoDBSidecar"
  deployment_method = "terraform"
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
  sidecar_id = cyral_sidecar.pg_sidecar.id
}

resource "cyral_repository" "pg_repo" {
  name = "pg_repo"
  type = "postgresql"
  host = "postgresql.mycompany.com"
  port = 5432
}

resource "cyral_repository_binding" "pg_repo_binding" {
  repository_id                 = cyral_repository.pg_repo.id
  sidecar_id                    = cyral_sidecar.pg_sidecar.id
  listener_port                 = 5432
  sidecar_as_idp_access_gateway = true
}

resource "cyral_repository" "mongodb_repo" {
  name = "mongodb_repo"
  type = "mongodb"

  # Specify the address or hostname of the endpoint of one node in the
  # MongoDB replica set. Cyral will automatically/dynamically identify
  # the remaining nodes of the replication cluster.
  host = "mycluster-shard-00-01.example.mongodb.net"

  port = 27017
  properties {
    mongodb_replica_set {
      max_nodes = 3

      # Specify the replica set identifier, a string value that
      # identifies the MongoDB replica set cluster. To find your
      # replica set ID, see our article:
      #
      # * https://cyral.freshdesk.com/a/solutions/articles/44002241594
      replica_set_id = "my-replica-set-id"
    }
  }
}

resource "cyral_repository_binding" "mongodb_repo_binding" {
  repository_id                 = cyral_repository.mongodb_repo.id
  sidecar_id                    = cyral_sidecar.mongodb_sidecar.id
  listener_port                 = 27017
  sidecar_as_idp_access_gateway = false
}

# Test Repo to use in examples.
resource "cyral_repository" "upgrade_test" {
  host = "postgresql.mycompany.com"
  name = "upgrade_test"
  port = 5432
  type = "postgresql"
}


resource "cyral_repository_local_account" "aws_iam" {
  repository_id = cyral_repository.upgrade_test.id
  aws_iam {
    database_name = "auth-db"
    local_account = "aws-iam-account"
    role_arn      = "some-role-arn"
  }
}
resource "cyral_repository_identity_map" "aws_iam_approval" {
  repository_id               = cyral_repository.upgrade_test.id
  repository_local_account_id = cyral_repository_local_account.aws_iam.id
  identity_type               = "user"
  identity_name               = "test-user-1"
  access_duration {
    days    = 5
    hours   = 0
    minutes = 0
    seconds = 0
  }
}
### AWS Secrets Manager
resource "cyral_repository_local_account" "aws_secrets_manager" {
  repository_id = cyral_repository.upgrade_test.id
  aws_secrets_manager {
    database_name = "auth-db"
    local_account = "aws-secrets-account"
    secret_arn    = "some-secret-arn"
  }
}
resource "cyral_repository_identity_map" "aws_secrets_access" {
  repository_id               = cyral_repository.upgrade_test.id
  repository_local_account_id = cyral_repository_local_account.aws_secrets_manager.id
  identity_type               = "group"
  identity_name               = "Everyone"
  access_duration {
    days    = 5
    hours   = 0
    minutes = 0
    seconds = 0
  }
}
### Hashicorp Vault
resource "cyral_repository_local_account" "vault" {
  repository_id = cyral_repository.upgrade_test.id
  hashicorp_vault {
    database_name = "auth-db"
    local_account = "vault-account"
    path          = "some-path"
  }
}
resource "cyral_repository_identity_map" "vault_id_map_approval" {
  repository_id               = cyral_repository.upgrade_test.id
  repository_local_account_id = cyral_repository_local_account.vault.id
  identity_type               = "user"
  identity_name               = "test-user-2"
  access_duration {
    days    = 10
    hours   = 4
    minutes = 2
    seconds = 50
  }
}

resource "cyral_repository_identity_map" "vault_id_map_access_rule" {
  repository_id               = cyral_repository.upgrade_test.id
  repository_local_account_id = cyral_repository_local_account.vault.id
  identity_type               = "user"
  identity_name               = "test-user-1"
}
### Environment variable
resource "cyral_repository_local_account" "env_var" {
  repository_id = cyral_repository.upgrade_test.id
  environment_variable {
    local_account = "env-var-account"
    variable_name = "CYRAL_DBSECRETS_ENV_VAR"
  }
}
resource "cyral_repository_identity_map" "env_id_map" {
  repository_id               = cyral_repository.upgrade_test.id
  repository_local_account_id = cyral_repository_local_account.env_var.id
  identity_type               = "user"
  identity_name               = "test-user-2"
}
### Kubernetes Secret
resource "cyral_repository_local_account" "kube" {
  repository_id = cyral_repository.upgrade_test.id
  kubernetes_secret {
    local_account = "kube-account"
    secret_name   = "some-secret"
    secret_key    = "some-key"
  }
}
resource "cyral_repository_identity_map" "kube_id_map_user" {
  repository_id               = cyral_repository.upgrade_test.id
  repository_local_account_id = cyral_repository_local_account.kube.id
  identity_type               = "user"
  identity_name               = "test-user-1"
}
resource "cyral_repository_identity_map" "kube_id_map_group" {
  repository_id               = cyral_repository.upgrade_test.id
  repository_local_account_id = cyral_repository_local_account.kube.id
  identity_type               = "group"
  identity_name               = "test-user-1"
}
### GCP Secret Manager
resource "cyral_repository_local_account" "gcp" {
  repository_id = cyral_repository.upgrade_test.id
  gcp_secret_manager {
    database_name = "auth-db"
    local_account = "gcp-account"
    secret_name   = "some-secret"
  }
}
