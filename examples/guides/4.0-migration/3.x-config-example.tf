terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = "~> 3.0"
    }
  }
}

provider "cyral" {
  control_plane = "[TENANT].app.cyral.com"
  client_id     = ""
  client_secret = ""
}

locals {
  repos = {
    postgresql = {
      host  = "postgres.cyral.com"
      ports = [5432]
      type  = "postgresql"
    }
    mongodb = {
      host             = "mongodb.cyral.com"
      ports            = range(27017, 27020)
      type             = "mongodb"
      server_type      = "replicaset"
      replica_set_name = "replica-set-id-123"
    }
    mysql = {
      host  = "mysql.cyral.com"
      ports = [3306]
      type  = "mysql"
    }
    s3 = {
      host = "s3.amazonaws.com"
      ports = [443]
      type = "s3"
    }
  }
}

resource "cyral_sidecar" "sidecar" {
  name              = "sidecar"
  deployment_method = "terraform"
}

resource "cyral_repository" "all_repositories" {
  for_each = local.repos
  name     = each.key
  host     = each.value.host
  port     = each.value.ports[0]
  type     = each.value.type

  dynamic "properties" {
    for_each = each.value.type == "mongodb" ? (each.value.server_type == "replicaset" ? [""] : []) : []
    content {
      mongodb_replica_set {
        max_nodes      = length(each.value.ports)
        replica_set_id = each.value.replica_set_name
      }
    }
  }
}

resource "cyral_repository_binding" "all_repo_binding" {
  for_each                      = cyral_repository.all_repositories
  repository_id                 = each.value.id
  listener_port                 = each.value.port
  sidecar_id                    = cyral_sidecar.sidecar.id
  sidecar_as_idp_access_gateway = true
}
