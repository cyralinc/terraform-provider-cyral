# Repository Resource

Provides a resource to handle repositories.

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

You may also use the same resource declaration to handle multiple repositories at once by using a `local` variable and `count` parameter:


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
```

## Argument Reference

* `host` - (Required): Repository host name (ex: `somerepo.cyral.com`)
* `port` - (Required): Repository access port (ex: `3306`)
* `type` - (Required): Repository type (see the list of supported types below)
  * Accepted values: 
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
* `name` - (Required): Repository name that will be used internally in Control Plane (ex: `your_repo_name`)

## Attribute Reference
* `id` - The ID of this resource.
