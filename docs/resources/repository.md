# Repository Resource

Provides a resource to [track repositories](https://cyral.com/docs/manage-repositories/repo-track).

See also [Cyral Repository Configuration Module](https://github.com/cyralinc/terraform-cyral-repository-config).
This module provides the repository configuration options as shown in Cyral UI.

## Example Usage

```hcl
resource "cyral_repository" "some_resource_name" {
    host = ""
    port = 0
    type = ""
    name = ""
}
```

You may also use the same resource declaration to handle multiple repositories at once by using a `local` variable and `for_each` parameter:

```hcl
locals {
  repos = {
    mymongodb = {
      host = "mongodb.cyral.com"
      port = 27017
      type = "mongodb"
    }
    mymariadb = {
      host = "mariadb.cyral.com"
      port = 3310
      type = "mariadb"
    }
    mypostgresql = {
      host = "postgresql.cyral.com"
      port = 5432
      type = "postgresql"
    }
  }
}

resource "cyral_repository" "repositories" {
  for_each = local.repos

  name  = each.key
  type  = each.value.type
  host  = each.value.host
  port  = each.value.port
}
```

## Argument Reference

- `host` - (Required): Repository host name (ex: `somerepo.cyral.com`)
- `port` - (Required): Repository access port (ex: `3306`)
- `type` - (Required): Repository type (see the list of supported types below)
  - Accepted values:
    - `bigquery`
    - `cassandra`
    - `denodo`
    - `dremio`
    - `galera`
    - `mariadb`
    - `mongodb`
    - `mysql`
    - `oracle`
    - `postgresql`
    - `redshift`
    - `s3`
    - `snowflake`
    - `sqlserver`
- `name` - (Required): Repository name that will be used internally in Control Plane (ex: `your_repo_name`)

## Attribute Reference

- `id` - The ID of this resource.
