# An example of a policy and a policy rule with a rego policy
# in `additional_checks`.
resource "cyral_policy" "this" {
  name = "My first policy"
  description = "This is my first policy"
  enabled = true
  data = ["EMAIL"]
  metadata_tags = ["Risk Level 1"]
}

resource "cyral_policy_rule" "this" {
  policy_id = cyral_policy.this.id
  deletes {
    additional_checks = <<EOT
is_valid_request {
  filter := request.filters[_]
  filter.field == "entity.user.is_real"
  filter.op == "="
  filter.value == false
}
EOT
    data = ["EMAIL"]
    rows = -1
    severity = "low"
  }
}
