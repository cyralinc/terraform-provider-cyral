# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy set using the data protection wizard to raise
# an alert and block updates and deletes on label CCN
resource "cyral_policy_set" "data_protection_policy" {
  name        = "data protection policy"
  description = "Raise an alert and block updates and deletes on label CCN"
  wizard_id   = "data-protection"
  parameters  = jsonencode(
    {
      "block" = true
      "alertSeverity" = "high"
      "governedOperations" = ["update", "delete"]
      "labels" = ["CCN"]
    }
  )
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
