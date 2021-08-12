# Sidecar Credentials

CRUD operations for Cyral sidecar credentials.

## Usage

```hcl
resource "cyral_sidecar_credentials" "SOME_RESOURCE_NAME" {
  sidecar_id = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
}
```

## Variables

| Name         | Default | Description                                                | Required |
| :----------- | :-----: | :--------------------------------------------------------  | :------: |
| `sidecar_id` |         | ID of the sidecar which the credentials will be generated. |   Yes    |

## Outputs

| Name            | Description                                     |
| :-------------- | :---------------------------------------------- |
| `id`            | Unique ID of the resource in the Control Plane. |
| `client_id`     | Sidecar Client ID.                              |
| `client_secret` | Sidecar Client Secret encoded using base 64     |
