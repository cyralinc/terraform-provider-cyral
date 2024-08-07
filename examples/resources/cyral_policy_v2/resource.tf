resource "cyral_repository" "myrepo" {
    type = "mongodb"
    name = "myrepo"

    repo_node {
        name = "node-1"
        host = "mongodb.cyral.com"
        port = 27017
    }

    mongodb_settings {
      server_type = "standalone"
    }
}

resource "cyral_policy_v2" "local_policy_example" {
  name        = "local_policy"
  description = "Local policy to allow gym users to read their own data"
  enabled     = true
  tags        = ["gym", "local"]
  scope {
    repo_ids = [cyral_repository.myrepo.id]
  }
  document    = jsonencode({
    governedData = {
      locations = ["gym_db.users"]
    }
    readRules = [
      {
        conditions = [
          {
            attribute = "identity.userGroups"
            operator  = "contains"
            value     = "users"
          }
        ]
        constraints = {
          datasetRewrite = "SELECT * FROM $${dataset} WHERE email = '$${identity.endUserEmail}'"
        }
      }
    ]
  })
  enforced    = true
  type        = "local"
}

resource "cyral_policy_v2" "global_policy_example" {
  name        = "global_policy"
  description = "Global policy for finance users with row limit for PII data"
  enabled     = true
  tags        = ["finance", "global"]
  scope {
    repo_ids = [cyral_repository.myrepo.id]
  }
  document    = jsonencode({
    governedData = {
      tags = ["PII"]
    }
    readRules = [
      {
        conditions = [
          {
            attribute = "identity.userGroups"
            operator  = "contains"
            value     = "finance"
          }
        ]
        constraints = {
          maxRows = 5
        }
      },
      {
        conditions = []
        constraints = {}
      }
    ]
  })
  enforced    = true
  type        = "global"
}
