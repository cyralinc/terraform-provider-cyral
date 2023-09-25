### Service account with all permissions
data "cyral_permission" "this" {}

resource "cyral_service_account" "this" {
  display_name = "cyral-service-account-1"
  permission_ids = [
    for p in data.cyral_permission.this.permission_list: p.id
  ]
}

output "client_id" {
  value = cyral_service_account.this.client_id
}

output "client_secret" {
  sensitive = true
  value = cyral_service_account.this.client_secret
}

### Service account with specific permissions
data "cyral_permission" "this" {}

locals {
  saPermissions = [
		"Modify Policies",
		"Modify Integrations",
  ]
}

resource "cyral_service_account" "this" {
  display_name = "cyral-service-account-1"
  permission_ids = [
    for p in data.cyral_permission.this.permission_list: p.id
    if contains(local.saPermissions, p.name)
  ]
}
