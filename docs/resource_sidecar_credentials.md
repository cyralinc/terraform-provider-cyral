# Sidecar Credentials

Create new credentials for Cyral sidecar.

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
| `client_secret` | Sidecar Client Secret.                          |

## Remote backend
Consider using a remote backend to encrypt the state of this resource if it sounds appropriate. For instance:

```hcl
terraform {
  required_providers {
    cyral = {
      source = "cyralinc/cyral"
    }
  }
  backend "s3" {
    bucket = "some-s3-bucket"
    key    = "terraform-state.json"
    region = "us-east-1"
    encrypt = true
    dynamodb_table = "some-dynamodb-table"
  }
}
```

This will use a S3 bucket as a remote backend, which will store and encrypt the terraform state while also providing state locking through the DynamoDB table. This is a partial configuration example, so a backend config file containing the AWS `access_key` and `secret_key` is needed. Thus you have to initialize terraform with the following command:

```
terraform init -backend-config=PATH
```

Where `PATH` is the path to the partial configuration file.

## See also

- [Remote Backends](https://www.terraform.io/docs/language/settings/backends/remote.html)
- [S3 remote backend](https://www.terraform.io/docs/language/settings/backends/s3.html)
- [Partial Configuration](https://www.terraform.io/docs/language/settings/backends/configuration.html#partial-configuration)
