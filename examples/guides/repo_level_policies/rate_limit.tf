# Creates pg data repository
resource "cyral_repository" "pg1" {
  type = "postgresql"
  name = "pg-1"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# Creates a policy instance from template to raise a 'high' alert
# and set a rate limit of 500 rows per hour for group 'Marketing'
# and any data labeled as CCN
resource "cyral_rego_policy_instance" "policy" {
  name        = "rate-limit-policy"
  category    = "SECURITY"
  description = "Raise a 'high' alert and set a rate limit of 500 rows per hour for group 'Marketing' and any data labeled as CCN"
  template_id = "rate-limit"
  parameters  = "{ \"rateLimit\": 500, \"block\": true, \"alertSeverity\": \"high\", \"labels\": [\"CCN\"], \"identities\": { \"included\": { \"groups\": [\"Marketing\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
  tags = ["tag1", "tag2"]
}
