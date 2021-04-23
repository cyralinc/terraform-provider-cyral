# Looker Integration

CRUD operations for Datadog integration.

## Usage

```hcl
resource "cyral_integration_looker" "looker" {
    client_id = ""
    client_secret = ""
    url = ""
}
```

## Variables

|  Name           |  Default  |  Description                                                          | Required |
|:----------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `client_id`     |           | Looker client id.                                                     | Yes      |
| `client_secret` |           | Looker client secret.                                                 | Yes      |
| `url`           |           | Looker integration url.                                               | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |