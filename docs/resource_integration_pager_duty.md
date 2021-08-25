# Pager Duty Integration

CRUD operations for Pager Duty integration.

## Usage

```hcl
resource "cyral_integration_pager_duty" "SOME_RESOURCE_NAME" {
    name = ""
    api_token = ""
}
```

## Variables

| Name        | Default | Description                                                     | Required |
|:------------|:-------:|:----------------------------------------------------------------|:--------:|
| `name`      |         | Integration name that will be used internally in Control Plane. | Yes      |
| `api_token` |         | API token for the Pager Duty integration.                       | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
