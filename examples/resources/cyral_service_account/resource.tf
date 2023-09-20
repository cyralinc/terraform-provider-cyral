### Service account with all permissions
resource "cyral_service_account" "sa_1" {
  display_name = "cyral-service-account-1"
  permissions {
    modify_sidecars_and_repositories = true
    modify_policies = true
    modify_integrations = true
    modify_users = true
    modify_roles = true
    view_users = true
    view_audit_logs = true
    repo_crawler = true
    view_datamaps = true
    view_roles = true
    view_policies = true
    approval_management = true
    view_integrations = true
  }
}

output "client_id" {
  value = cyral_service_account.sa_1.client_id
}

output "client_secret" {
  sensitive = true
  value = cyral_service_account.sa_1.client_secret
}
