# Sidecar

This resource provides CRUD operations in Cyral sidecars, allowing users to Create, Read, Update and Delete sidecars.

## Usage

```hcl
resource "cyral_sidecar" "SOME_SOURCE_NAME" {
    name = ""
}
```

## Variables

|  Name         |  Default  |  Description                                                                         | Required |
|:--------------|:---------:|:-------------------------------------------------------------------------------------|:--------:|
| `name`        |           | Sidecar name that will be used internally in Control Plane (ex: `your_sidecar_name`) | Yes      |
