# Repository Binding

This resource binds repositories to sidecars.

## Usage

```hcl
resource "cyral_repository_binding" "SOME_RESOURCE_NAME" {
    enabled = true|false
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    sidecar_id    = cyral_repository.SOME_SIDECAR_RESOURCE_NAME.id
    listener_port = 0
    listener_host = "0.0.0.0"
}
```

You may also use the same resource declaration to bind multiple repositories at once:

```hcl
resource "cyral_repository" "mongodb_repo" {
    type = "mongodb"
    host = "mongodb.cyral.com"
    port = 27017
    name = "mymongodb"
}

resource "cyral_repository" "mariadb_repo" {
    type = "mariadb"
    host = "mariadb.cyral.com"
    port = 3307
    name = "mymariadb"
}

resource "cyral_sidecar" "my_sidecar_name" {
    name = "mysidecar"
    deployment_method = "cloudFormation"
    publicly_accessible = true
}

locals {
    repositories = [cyral_repository.mongodb_repo, cyral_repository.mariadb_repo]
}

resource "cyral_repository_binding" "repo_binding" {
    count         = length(local.repositories)
    repository_id = local.repositories[count.index].id
    listener_port = local.repositories[count.index].port
    sidecar_id    = cyral_sidecar.my_sidecar_name.id
}
```

## Variables

|  Name           |  Default    |  Description                                                                         | Required |
|:----------------|:-----------:|:-------------------------------------------------------------------------------------|:--------:|
| `enabled`       | `true`      | Enable|Disable the repository in the target sidecar. It is important to notice that the resource will always be created, but will remain inactive if set to `false`.  | No       |
| `repository_id` |             | ID of the repository that will be bound to the sidecar.                              | Yes      |
| `sidecar_id`    |             | ID of the sidecar that the repository(ies) will be bound to.                         | Yes      |
| `listener_port` |             | Port in which the sidecar will listen for the given repository.                      | Yes      |
| `listener_host` | `"0.0.0.0"` | Address in which the sidecar will listen for the given repository. By default, the sidecar will listen in all interfaces. | No       |

