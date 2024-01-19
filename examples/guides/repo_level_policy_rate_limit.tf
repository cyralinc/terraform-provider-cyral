# Creates pg data repository
resource "cyral_repository" "repo" {
  type = "postgresql"
  name = "my_pg"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "rate-limit-policy"
  category    = "SECURITY"
  description = "Implement a threshold on label CCN for group Marketing of 500 rows per hour"
  template_id = "rate-limit"
  parameters  = "{ \"rateLimit\": 500, \"block\": true, \"alertSeverity\": \"high\", \"labels\": [\"CCN\"], \"identities\": { \"included\": { \"groups\": [\"Marketing\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
