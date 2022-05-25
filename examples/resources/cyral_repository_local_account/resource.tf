### AWS IAM
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    aws_iam {
        database_name = ""
        local_account = ""
        role_arn = ""
    }
}

### AWS Secrets Manager
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    aws_secrets_manager {
        database_name = ""
        local_account = ""
        secret_arn = ""
    }
}

### Cyral Storage
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    cyral_storage {
        database_name = ""
        local_account = ""
        password = ""
    }
}

### Hashicorp Vault
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    hashicorp_vault {
        database_name = ""
        local_account = ""
        path = ""
    }
}

### Environment variable
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    environment_variable {
        local_account = ""
        variable_name = ""
    }
}

### GCP Secret Manager
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    gcp_secret_manager {
        database_name = ""
        local_account = ""
        secret_name = ""
    }
}
