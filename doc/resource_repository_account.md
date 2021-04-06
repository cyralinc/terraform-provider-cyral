# Datamap

CRUD operations for Cyral Repository Account.

## Usage
Although, all authentication schemas are listed below, the API expects only one to be used.
```hcl
resource "cyral_repository_account" "my-repo-account" {
    aws_iam {
        database_name = ""
        repo_account = ""
        role_arn = ""
    }
    aws_secrets_manager {
        database_name = ""
        repo_account = ""
        secret_arn = ""
    }
    cyral_storage {
        database_name = ""
        password = ""
        repo_account = ""
    }
    hashicorp_vault {
        database_name = ""
        path = ""
        repo_account = ""
    }
}
```

## Variables

|  Name           |  Default  |  Description                                                                         | Required |
|:----------------|:---------:|:-------------------------------------------------------------------------------------|:--------:|
| `database_name`       |           | Repo Account Database Name.                    | No      |
| `repo_account`       |           | Repo Account Name.                    | No      |
| `role_arn`       |           | AWS IAM role ARN.                    | Yes      |
| `secret_arn`       |           | AWS Secret Manager ARN.                    | Yes      |
| `password`       |           | Cyral Storage Password.                    | Yes      |
| `path`       |           | Hashicorp Vault path.                    | Yes      |