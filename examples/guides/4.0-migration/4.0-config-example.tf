terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = "~> 4.0"
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
      mysql_settings = {
        character_set = ""
        db_version = ""
      }
    }
    s3 = {
      host = "s3.amazonaws.com"
      ports = [443]
      type = "s3"
      s3_settings = {
        proxy_mode = true
      }
    }
  }
  type_port_map = merge([
    for key, repo in local.repos : {
      for port in repo.ports :
      "${repo.type}_${port}" => {
          type = repo.type
          port = port
      }
    }
  ]...)
}

resource "cyral_sidecar" "sidecar" {
  name              = "sidecar"
  deployment_method = "terraform"
}

resource "cyral_repository" "all_repositories" {
  for_each = local.repos
  name     = each.key
  type     = each.value.type

  connection_draining {
    auto      = false
    wait_time = 0
  }

  dynamic "repo_node" {
    for_each = range(0, length(each.value.ports))
    content {
      dynamic = repo_node.value == 0 ? false : true
      host = repo_node.value == 0 ? each.value.host : ""
      port = repo_node.value == 0 ? each.value.ports[0] : 0
    }
  }

  dynamic "mongodb_settings" {
    for_each = each.value.type == "mongodb" ? (each.value.server_type == "replicaset" ? [""] : []) : []
    content {
      replica_set_name = each.value.replica_set_name
      server_type = each.value.server_type
    }
  }
}

resource "cyral_sidecar_listener" "sidecar_all_listeners" {
  for_each = local.type_port_map
  repo_types = [each.value.type]
  sidecar_id = cyral_sidecar.sidecar.id

  network_address {
    port = each.value.port
  }

  dynamic "mysql_settings" {
    for_each = each.value.type == "mysql" ? [""] : []
    content {
      character_set = local.repos[each.value.type].mysql_settings.character_set
      db_version = local.repos[each.value.type].mysql_settings.db_version
    }
  }

  dynamic "s3_settings" {
    for_each = each.value.type == "s3" ? [""] : []
    content {
      proxy_mode = local.repos[each.value.type].s3_settings.proxy_mode
    }
  }
}

resource "cyral_repository_binding" "all_repo_binding" {
  for_each = local.repos
  repository_id = cyral_repository.all_repositories[each.key].id
  sidecar_id = cyral_sidecar.sidecar.id

  dynamic "listener_binding" {
    for_each = each.value.ports
    content {
      listener_id = cyral_sidecar_listener.sidecar_all_listeners["${each.value.type}_${listener_binding.value}"].listener_id
      node_index = listener_binding.key
    }
  }
}

resource "cyral_repository_access_gateway" "all_repo_binding_all_access_gateways" {
  for_each = local.repos
  repository_id = cyral_repository.all_repositories[each.key].id
  sidecar_id = cyral_sidecar.sidecar.id
  binding_id = cyral_repository_binding.all_repo_binding[each.key].binding_id
}
