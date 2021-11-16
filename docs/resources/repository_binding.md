# Repository Binding Resource

Allows binding repositories to sidecars.

## Example Usage

### Bind a single repository

```hcl
resource "cyral_repository_binding" "some_resource_name" {
    enabled = true|false
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    sidecar_id    = cyral_repository.SOME_SIDECAR_RESOURCE_NAME.id
    listener_port = 0
    listener_host = "0.0.0.0"
}
```

### Bind multiple repositories

It is possible to create and bind multiple repositories at once by using a `local` variable and `count` parameter:

```hcl
locals {
    repos = [
        ["mongodb", "mongodb.cyral.com", 27017, "mymongodb"],
        ["mariadb", "mariadb.cyral.com", 3310, "mymariadb"],
        ["postgresql", "postgresql.cyral.com", 5432, "mypostgresql"]
    ]
}

resource "cyral_repository" "repositories" {
    count = length(local.repos)

    type  = local.repos[count.index][0]
    host  = local.repos[count.index][1]
    port  = local.repos[count.index][2]
    name  = local.repos[count.index][3]
}

resource "cyral_sidecar" "my_sidecar_name" {
    name = "mysidecar"
    tags = ["deploymentMethod:cloudFormation", "tag1"]
}

resource "cyral_repository_binding" "repo_binding" {
    enabled       = true
    count         = length(local.repos)
    repository_id = cyral_repository.repositories[count.index].id
    listener_port = cyral_repository.repositories[count.index].port
    sidecar_id    = cyral_sidecar.my_sidecar_name.id
}
```

## Argument Reference

* `enabled` - (Optional) Enable|Disable the repository in the target sidecar. It is important to notice that the resource will always be created, but will remain inactive if set to `false`.
* `repository_id` - (Required) ID of the repository that will be bound to the sidecar.
* `sidecar_id` - (Required) ID of the sidecar that the repository(ies) will be bound to.
* `listener_port` - (Required) Port in which the sidecar will listen for the given repository.
* `listener_host` - (Optional) Address in which the sidecar will listen for the given repository. By default, the sidecar will listen in all interfaces.

## Attribute Reference

* `id` - The ID of this resource.
