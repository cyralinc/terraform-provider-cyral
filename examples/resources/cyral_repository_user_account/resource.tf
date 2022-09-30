# Test Repo to use in examples.
resource "cyral_repository" "tf_test_repo" {
    host = "postgresql.mycompany.com"
    name = "tf_test_repo"
    port = 5432
    type = "postgresql"
}


# cyral_repository_user_account with auth scheme aws_iam
resource "cyral_repository_user_account" "aws_iam" {
    name          = "hbf_aws_iam"
    repository_id = cyral_repository.tf_test_repo.id

    auth_scheme {
        aws_iam {
            role_arn = "role_arn"
        }
    }
}

# cyral_repository_user_account with auth scheme aws_secrets will be created
resource "cyral_repository_user_account" "aws_secrets" {

    name          = "hbf_aws_secrets"
    repository_id = cyral_repository.tf_test_repo.id

    auth_scheme {
        aws_secrets_manager {
            secret_arn = "secret_arn"
        }
    }
}

# cyral_repository_user_account with auth scheme env_var will be created
resource "cyral_repository_user_account" "env_var" {

    name          = "hbf_env_var"
    repository_id = cyral_repository.tf_test_repo.id

    auth_scheme {

        environment_variable {
            variable_name = "some-env-var"
        }
    }
}

# cyral_repository_user_account with auth scheme gcp_secrets will be created
resource "cyral_repository_user_account" "gcp_secrets" {

    name          = "hbf_gcp_secrets"
    repository_id = cyral_repository.tf_test_repo.id

    auth_scheme {

        gcp_secrets_manager {
            secret_name = "secret_name"
        }
    }
}

# cyral_repository_user_account with auth scheme hashicorp will be created
resource "cyral_repository_user_account" "hashicorp" {

    name          = "hbf_hashicorp"
    repository_id = cyral_repository.tf_test_repo.id

    auth_scheme {
        hashicorp_vault {
            path = "some-path"
            is_dynamic_user_account = false
        }
    }
}

# cyral_repository_user_account with auth scheme kubernetes will be created
resource "cyral_repository_user_account" "kubernetes" {

    name          = "hbf_kubernetes"
    repository_id = cyral_repository.tf_test_repo.id

    auth_scheme {

        kubernetes_secret {
            secret_key  = "secret_key"
            secret_name = "secret_name"
        }
    }
}
