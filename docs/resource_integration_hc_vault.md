# Hashicorp Vault Integration

CRUD operations for Hashicorp Vault integration.

## Usage

```hcl
resource "cyral_integration_hc_vault" "my-vault" {
  name = ""
  auth_type = ""
  server = ""
  auth_method = ""
}
```

## Variables

| Name          | Default | Description                                                     | Required |
|:--------------|:-------:|:----------------------------------------------------------------|:--------:|
| `name`        |         | Integration name that will be used internally in Control Plane. | Yes      |
| `auth_type`   |         | Authentication type for the integration.                        | Yes      |
| `server`      |         | The server on which the vault service is running.               | Yes      |
| `auth_method` |         | The authentication method for the integration.                  | Yes      |

## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
