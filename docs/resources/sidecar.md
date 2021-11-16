# Sidecar Resource

Provides a resource to handle sidecars.

## Example Usage

```hcl
resource "cyral_sidecar" "some_resource_name" {
    name = ""
    deployment_method = "someValidMethod"
    labels = ["label1", "label2"]
}
```

## Argument Reference

* `name` - (Required) Sidecar name that will be used internally in Control Plane (ex: `your_sidecar_name`).
* `deployment_method` - (Required) Deployment method that will be used by this sidecar (valid values: `docker`, `cloudFormation`, `terraform`, `helm`, `helm3`, `automated`, `custom`, `terraformGKE`).
* `labels` - (Optional) Labels that can be attached to the sidecar and shown in the `Tags` field in the UI.

## Attribute Reference

* `id` - The ID of this resource.
