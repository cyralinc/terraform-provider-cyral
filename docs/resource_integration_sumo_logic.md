# Sumo Logic Integration

CRUD operations for Sumo Logic integration.

## Usage

```hcl
resource "cyral_integration_sumo_logic" "SOME_RESOURCE_NAME" {
    name = ""
    address = ""
}
```

## Variables

|  Name         |  Default  |  Description                                                          | Required |
|:--------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`        |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `address`     |           | Sumo Logic Address.                                                   | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |