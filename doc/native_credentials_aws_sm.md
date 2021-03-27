# Deploy Native Repository Credentials to AWS Secrets Manager

Use the following code as an example to deploy a native repository credential that will be used by 

```hcl
# See the Cyral provider documentation for more
# information on how to initialize it correctly.
provider "cyral" {
    control_plane = "mycontrolplane.cyral.com:8000"
}

resource "cyral_repository" "mongodb_repo" {
    type = "mongodb"
    host = "mongodb.cyral.com"
    port = 27017
    name = "mymongodb"
}

# See the AWS provider documentation for more
# information on how to initialize it correctly.
provider "aws" {
    # By deploying the secret to the same account and region of your
    # sidecar and using the name suggested in my_repository_secret, 
    # the sidecar will gain access to the secret automatically.
    region = "us-east-1"
}

resource "aws_secretsmanager_secret" "my_repository_secret" {
    name = join("", [
      "/cyral/dbsecrets/",
      cyral_repository.mongodb_repo.id
    ])
}

resource "aws_secretsmanager_secret_version" "my_repository_secret_version" {
    secret_id     = aws_secretsmanager_secret.my_repository_secret.id
    secret_string = jsonencode(local.repository_credentials)
}
```
