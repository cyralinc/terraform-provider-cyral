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
  name        = "service account abuse policy"
  category    = "SECURITY"
  description = "Always require user attribution for service acount 'john'"
  template_id = "service-account-abuse"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"serviceAccounts\": [\"john\"]}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
}
