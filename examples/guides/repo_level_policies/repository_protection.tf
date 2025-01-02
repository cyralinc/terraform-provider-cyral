# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy set using the repository protection wizard to alert if more than
# 100 rows are updated or deleted per query on all repository data by anyone except group 'Admin'
resource "cyral_policy_set" "repository_protection_policy" {
  name        = "repository protection policy"
  description = "Alert if more than 100 rows are updated or deleted per query on all repository data by anyone except group 'Admin'"
  wizard_id   = "repository-protection"
  parameters  = jsonencode(
    {
      "rowLimit" = 100
      "datasets" = "*"
      "governedOperations" = ["update", "delete"]
      "identities" = { "excluded": { "groups" = ["Admin"] } }
    }
  )
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
}
