# Creates MySQL data repository
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-firewall-policy"
  category    = "SECURITY"
  description = "Filter 'finance.cards' when someone (except 'Admin' group) reads it"
  template_id = "data-firewall"
  parameters  = "{ \"dataSet\": \"finance.cards\", \"dataFilter\": \" finance.cards.country = 'US' \", \"labels\": [\"CCN\"], \"excludedIdentities\": { \"groups\": [\"Admin\"] } }"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
