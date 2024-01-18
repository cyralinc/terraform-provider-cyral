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
  name        = "data-masking-policy"
  category    = "SECURITY"
  description = "Masks label CCN for identities in Marketing group"
  template_id = "data-masking"
  parameters  = "{ \"maskType\": \"NULL_MASK\", \"labels\": [\"CCN\"], \"identities\": { \"included\": { \"groups\": [\"Marketing\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
