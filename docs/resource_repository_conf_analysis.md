# Repository Analysis Configuration

CRUD operations for Repository Analysis Configuration.

## Note

This resource allows configuring both `Log Settings` and `Advanced` (Logs, Alerts, Analysis and Enforcement) configurations for Data Repositories.

## Usage

### All Config enabled

```hcl
resource "cyral_repository_conf_analysis" "all_conf_analysis_enabled" {
  repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
  redact = "all"
  alert_on_violation = true
  disable_pre_configured_alerts = false
  block_on_violation = true
  disable_filter_analysis = false
  rewrite_on_violation = true
  comment_annotation_groups = [ "identity" ]
  log_groups = [ "everything" ]
}
```

### All Config disabled

```hcl
resource "cyral_repository_conf_analysis" "all_conf_analysis_disabled" {
  repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
  redact = "none"
  alert_on_violation = false
  disable_pre_configured_alerts = true
  block_on_violation = false
  disable_filter_analysis = true
  rewrite_on_violation = false
  comment_annotation_groups = []
  log_groups = []
}
```

## Variables

| Name                            | Default | Description                                                                                                                                                                                                                                                                                                                                                                       | Required |
| :------------------------------ | :-----: | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------: |
| `repository_id`                 |         | The ID of an existing data repository resource that will be configured.                                                                                                                                                                                                                                                                                                           |   Yes    |
| `redact`                        | `"all"` | Valid values are: `all`, `none` and `watched`. If set to `all` it will enable the redact of all literal values, `none` will disable it, and `watched` will only redact values from tracked fields set in the Datamap.                                                                                                                                                             |    No    |
| `alert_on_violation`            | `true`  | If set to `true` it will enable alert on policy violations.                                                                                                                                                                                                                                                                                                                       |    No    |
| `disable_pre_configured_alerts` | `false` | If set to `false` it will keep preconfigured alerts enabled.                                                                                                                                                                                                                                                                                                                      |    No    |
| `block_on_violation`            | `false` | If set to `true` it will enable query blocking in case of a policy violation.                                                                                                                                                                                                                                                                                                     |    No    |
| `disable_filter_analysis`       | `false` | If set to `false` it will keep filter analysis enabled.                                                                                                                                                                                                                                                                                                                           |    No    |
| `rewrite_on_violation`          | `false` | If set to `true` it will enable rewriting queries on violations.                                                                                                                                                                                                                                                                                                                  |    No    |
| `comment_annotation_groups`     |  `[]`   | Valid values are: `identity`, `client`, `repo`, `sidecar`. The default behavior is to set only the `identity` when this option is enabled, but you can also opt to add the contents of `client`, `repo`, `sidecar` logging blocks as query comments. See also [Logging additional data as comments on a query](https://support.cyral.com/support/solutions/articles/44002218978). |    No    |
| `log_groups`                    |  `[]`   | Responsible for configuring the Log Settings. Valid values are documented below.                                                                                                                                                                                                                                                                                                  |    No    |

The `log_groups` list support the following values:

| Name                  | Description                                                                                                       |
| :-------------------- | :---------------------------------------------------------------------------------------------------------------- |
| `everything`          | Enables all the Log Settings.                                                                                     |
| `dql`                 | Enables the `DQLs` setting for `all requests`.                                                                    |
| `dml`                 | Enables the `DMLs` setting for `all requests`.                                                                    |
| `ddl`                 | Enables the `DDLs` setting for `all requests`.                                                                    |
| `sensitive & dql`     | Enables the `DQLs` setting for `logged fields`.                                                                   |
| `sensitive & dml`     | Enables the `DMLs` setting for `logged fields`.                                                                   |
| `sensitive & ddl`     | Enables the `DDLs` setting for `logged fields`.                                                                   |
| `privileged`          | Enables the `Privileged commands` setting.                                                                        |
| `port-scan`           | Enables the `Port scans` setting.                                                                                 |
| `auth-failure`        | Enables the `Authentication failures` setting.                                                                    |
| `full-table-scan`     | Enables the `Full scans` setting.                                                                                 |
| `violations`          | Enables the `Policy violations` setting.                                                                          |
| `connections`         | Enables the `Connection activity` setting.                                                                        |
| `sensitive`           | Log all queries manipulating sensitive fields (watches)                                                           |
| `data-classification` | Log all queries whose response was automatically classified as sensitive (credit card numbers, emails and so on). |
| `audit`               | Log `sensitive`, `DQLs`, `DDLs`, `DMLs` and `privileged`.                                                         |
| `error`               | Log analysis errors.                                                                                              |
| `new-connections`     | Log new connections.                                                                                              |
| `closed-connections`  | Log closed connections.                                                                                           |

## Outputs

| Name | Description                                     |
| :--- | :---------------------------------------------- |
| `id` | Unique ID of the resource in the Control Plane. |
