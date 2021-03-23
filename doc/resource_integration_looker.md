# Repository

CRUD operations for Datadog integration.

## Usage

```hcl
resource "cyral_integration_looker" "looker" {
    name = ""
    client_id = ""
    client_secret = ""
    url = ""
}
```

## Variables

|  Name         |  Default  |  Description                                                          | Required |
|:--------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`        |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `client_id`        |           | Looker client id.       | Yes      |
| `client_secret`        |           | Looker client secret.       | Yes      |
| `url`        |           | Looker integration url.       | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |