
### All Config enabled
resource "cyral_repository_conf_analysis" "all_conf_analysis_enabled" {
  repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
  redact = "all"
  alert_on_violation = true
  disable_pre_configured_alerts = false
  block_on_violation = true
  disable_filter_analysis = false
  enable_dataset_rewrites = true
  enable_data_masking = true
  mask_all_occurrences = true
  comment_annotation_groups = [ "identity" ]
  log_groups = [ "everything" ]
}

### All Config disabled
resource "cyral_repository_conf_analysis" "all_conf_analysis_disabled" {
  repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
  redact = "none"
  alert_on_violation = false
  disable_pre_configured_alerts = true
  block_on_violation = false
  disable_filter_analysis = true
  enable_dataset_rewrites = false
  enable_data_masking = false
  mask_all_occurrences = false
  comment_annotation_groups = []
  log_groups = []
}
