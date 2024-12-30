# Creates pg data repository
resource "cyral_repository" "pg1" {
  type = "postgresql"
  name = "pg-1"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# Creates a policy set using the read limit wizard to limits to 100 the
# amount of rows that can be read per query on the entire
# repository for group 'Devs'
resource "cyral_policy_set" "read_limit_policy" {
  name        = "read limit policy"
  description = "Limits to 100 the amount of rows that can be read per query on the entire repository for group 'Devs'"
  wizard_id   = "read-limit"
  parameters  = jsonencode(
    {
      "rowLimit" = 100
      "enforce" = true
      "datasets" = "*"
      "identities" = { "included": { "groups" = ["Devs"] } }
    }
  )
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
}
