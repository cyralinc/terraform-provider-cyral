# Creates MySQL data repository
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-protection-policy"
  category    = "SECURITY"
  description = "Protect label CCN for update and delete queries"
  template_id = "data-protection"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorDeletes\": true, \"labels\": [\"CCN\"]}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
