# Return all existing CloudWatch integrations
data "cyral_integration_logging" "some_data_source_name" {
  type = "CLOUDWATCH"
}

# Returns all existing integrations. Attribute `type` defaults to `ANY` if not set.
data "cyral_integration_logging" "another_data_source_name" {

}
