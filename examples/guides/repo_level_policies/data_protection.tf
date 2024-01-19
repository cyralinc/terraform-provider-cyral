# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy instance from template to raise a 'high' alert
# and block updates and deletes on label CCN
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-protection-policy"
  category    = "SECURITY"
  description = "Raise a 'high' alert and block updates and deletes on label CCN"
  template_id = "data-protection"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorDeletes\": true, \"labels\": [\"CCN\"]}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
