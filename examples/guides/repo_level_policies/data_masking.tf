# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy instance from template to apply null masking to
# any data labeled as CCN for users in group 'Marketing'
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-masking-policy"
  category    = "SECURITY"
  description = "Apply null masking to any data labeled as CCN for users in group 'Marketing'"
  template_id = "data-masking"
  parameters  = "{ \"maskType\": \"NULL_MASK\", \"labels\": [\"CCN\"], \"identities\": { \"included\": { \"groups\": [\"Marketing\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
