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
# 'finance.cards', returning only data where
# finance.cards.country = 'US' for users not in 'Admin' group
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-firewall-policy"
  category    = "SECURITY"
  description = "Returns only data where finance.cards.country = 'US' in table 'finance.cards' for users not in 'Admin' group"
  template_id = "data-firewall"
  parameters  = "{ \"dataSet\": \"finance.cards\", \"dataFilter\": \" finance.cards.country = 'US' \", \"labels\": [\"CCN\"], \"excludedIdentities\": { \"groups\": [\"Admin\"] } }"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
