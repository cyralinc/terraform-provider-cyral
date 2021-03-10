# Repository

This resource provides CRUD operations in Cyral repositories, allowing users to Create, Read, Update and Delete repos.

## Usage

```hcl
resource "cyral_repository" "SOME_RESOURCE_NAME" {
    host = ""
    port = 0
    type = ""
    name = ""
}
```

## Variables

|  Name         |  Default  |  Description                                                                         | Required |
|:--------------|:---------:|:-------------------------------------------------------------------------------------|:--------:|
| `host`        |           | Repository host name (ex: `somerepo.cyral.com`)                                      | Yes      |
| `port`        |           | Repository access port (ex: `3306`)                                                  | Yes      |
| `type`        |           | Repository type (see the list of supported types below)                              | Yes      |
| `name`        |           | Repository name that will be used internally in Control Plane (ex: `your_repo_name`) | Yes      |

### Supported Repository Types

All types below are case sensitive:

- `bigquery`
- `cassandra`
- `dremio`
- `galera`
- `mariadb`
- `mongodb`
- `mysql`
- `oracle`
- `postgresql`
- `s3`
- `snowflake`
- `sqlserver`
