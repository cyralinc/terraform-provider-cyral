# Sidecar

This resource provides CRUD operations in Cyral sidecars, allowing users to Create, Read, Update and Delete sidecars.

## Usage

```hcl
resource "cyral_sidecar" "SOME_SOURCE_NAME" {
    name = ""
    deployment_method = ""
    publicly_accessible = true|false
    log_integration_id = ""
    metrics_integration_id = ""
}
```

## Variables

|  Name                    |  Default    |  Description                                                                         | Required |
|:-------------------------|:-----------:|:-------------------------------------------------------------------------------------|:--------:|
| `name`                   |             | Sidecar name that will be used internally in Control Plane (ex: `your_sidecar_name`) | Yes      |
| `deployment_method`      |             | Deployment method that will be used by this sidecar (valid values: `docker`, `cloudformation`, `terraform`, `helm`, `helm3`, `automated`, `custom`, `terraformGKE`) | Yes      |
| `publicly_accessible`    | `false`     | Defines if the sidecar will be publicly accessible on the internet or not (valid values: `true`, `false`)| No       |
| `log_integration_id`     | `"default"` | ID of the log integration that will be used by this sidecar.                         | No       |
| `metrics_integration_id` | `"default"` | ID of the metrics integration that will be used by this sidecar.                     | No       |
