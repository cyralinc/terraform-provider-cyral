# Repository

CRUD operations for Microsoft Teams integration.

## Usage

```hcl
resource "cyral_integration_microsoft_teams" "SOME_RESOURCE_NAME" {
    name = ""
    url = ""
}
```

## Variables

|  Name         |  Default  |  Description                                                          | Required |
|:--------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`        |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `url`         |           | Microsoft Teams webhook URL.                                          | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |