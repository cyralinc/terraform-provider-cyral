resource "cyral_policy_v2" "local_policy_example" {
  name        = "local_policy"
  description = "Local policy to allow gym users to read their own data"
  enabled     = true
  tags        = ["gym", "local"]
  scope {
    repo_ids = ["2gaWEAyeKbETyUy1LSx985gVqrk"]
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
    repo_ids = ["2gaWEAyeKbETyUy1LSx985gVqrk"]
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


resource "cyral_policy_v2" "approval_policy_example" {
  name        = "approval_policy"
  description = "Approval policy for emergency access"
  enabled     = true
  tags        = ["approval", "emergency"]
  scope {
    repo_ids = ["gym_repo"]
  }
  valid_from  = "2024-06-12T00:00:00Z"
  valid_until = "2024-06-13T00:00:00Z"
  document    = jsonencode({
    identity = {
      email    = "gym_owner@example.com"
    }
    conditions = [
      {
        attribute = "role"
        operator  = "equals"
        value     = "admin"
      }
    ]
    approvals = {
      read = {
        tags = ["PII"]
        locations = ["gym_db.users"]

      }
      update = {
        tags = ["PII"]
        locations = ["gym_db.users"]
      }
      delete = {
        tags = ["PII"]
        locations = ["gym_db.users"]
      }
    }
  })
  type        = "approval"
}
