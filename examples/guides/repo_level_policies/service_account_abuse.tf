# Creates pg data repository
resource "cyral_repository" "pg1" {
  type = "postgresql"
  name = "pg-1"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# Creates a policy set using the service account abuse wizard to alert and block
# whenever the service accounts john is used without end user attribution.
resource "cyral_policy_set" "service_account_abuse_policy" {
  name        = "service account abuse policy"
  description = "Alert and block whenever the service accounts john is used without end user attribution"
  wizard_id   = "service-account-abuse"
  wizard_parameters  = jsonencode(
    {
      "block" = true
      "alertSeverity" = "high"
      "serviceAccounts" = ["john"]
    }
  )
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
}
