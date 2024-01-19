# Creates pg data repository
resource "cyral_repository" "pg1" {
  type = "postgresql"
  name = "pg-1"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# Creates a policy instance from template to limits to 100 the
# amount of rows that can be read per query on the entire
# repository for group 'Devs'
resource "cyral_rego_policy_instance" "policy" {
  name        = "read-limit-policy"
  category    = "SECURITY"
  description = "Limits to 100 the amount of rows that can be read per query on the entire repository for group 'Devs'"
  template_id = "read-limit"
  parameters  = "{ \"rowLimit\": 100, \"block\": true, \"alertSeverity\": \"high\", \"appliesToAllData\": true, \"identities\": { \"included\": { \"groups\": [\"Devs\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
}
