# Sidecar Credentials

Create new credentials for Cyral sidecar.

## Usage

```hcl
resource "cyral_sidecar_credentials" "SOME_RESOURCE_NAME" {
  sidecar_id = cyral_sidecar.SOME_SIDECAR_RESOURCE_NAME.id
}
```

Consider using a remote backend to encrypt the state of this resource if it sounds appropriate.

## See also

- [Remote Backends](https://www.terraform.io/docs/language/settings/backends/remote.html)

## Variables

| Name         | Default | Description                                                | Required |
| :----------- | :-----: | :--------------------------------------------------------  | :------: |
| `sidecar_id` |         | ID of the sidecar which the credentials will be generated. |   Yes    |

## Outputs

| Name            | Description                                     |
| :-------------- | :---------------------------------------------- |
| `id`            | Unique ID of the resource in the Control Plane. |
| `client_id`     | Sidecar Client ID.                              |
| `client_secret` | Sidecar Client Secret.                          |
