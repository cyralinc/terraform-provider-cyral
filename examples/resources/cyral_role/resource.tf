### Role with Default Configuration
resource "cyral_role" "some_resource_name" {
  name="some-role-name"
}

### Role with Custom Permissions Configuration
resource "cyral_role" "some_resource_name" {
  name="some-role-name"
  permissions {
    modify_sidecars_and_repositories = true
    modify_users = true
    modify_policies = true
    view_sidecars_and_repositories = true
    view_audit_logs = false
    modify_integrations = false
    modify_roles = false
    view_datamaps = false
  }
}
