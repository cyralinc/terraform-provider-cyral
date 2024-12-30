# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy set using the data masking wizard to apply null masking to
# any data labeled as CCN for users in group 'Marketing'
resource "cyral_policy_set" "data_masking_policy" {
  name        = "data masking policy"
  description = "Apply null masking to any data labeled as CCN for users in group 'Marketing'"
  wizard_id   = "data-masking"
  parameters  = jsonencode(
    {
      "maskType" = "null"
      "labels" = ["CCN"]
      "identities" = { "included": { "groups" = ["Marketing"] } }
    }
  )
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
