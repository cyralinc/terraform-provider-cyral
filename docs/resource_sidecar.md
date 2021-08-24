# Sidecar

CRUD operations for Cyral sidecars.

## Usage

```hcl
resource "cyral_sidecar" "SOME_RESOURCE_NAME" {
    name = ""
    deployment_method = ""
    log_integration_id = ""
    metrics_integration_id = ""
    aws_configuration {
        publicly_accessible = true|false
        aws_region = "us-east-1"
        key_name = "ec2-key-name"
        vpc = "vpc-id"
        subnets = "subnetid1,subnetid2,subnetidN"
    }
}
```

## Variables

|  Name                    |  Default    |  Description                                                                         | Required |
|:-------------------------|:-----------:|:-------------------------------------------------------------------------------------|:--------:|
| `name`                   |             | Sidecar name that will be used internally in Control Plane (ex: `your_sidecar_name`) | Yes      |
| `deployment_method`      |             | Deployment method that will be used by this sidecar (valid values: `docker`, `cloudFormation`, `terraform`, `helm`, `helm3`, `automated`, `custom`, `terraformGKE`) | Yes      |
| `log_integration_id`     | `"default"` | ID of the log integration that will be used by this sidecar.                         | No       |
| `metrics_integration_id` | `"default"` | ID of the metrics integration that will be used by this sidecar.                     | No       |
| `aws_configuration`      |             | AWS parameters for `cloudFormation` and `terraform` deployment methods.              | No       |
| `publicly_accessible`    |             | Defines a public IP and an internet-facing LB if set to `true`.                      | Yes      |
| `aws_region`             |             | AWS region that will be used to deploy the sidecar.                                  | No       |
| `key_name`               |             | Key-pair name that will be associated to the sidecar EC2 instances.                  | No       |
| `vpc`                    |             | ID of the VPC that the sidecar will be deployed to.                                  | No       |
| `subnets`                |             | Comma-separated list of subnet ids that the sidecar will be deployed to.             | No       |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
