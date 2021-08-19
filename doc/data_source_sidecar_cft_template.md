# Sidecar CFT Template

Returns the CloudFormation deployment template for a given sidecar.

## Usage

```hcl
data "cyral_sidecar_cft_template" "SOME_DATA_SOURCE_NAME" {
    sidecar_id = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
    log_integration_id = SOME_CYRAL_INTEGRATION.SOME_INTEGRATION_NAME.id
    metrics_integration_id = SOME_CYRAL_INTEGRATION.SOME_INTEGRATION_NAME.id
    aws_configuration {
        publicly_accessible = true|false
        key_name = "some-ec2-key-name"
    }
}
```

## Observation
This template only support sidecars with `cloudFormation` deployment method. For other methods, look up to use the generic module built into terraform.

## Variables

|  Name         |  Default  |  Description                                               | Required |
|:--------------|:---------:|:-----------------------------------------------------------|:--------:|
| `sidecar_id`  |           | ID of the sidecar which the template will be generated.                  | Yes      |
| `log_integration_id`     | `"default"` | ID of the log integration that will be used by this template.                         | No       |
| `metrics_integration_id` | `"default"` | ID of the metrics integration that will be used by this template.                     | No       |
| `aws_configuration`      |             | AWS parameters for `cloudFormation` deployment method.              | Yes       |
| `publicly_accessible`    |             | Defines a public IP and an internet-facing LB if set to `true`.                      | Yes      |
| `key_name`               |             | Key-pair name that will be associated to the sidecar EC2 instances.                  | No       |

## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `template`   |  The output variable that will contain the template.                | 