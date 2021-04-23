# ELK Integration

CRUD operations for ELK integration.

## Usage

```hcl
resource "cyral_integration_elk" "SOME_RESOURCE_NAME" {
    name = ""
    kibana_url = ""
    es_url = ""
}
```

## Variables

|  Name         |  Default  |  Description                                                          | Required |
|:--------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`        |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `kibana_url`  |           | Kibana URL.                                                           | Yes      |
| `es_url`      |           | Elastic Search URL.                                                   | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |