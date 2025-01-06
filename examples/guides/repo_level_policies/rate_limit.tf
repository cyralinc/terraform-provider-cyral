# Creates pg data repository
resource "cyral_repository" "pg1" {
  type = "postgresql"
  name = "pg-1"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# Creates a policy set using the rate limit wizard to raise an alert
# and set a rate limit of 500 rows per hour for group 'Marketing'
# and any data labeled as CCN
resource "cyral_policy_set" "rate_limit_policy" {
  name        = "rate limit policy"
  description = "Raise an alert and set a rate limit of 500 rows per hour for group 'Marketing' and any data labeled as CCN"
  wizard_id   = "rate-limit"
  wizard_parameters  = jsonencode(
    {
      "rateLimit" = 500
      "enforce" = true
      "labels" = ["CCN"]
      "identities" = { "included": { "groups" = ["Marketing"] } }
    }
  )
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
  tags = ["tag1", "tag2"]
}
