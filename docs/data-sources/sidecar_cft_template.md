# Sidecar CFT Template

Retrieves the CloudFormation deployment template for a given sidecar.

This data source only support sidecars with `cloudFormation` deployment method. For other methods, look up to use the generic module built into terraform.

## Example Usage

```hcl
data "cyral_sidecar_cft_template" "some_data_source_name" {
    sidecar_id = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    log_integration_id = SOME_CYRAL_INTEGRATION.SOME_INTEGRATION_NAME.id
    metrics_integration_id = SOME_CYRAL_INTEGRATION.SOME_INTEGRATION_NAME.id
    aws_configuration {
        publicly_accessible = true|false
        key_name = "some-ec2-key-name"
    }
}
```

## Argument Reference

* `sidecar_id` - (Required) ID of the sidecar which the template will be generated.
* `log_integration_id` - (Optional) ID of the log integration that will be used by this template.
* `metrics_integration_id` - (Optional) ID of the metrics integration that will be used by this template.
* `aws_configuration` - (Required) AWS parameters for `cloudFormation` deployment method.
* `publicly_accessible` - (Required) Defines a public IP and an internet-facing LB if set to `true`.
* `key_name` - (Optional) Key-pair name that will be associated to the sidecar EC2 instances.

## Attribute Reference

* `template` - The output variable that will contain the template.