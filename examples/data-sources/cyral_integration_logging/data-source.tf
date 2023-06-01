data "cyral_integration_logging" "some_data_source_name" {
  # Filter you can apply
  type = "CLOUDWATCH"
}

data "cyral_integration_logging" "another_data_source_name" {
  # type defaults to `ANY` if not set.
}
