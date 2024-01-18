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
  name        = "user-segmentation-policy"
  category    = "SECURITY"
  description = "Applies a data filter in 'finance.cards' when someone from group 'Marketing' reads data labeled as 'CCN'"
  template_id = "user-sementation"
  parameters  = "{ \"dataSet\": \"finance.cards\", \"dataFilter\": \" finance.cards.country = 'US' \", \"labels\": [\"CCN\"], \"includedIdentities\": { \"groups\": [\"Marketing\"] } }"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
