# Datamap

CRUD operations for Cyral Repository Local Account.

## Usage

Although, all authentication schemas are listed below, the API expects only one to be used.

```hcl
resource "cyral_repository_local_account" "my-repo-account" {
    aws_iam {
        database_name = ""
        local_account = ""
        role_arn = ""
    }
    aws_secrets_manager {
        database_name = ""
        local_account = ""
        secret_arn = ""
    }
    cyral_storage {
        database_name = ""
        local_account = ""
        password = ""
    }
    hashicorp_vault {
        database_name = ""
        local_account = ""
        path = ""
    }
}
```

## Variables

### aws_iam

|  Name           |  Default  |  Description                                         | Required |
|:----------------|:---------:|:-----------------------------------------------------|:--------:|
| `database_name` |           | Database name that the local account corresponds to. | Yes      |
| `local_account` |           | Local repository account name.                       | Yes      |
| `role_arn`      |           | AWS IAM role ARN.                                    | Yes      |

### aws_secrets_manager

|  Name           |  Default  |  Description                                         | Required |
|:----------------|:---------:|:-----------------------------------------------------|:--------:|
| `database_name` |           | Database name that the local account corresponds to. | Yes      |
| `local_account` |           | Local repository account name.                       | Yes      |
| `secret_arn`    |           | AWS Secret Manager ARN.                              | Yes      |



### cyral_storage

|  Name           |  Default  |  Description                                         | Required |
|:----------------|:---------:|:-----------------------------------------------------|:--------:|
| `database_name` |           | Database name that the local account corresponds to. | Yes      |
| `local_account` |           | Local repository account name.                       | Yes      |
| `password`      |           | Cyral Storage Password.                              | Yes      |

### hashicorp_vault

|  Name           |  Default  |  Description                                         | Required |
|:----------------|:---------:|:-----------------------------------------------------|:--------:|
| `database_name` |           | Database name that the local account corresponds to. | Yes      |
| `local_account` |           | Local repository account name.                       | Yes      |
| `path`          |           | Hashicorp Vault path.                                | Yes      |
