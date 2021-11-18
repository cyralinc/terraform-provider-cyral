# Sidecar Resource

Provides a resource to [manage sidecars](https://cyral.com/docs/sidecars/sidecar-manage).

## Example Usage

```hcl
resource "cyral_sidecar" "some_resource_name" {
    name = ""
    deployment_method = "someValidMethod"
    labels = ["label1", "label2"]
    user_endpoint = ""
}
```

## Argument Reference

* `name` - (Required) Sidecar name that will be used internally in Control Plane (ex: `your_sidecar_name`).
* `deployment_method` - (Required) Deployment method that will be used by this sidecar (valid values: `docker`, `cloudFormation`, `terraform`, `helm`, `helm3`, `automated`, `custom`, `terraformGKE`).
* `labels` - (Optional) Labels that can be attached to the sidecar and shown in the `Tags` field in the UI.
* `user_endpoint` - (Optional) User-defined endpoint (also referred as `alias`) that can be used to override the sidecar DNS endpoint shown in the UI.

## Attribute Reference

* `id` - The ID of this resource.
