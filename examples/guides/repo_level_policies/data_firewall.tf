# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy set using the data firewall wizard to filter table
# 'finance.cards', returning only data where
# finance.cards.country = 'US' for users not in 'Admin' group
resource "cyral_policy_set" "data_firewall_policy" {
  name        = "data firewall policy"
  description = "Returns only data where finance.cards.country = 'US' in table 'finance.cards' for users not in 'Admin' group"
  wizard_id   = "data-firewall"
  wizard_parameters  = jsonencode(
    {
      "dataset" = "finance.cards"
      "dataFilter" = " finance.cards.country = 'US' "
      "labels" = ["CCN"]
      "excludedIdentities" = { "groups" = ["Admin"] }
    }
  )
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
