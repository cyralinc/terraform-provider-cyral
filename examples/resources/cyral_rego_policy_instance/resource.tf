### Global rego policy instance
resource "cyral_rego_policy_instance" "policy" {
  name = "User Management"
  category = "SECURITY"
  description = "Policy to govern user management operations"
  template_id = "object-protection"
  parameters = jsonencode(
    {
      "objectType" = "role/user"
      "block" = true
      "monitorCreates" = true
      "monitorAlters" = true
      "monitorDrops" = true
      "identities" = {
        "excluded" = {
          "groups" = ["dba"]
        }
      }
    }
  )
  enabled = true
  tags = ["tag1", "tag2"]
}

output "policy_last_updated" {
  value = cyral_rego_policy_instance.policy.last_updated
}

output "policy_created" {
  value = cyral_rego_policy_instance.policy.created
}

### Repo-level policy
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
      host = "mysql.cyral.com"
      port = 3306
  }
}

resource "cyral_rego_policy_instance" "policy" {
  name = "User Management"
  category = "SECURITY"
  description = "Policy to govern user management operations"
  template_id = "object-protection"
  parameters = jsonencode(
    {
      "objectType" = "role/user"
      "block" = true
      "monitorCreates" = true
      "monitorAlters" = true
      "monitorDrops" = true
      "identities" = {
        "excluded" = {
          "groups" = ["dba"]
        }
      }
    }
  )
  enabled = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}

### Rego policy instance with duration
resource "cyral_rego_policy_instance" "policy" {
  name = "User Management"
  category = "SECURITY"
  description = "Policy to govern user management operations"
  template_id = "object-protection"
  parameters = jsonencode(
    {
      "objectType" = "role/user"
      "block" = true
      "monitorCreates" = true
      "monitorAlters" = true
      "monitorDrops" = true
      "identities" = {
        "excluded" = {
          "groups" = ["dba"]
        }
      }
    }
  )
  enabled = true
  tags = ["tag1", "tag2"]
  duration = "10s"
}
