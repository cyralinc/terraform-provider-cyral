# Sidecar

CRUD operations for Cyral sidecars.

## Usage

```hcl
resource "cyral_sidecar" "SOME_RESOURCE_NAME" {
    name = ""
    deployment_method = "someValidMethod"
}
```

## Variables

|  Name                    |  Default    |  Description                                                                         | Required |
|:-------------------------|:-----------:|:-------------------------------------------------------------------------------------|:--------:|
| `name`                   |             | Sidecar name that will be used internally in Control Plane (ex: `your_sidecar_name`) | Yes      |
| `deployment_method`      |             | Deployment method that will be used by this sidecar (valid values: `docker`, `cloudFormation`, `terraform`, `helm`, `helm3`, `automated`, `custom`, `terraformGKE`) | Yes      |

## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
