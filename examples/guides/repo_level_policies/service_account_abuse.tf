# Creates pg data repository
resource "cyral_repository" "pg1" {
  type = "postgresql"
  name = "pg-1"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# Creates a policy instance from template to alert and block
# whenever the following service accounts john try to read,
# update, or delete data from the repository without end
# user attribution.
resource "cyral_rego_policy_instance" "policy" {
  name        = "service account abuse policy"
  category    = "SECURITY"
  description = "Alert and block whenever the following service accounts john try to read, update, or delete data from the repository without end user attribution"
  template_id = "service-account-abuse"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"serviceAccounts\": [\"john\"]}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
}
