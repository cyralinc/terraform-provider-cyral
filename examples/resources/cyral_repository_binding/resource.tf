
### Bind a single repository
resource "cyral_repository_binding" "some_resource_name" {
    enabled = true
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    sidecar_id    = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    listener_port = 0
    listener_host = "0.0.0.0"
    sidecar_as_idp_access_gateway = false
}

### Bind multiple repositories
locals {
  repos = {
    mymongodb = {
      host          = "mongodb.cyral.com"
      port          = 27017
      type          = "mongodb"
      listener_port = 27117
    }
    mymariadb = {
      host          = "mariadb.cyral.com"
      port          = 3306
      listener_port = 3310
      type          = "mariadb"
    }
    mypostgresql = {
      host          = "postgresql.cyral.com"
      port          = 5436
      listener_port = 5432
      type          = "postgresql"
    }
  }
}

resource "cyral_repository" "repositories" {
  for_each = local.repos

  name = each.key
  type = each.value.type
  host = each.value.host
  port = each.value.port
}

resource "cyral_sidecar" "my_sidecar_name" {
  name = "mysidecar"
  tags = ["deploymentMethod:cloudFormation", "tag1"]
}

resource "cyral_repository_binding" "repo_binding" {
  for_each = local.repos

  enabled       = true
  repository_id = cyral_repository.repositories[each.key].id
  listener_port = each.value.port
  sidecar_id    = cyral_sidecar.my_sidecar_name.id
}