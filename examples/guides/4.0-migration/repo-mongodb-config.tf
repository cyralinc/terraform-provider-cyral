locals {
  repos = {
    mongodb = {
      host             = "mongodb.com"
      ports            = range(27017, 27020)
      type             = "mongodb"
      server_type      = "replicaset"
      replica_set_name = "replica-set-id-123"
    }
    # some other repos definitions...
  }
}

resource "cyral_repository" "all_repositories" {
  for_each = local.repos
  name  = each.key
  type  = each.value.type

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
