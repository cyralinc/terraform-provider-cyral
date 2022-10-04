### AWS Secrets Manager with automatic approval
resource "cyral_repository_local_account" "some_resource_name" {
    repository_id = cyral_repository.SOME_REPOSITORY_RESOURCE_NAME.id
    config {
        auto_approve_access = true
        # Automatically approve 5 minutes access requests
        max_auto_approve_duration = "PT5M"
    }
    aws_secrets_manager {
        database_name = ""
        local_account = ""
        secret_arn = ""
    }
}
