resource "cyral_datalabel" "NAME" {
  name        = "NAME"
  description = "Customer name"
  tags        = ["PII", "SENSITIVE"]
  classification_rule {
    rule_type = "REGO"
    rule_code = "some-rego-code"
    rule_status = "ENABLED"
  }
}
