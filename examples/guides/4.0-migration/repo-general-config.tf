locals {
  repos = {
    postgresql = {
      host  = "postgres.com"
      ports = [5432]
      type  = "postgresql"
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
}
