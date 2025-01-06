# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy set using the user segmentation wizard to filter table
# 'finance.cards' when users in group 'Marketing' read label
# CCN, returning only data where finance.cards.country = 'US'
resource "cyral_policy_set" "user_segmentation_policy" {
  name        = "user segmentation policy"
  description = "Filter table 'finance.cards' when users in group 'Marketing' read label CCN, returning only data where finance.cards.country = 'US'"
  wizard_id   = "user-segmentation"
  wizard_parameters  = jsonencode(
    {
      "dataset" = "finance.cards"
      "dataFilter" = " finance.cards.country = 'US' "
      "labels" = ["CCN"]
      "includedIdentities" = { "groups" = ["Marketing"] }
    }
  )
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
