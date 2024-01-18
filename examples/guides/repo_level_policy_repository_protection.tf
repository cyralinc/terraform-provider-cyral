# Creates MySQL data repository
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
    host = "mysql.cyral.com"
    port = 5432
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "repository-protection-policy"
  category    = "SECURITY"
  description = "Limits to 100 the amount of rows that can be updated or deleted per query on all repository data for anyone except group 'Admin'"
  template_id = "repository-protection"
  parameters  = "{ \"rowLimit\": 100, \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorDeletes\": true, \"identities\": { \"excluded\": { \"groups\": [\"Admin\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
}
