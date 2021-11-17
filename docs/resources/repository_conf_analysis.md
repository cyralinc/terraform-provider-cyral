# Repository Analysis Configuration Resource

Provides Repository Analysis Configuration. This resource allows configuring both
[Log Settings](https://cyral.com/docs/manage-repositories/repo-log-volume) and [Advanced settings](https://cyral.com/docs/manage-repositories/repo-advanced-settings) (Logs, Alerts, Analysis and Enforcement) configurations
for Data Repositories.

## Example Usage

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

## Argument Reference

* `repository_id` - (Required) The ID of an existing data repository resource that will be configured.
* `redact` - (Optional) Valid values are: `all`, `none` and `watched`. If set to `all` it will enable the redact of all literal values, `none` will disable it, and `watched` will only redact values from tracked fields set in the Datamap.
* `alert_on_violation` - (Optional) If set to `true` it will enable alert on policy violations.
* `disable_pre_configured_alerts` - (Optional) If set to `false` it will keep preconfigured alerts enabled.
* `block_on_violation` - (Optional) If set to `true` it will enable query blocking in case of a policy violation.
* `disable_filter_analysis` - (Optional) If set to `false` it will keep filter analysis enabled.
* `rewrite_on_violation` - (Optional) If set to `true` it will enable rewriting queries on violations.
* `comment_annotation_groups` - (Optional) Valid values are: `identity`, `client`, `repo`, `sidecar`. The default behavior is to set only the `identity` when this option is enabled, but you can also opt to add the contents of `client`, `repo`, `sidecar` logging blocks as query comments. See also [Logging additional data as comments on a query](https://support.cyral.com/support/solutions/articles/44002218978).
* `log_groups` - (Optional) Responsible for configuring the Log Settings. Valid values are documented below. The `log_groups` list support the following values:
  * `everything` - Enables all the Log Settings.
  * `dql` - Enables the `DQLs` setting for `all requests`.
  * `dml` - Enables the `DMLs` setting for `all requests`.
  * `ddl` - Enables the `DDLs` setting for `all requests`.
  * `sensitive & dql` - Enables the `DQLs` setting for `logged fields`.
  * `sensitive & dml` - Enables the `DMLs` setting for `logged fields`.
  * `sensitive & ddl` - Enables the `DDLs` setting for `logged fields`. 
  * `privileged` - Enables the `Privileged commands` setting.
  * `port-scan` - Enables the `Port scans` setting.
  * `auth-failure` - Enables the `Authentication failures` setting.
  * `full-table-scan` - Enables the `Full scans` setting.
  * `violations` - Enables the `Policy violations` setting.
  * `connections` - Enables the `Connection activity` setting.
  * `sensitive` - Log all queries manipulating sensitive fields (watches)
  * `data-classification` - Log all queries whose response was automatically classified as sensitive (credit card numbers, emails and so on).
  * `audit` - Log `sensitive`, `DQLs`, `DDLs`, `DMLs` and `privileged`.
  * `error` - Log analysis errors.
  * `new-connections` - Log new connections.
  * `closed-connections` - Log closed connections.

## Attribute Reference

* `id` - The ID of this resource.
