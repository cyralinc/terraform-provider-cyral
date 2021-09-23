# Repository Analysis Configuration

CRUD operations for Repository Analysis Configuration.

## Note

This resource allows configuring both `Log Settings` and `Advanced` (Logs, Alerts, Analysis and Enforcement) configurations for Data Repositories.

## Usage

### All Config enable

```hcl
resource "cyral_repository_conf_analysis" "all_conf_analysis_enable" {
  repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
  redact = "all"
  alert_on_violation = true
  disable_pre_configured_alerts = false
  block_on_violation = true
  disable_filter_analysis = false
  rewrite_on_violation = true
  tag_sensitive_data = false
  ignore_identifier_case = false
  analyze_where_clause = false
  comment_annotation_groups = [ "identity" ]
  log_groups = [ "everything" ]
}
```

### All Config disable

```hcl
resource "cyral_repository_conf_analysis" "all_conf_analysis_disable" {
  repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
  redact = "none"
  alert_on_violation = false
  disable_pre_configured_alerts = true
  block_on_violation = false
  disable_filter_analysis = true
  rewrite_on_violation = false
  tag_sensitive_data = false
  ignore_identifier_case = false
  analyze_where_clause = false
  comment_annotation_groups = []
  log_groups = []
}
```

## Variables

| Name                            | Default | Description                                                                                                                                                                                                                                                                                                                                                                                                                            | Required |
| :------------------------------ | :-----: | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------: |
| `repository_id`                 |         | The ID of an existing data repository resource that will be configured.                                                                                                                                                                                                                                                                                                                                                                |   Yes    |
| `redact`                        | `"all"` | Refers to the `Redact literal values` option in the UI. Valid values are: `all`, `none` and `watched`. If set with `all` it will enable the redact of all literal values, `none` will disable it, and `watched` will only redact values from tracked fields set in the Datamap.                                                                                                                                                  |    No    |
| `alert_on_violation`            | `true`  | Refers to the `Alert on policy violations` option in the UI. If set with `true` it will enable alert on policy violations.                                                                                                                                                                                                                                                                                                          |    No    |
| `disable_pre_configured_alerts` | `false` | Refers to the `Enable preconfigured alerts` option in the UI. If set with `false` it will keep preconfigured alerts enabled.                                                                                                                                                                                                                                                                                                        |    No    |
| `block_on_violation`            | `false` | Refers to the `Block on violations` option in the UI. If set with `true` it will enable query blocking in case of a policy violation.                                                                                                                                                                                                                                                                                                   |    No    |
| `disable_filter_analysis`       | `false` | Refers to the `Perform filter analysis` option in the UI. If set with `false` it will keep filter analysis enabled.                                                                                                                                                                                                                                                                                                                 |    No    |
| `rewrite_on_violation`          | `false` | Refers to the `Rewrite queries on violations` option in the UI. If set with `true` it will enable rewriting queries on violations.                                                                                                                                                                                                                                                                                                  |    No    |
| `tag_sensitive_data`            | `false` | If set with `true` it will enable the detection of sensitive data. - `Not avaiable through the UI`.                                                                                                                                                                                                                                                                                                                                                                  |    No    |
| `ignore_identifier_case`        | `false` | This is used only for `MYSQL` data repositories. If set with `true` it will ignore the case of the query identifiers. - `Not avaiable through the UI`.                                                                                                                                                                                                                                                                                                               |    No    |
| `analyze_where_clause`          | `false` | If set with `true` it will enable the analysis of the block WHERE of the queries. - `Not avaiable through the UI`.                                                                                                                                                                                                                                                                                                                                                   |    No    |
| `comment_annotation_groups`     |   `[]`      | Refers to the `Enhance database logs` option in the UI. Valid values are: `identity`, `client`, `repo`, `sidecar`. The UI default behavior sets only the `identity` when this option is enabled, but you can also opt to add the contents of `client`, `repo`, `sidecar` logging blocks as query comments. See also [Logging additional data as comments on a query](https://support.cyral.com/support/solutions/articles/44002218978) |    No    |
| `log_groups`                    |    `[]`     | Refers to the `Log Settings` configuration section in the UI. Valid values are documented below.                                                                                                                                                                                                                                                                                                                                       |    No    |

The `log_groups` list support the following values:

| Name                  | Description                                                                               |
| :-------------------- | :---------------------------------------------------------------------------------------- |
| `everything`          | Enables all the Log Settings.                                                             |
| `dql`                 | Enables the `DQLs` setting for `all requests`.                                            |
| `dml`                 | Enables the `DMLs` setting for `all requests`.                                            |
| `ddl`                 | Enables the `DDLs` setting for `all requests`.                                            |
| `sensitive & dql`     | Enables the `DQLs` setting for `logged fields`.                                           |
| `sensitive & dml`     | Enables the `DMLs` setting for `logged fields`.                                           |
| `sensitive & ddl`     | Enables the `DDLs` setting for `logged fields`.                                           |
| `privileged`          | Enables the `Privileged commands` setting.                                                |
| `port-scan`           | Enables the `Port scans` setting.                                                         |
| `auth-failure`        | Enables the `Authentication failures` setting.                                            |
| `full-table-scan`     | Enables the `Full scans` setting.                                                         |
| `violations`          | Enables the `Policy violations` setting.                                                  |
| `connections`         | Enables the `Connection activity` setting.                                                |
| `sensitive`           | Log all queries manipulating sensitive fields (watches) - `Not avaiable through the UI`.  |
| `data-classification` | Log all queries whose response was automatically classified as sensitive (credit card numbers, emails and so on) - `Not avaiable through the UI`.                                                          |
| `audit`               | Log `sensitive`, `DQLs`, `DDLs`, `DMLs` and `privileged` - `Not avaiable through the UI`. |
| `error`               | Log analysis errors - `Not avaiable through the UI`.                                      |
| `new-connections`     | Log new connections - `Not avaiable through the UI`.                                      |
| `closed-connections`  | Log closed connections - `Not avaiable through the UI`.                                   |

## Outputs

| Name | Description                                     |
| :--- | :---------------------------------------------- |
| `id` | Unique ID of the resource in the Control Plane. |
