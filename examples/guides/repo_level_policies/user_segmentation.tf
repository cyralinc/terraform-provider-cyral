# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy instance from template to filter table
# 'finance.cards' when users in group 'Marketing' read label
# CCN, returning only data where finance.cards.country = 'US'
resource "cyral_rego_policy_instance" "policy" {
  name        = "user-segmentation-policy"
  category    = "SECURITY"
  description = "Filter table 'finance.cards' when users in group 'Marketing' read label CCN, returning only data where finance.cards.country = 'US'"
  template_id = "user-segmentation"
  parameters  = "{ \"dataSet\": \"finance.cards\", \"dataFilter\": \" finance.cards.country = 'US' \", \"labels\": [\"CCN\"], \"includedIdentities\": { \"groups\": [\"Marketing\"] } }"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
