### Global rego policy instance
resource "cyral_rego_policy_instance" "policy" {
  name = "some-rate-limit-policy"
  category = "SECURITY"
  description = "Some policy description."
  template_id = "rate-limit"
  parameters = "{\"rateLimit\":7,\"labels\":[\"EMAIL\"],\"alertSeverity\":\"high\",\"block\":false}"
  enabled = true
  tags = ["tag1", "tag2"]
}

output "policy_last_updated" {
  value = cyral_rego_policy_instance.policy.last_updated
}

output "policy_created" {
  value = cyral_rego_policy_instance.policy.created
}

### Repo-level policy
resource "cyral_rego_policy_instance" "policy" {
  name = "some-data-masking-policy"
  category = "SECURITY"
  description = "Some policy description."
  template_id = "data-masking"
  parameters = "{\"labels\":[\"ADDRESS\"],\"maskType\":\"NULL_MASK\"}"
  enabled = true
  scope {
    repo_ids = ["2U4prk5o6yi1rTvvXyImz8lgbgG"]
  }
  tags = ["tag1", "tag2"]
}

### Rego policy instance with duration
resource "cyral_rego_policy_instance" "policy" {
  name = "some-data-masking-policy"
  category = "SECURITY"
  description = "Some policy description."
  template_id = "data-masking"
  parameters = "{\"labels\":[\"ADDRESS\"],\"maskType\":\"NULL_MASK\"}"
  enabled = true
  tags = ["tag1", "tag2"]
  duration = "10s"
}
