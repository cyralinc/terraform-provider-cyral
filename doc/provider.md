# Provider

The provider is the base element and it must be used to inform application-wide parameters.

## Usage

```hcl
resource "cyral_repository" "SOME_SOURCE_NAME" {
    host = ""
    port = 0
    type = ""
    name = ""
    require_tls = false
}
```

## Variables

|  Name         |  Default  |  Description                                                                         | Required |
|:--------------|:---------:|:-------------------------------------------------------------------------------------|:--------:|
| `host`        |           | Repository host name (ex: `somerepo.cyral.com`)                                      | Yes      |
| `port`        |           | Repository access port (ex: `3306`)                                                  | Yes      |
| `type`        |           | Repository type (see the list of supported types below)                              | Yes      |
| `name`        |           | Repository name that will be used internally in Control Plane (ex: `your_repo_name`) | Yes      |
| `require_tls` | `false`   | Control plane host (ex: `yourcp.cyral.com`)                                          | No       |

### Supported Repository Types

All types below are case sensitive:

- `bigquery`
- `cassandra`
- `dremio`
- `galera`
- `mariadb`
- `mongodb`
- `mysql`
- `postgresql`
- `snowflake`
- `sqlserver`
