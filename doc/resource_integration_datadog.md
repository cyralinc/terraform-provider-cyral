# Repository

CRUD operations for Datadog integration.

## Usage

```hcl
resource "cyral_integration_datadog" "SOME_RESOURCE_NAME" {
    name = ""
    api_key = ""
}
```

## Variables

|  Name         |  Default  |  Description                                                          | Required |
|:--------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`        |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `api_key`     |           | Datadog API key.                                                      | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
