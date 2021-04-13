# Datamap

CRUD operations for Cyral Repository Local Account.

## Usage

Although all credential options are listed below, the API expects only one to be used at a time.

```hcl
resource "cyral_repository_local_account" "my-repo-account" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
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

|  Name                |  Default  |  Description                                                         | Required |
|:---------------------|:---------:|:---------------------------------------------------------------------|:--------:|
| `repository_id`      |           | ID of the repository that will be used by this account.              | Yes      |
| `aws_iam`            |           | Credential option to set the local account from AWS IAM.             | No       |
| `aws_secrets_manager`|           | Credential option to set the local account from AWS Secrets Manager. | No       |
| `cyral_storage`      |           | Credential option to set the local account from Cyral Storage.       | No       |
| `hashicorp_vault`    |           | Credential option to set the local account from Hashicorp Vault.     | No       |

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
