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
  name        = "dataset-protection"
  category    = "SECURITY"
  description = "Blocks reads and updates over schema 'finance' and dataset 'cyral.customers'."
  template_id = "dataset-protection"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorReads\": true, \"datasets\": {\"disallowed\": [\"finance.*\", \"cyral.customers\"]}}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
}
