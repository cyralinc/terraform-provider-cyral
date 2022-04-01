data "cyral_sidecar_cft_template" "some_data_source_name" {
    sidecar_id = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    log_integration_id = SOME_CYRAL_INTEGRATION.SOME_INTEGRATION_NAME.id
    metrics_integration_id = SOME_CYRAL_INTEGRATION.SOME_INTEGRATION_NAME.id
    aws_configuration {
        publicly_accessible = true|false
        key_name = "some-ec2-key-name"
    }
}
