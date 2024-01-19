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
# and block updates and reads on schema 'finance' and dataset
# 'cyral.customers'
resource "cyral_rego_policy_instance" "policy" {
  name        = "dataset-protection"
  category    = "SECURITY"
  description = "Raise a 'high' alert and block updates and reads on schema 'finance' and dataset 'cyral.customers'"
  template_id = "dataset-protection"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorReads\": true, \"datasets\": {\"disallowed\": [\"finance.*\", \"cyral.customers\"]}}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
}
