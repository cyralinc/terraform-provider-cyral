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
  name        = "read-limit-policy"
  category    = "SECURITY"
  description = "Limits to 100 the amount of rows that can be read per query on all repository data for group 'Devs'"
  template_id = "read-limit"
  parameters  = "{ \"rowLimit\": 100, \"block\": true, \"alertSeverity\": \"high\", \"appliesToAllData\": true, \"identities\": { \"included\": { \"groups\": [\"Devs\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
}
