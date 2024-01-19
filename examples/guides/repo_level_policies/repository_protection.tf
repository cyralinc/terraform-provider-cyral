# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy instance from template to limits to 100 the
# amount of rows that can be updated or deleted per query on
# all repository data for anyone except group 'Admin'
resource "cyral_rego_policy_instance" "policy" {
  name        = "repository-protection-policy"
  category    = "SECURITY"
  description = "Limits to 100 the amount of rows that can be updated or deleted per query on all repository data for anyone except group 'Admin'"
  template_id = "repository-protection"
  parameters  = "{ \"rowLimit\": 100, \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorDeletes\": true, \"identities\": { \"excluded\": { \"groups\": [\"Admin\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
}
