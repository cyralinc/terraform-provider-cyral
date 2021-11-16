# Pager Duty Integration Resource

CRUD operations for Pager Duty integration.

## Example Usage

```hcl
resource "cyral_integration_pager_duty" "some_resource_name" {
    name = ""
    api_token = ""
}
```

## Argument Reference

| Name        | Default | Description                                                     | Required |
|:------------|:-------:|:----------------------------------------------------------------|:--------:|
| `name`      |         | Integration name that will be used internally in Control Plane. | Yes      |
| `api_token` |         | API token for the Pager Duty integration.                       | Yes      |


## Attribute Reference

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
