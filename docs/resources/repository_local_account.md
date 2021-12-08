# Repository Local Account Resource

Provides a resource to handle repository local accounts.

## Example Usage

### AWS IAM

```hcl
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    aws_iam {
        database_name = ""
        local_account = ""
        role_arn = ""
    }
}
```

### AWS Secrets Manager

```hcl
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    aws_secrets_manager {
        database_name = ""
        local_account = ""
        secret_arn = ""
    }
}
```

### Cyral Storage

```hcl
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    cyral_storage {
        database_name = ""
        local_account = ""
        password = ""
    }
}
```

### Hashicorp Vault

```hcl
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    hashicorp_vault {
        database_name = ""
        local_account = ""
        path = ""
    }
}
```

### Environment variable

```hcl
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    environment_variable {
        local_account = ""
        variable_name = ""
    }
}
```

## Argument Reference

* `repository_id` - (Required) ID of the repository that will be used by the local account.
* `aws_iam` - (Optional) Credential option to set the local account from AWS IAM. See [aws_iam](#aws_iam) below for more details.
* `aws_secrets_manager` - (Optional) Credential option to set the local account from AWS Secrets Manager. See [aws_secrets_manager](#aws_secrets_manager) below for more details.
* `cyral_storage` - (Optional) Credential option to set the local account from Cyral Storage. See [cyral_storage](#cyral_storage) below for more details.
* `hashicorp_vault` - (Optional) Credential option to set the local account from Hashicorp Vault. See [hashicorp_vault](#hashicorp_vault) below for more details.

### aws_iam

The `aws_iam` object supports the following arguments:

* `database_name` - (Optional) Database name that the local account corresponds to.
* `local_account` - (Required) Local account name.
* `role_arn` - (Required) AWS IAM role ARN.

### aws_secrets_manager

The `aws_secrets_manager` object supports the following arguments:

* `database_name` - (Optional) Database name that the local account corresponds to.
* `local_account` - (Required) Local account name.
* `secret_arn` - (Required) ARN of the AWS Secret Manager that stores the credential.

### cyral_storage

The `cyral_storage` object supports the following arguments:

* `database_name` - (Optional) Database name that the local account corresponds to.
* `local_account` - (Required) Local account name.
* `password` - (Required) Local account password.

### hashicorp_vault

The `hashicorp_vault` object supports the following arguments:

* `database_name` - (Optional) Database name that the local account corresponds to.
* `local_account` - (Required) Local account name.
* `path` - (Required) Hashicorp Vault path.

### environment_variable

The `environment_variable` object supports the following arguments:

* `database_name` - (Optional) Database name that the local account corresponds to.
* `local_account` - (Required) Local account name.
* `environment_name` - (Required) Name of the environment variable that will store credentials.

## Attribute Reference

* `id` - The ID of this resource.