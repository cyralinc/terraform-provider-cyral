# Repository Binding Resource

Allows [binding repositories to sidecars](https://cyral.com/docs/sidecars/sidecar-assign-repo).

## Example Usage

### Bind a single repository

```hcl
resource "cyral_repository_binding" "some_resource_name" {
    enabled = true|false
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    sidecar_id    = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    listener_port = 0
    listener_host = "0.0.0.0"
}
```

### Bind multiple repositories

It is possible to create and bind multiple repositories at once by using a `local` variable and `for_each` parameter:

```hcl
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
```

## Argument Reference

- `enabled` - (Optional) Enable|Disable the repository in the target sidecar. It is important to notice that the resource will always be created, but will remain inactive if set to `false`.
- `repository_id` - (Required) ID of the repository that will be bound to the sidecar.
- `sidecar_id` - (Required) ID of the sidecar that the repository(ies) will be bound to.
- `listener_port` - (Required) Port in which the sidecar will listen for the given repository.
- `listener_host` - (Optional) Address in which the sidecar will listen for the given repository. By default, the sidecar will listen in all interfaces.

## Attribute Reference

- `id` - The ID of this resource.
