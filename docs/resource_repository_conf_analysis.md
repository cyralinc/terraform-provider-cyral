# Repository Analysis Configuration

CRUD operations for Repository Analysis Configuration.

## Usage

```hcl
resource "cyral_repository_conf_analysis" "some_conf_analysis_resource_name" {
  repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
  redact = "all"
  tag_sensitive_data = false
  ignore_identifier_case = false
  analyze_where_clause = false
  alert_on_violation = true
  disable_pre_configured_alerts = false
  block_on_violation = false
  disable_filter_analysis = false
  rewrite_on_violation = false
  comment_annotation_groups = [ "identity" ]
  log_groups = [ "everything" ]
}
```

## Variables

| Name                            | Default | Description | Required |
| :------------------------------ | :-----: | :---------- | :------: |
| `repository_id`                 |         |             |   Yes    |
| `redact`                        | `"all"` |             |    No    |
| `tag_sensitive_data`            | `false` |             |    No    |
| `ignore_identifier_case`        | `false` |             |    No    |
| `analyze_where_clause`          | `false` |             |    No    |
| `alert_on_violation`            | `true`  |             |    No    |
| `disable_pre_configured_alerts` | `false` |             |    No    |
| `block_on_violation`            | `false` |             |    No    |
| `disable_filter_analysis`       | `false` |             |    No    |
| `rewrite_on_violation`          | `false` |             |    No    |
| `comment_annotation_groups`     |         |             |    No    |
| `log_groups`                    |         |             |    No    |

## Outputs

| Name | Description                                     |
| :--- | :---------------------------------------------- |
| `id` | Unique ID of the resource in the Control Plane. |
