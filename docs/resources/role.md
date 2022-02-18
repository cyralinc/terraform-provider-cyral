# Role Resource

Provides a resource to manage [roles for Cyral control plane users](https://cyral.com/docs/account-administration/acct-manage-cyral-roles/#create-and-manage-administrator-roles-for-cyral-control-plane-users). See also: [Role SSO Groups](./role_sso_groups.md).

## Example Usage

### Role with Default Configuration

```hcl
resource "cyral_role" "some_resource_name" {
  name="some-role-name"
}
```

### Role with Custom Permissions Configuration

```hcl
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
```

## Argument Reference

- `name` - (Required) The name of the role.
- `permissions` - (Optional) A block responsible for configuring the role permissions. See [permissions](#permissions) below for more details.

### permissions

The `permissions` object supports the following:

- `modify_sidecars_and_repositories` - (Optional) Allows modifying sidecars and repositories for this role. Defaults to `false`.
- `modify_users` - (Optional) Allows modifying users for this role. Defaults to `false`.
- `modify_policies` - (Optional) Allows modifying policies for this role. Defaults to `false`.
- `view_sidecars_and_repositories` - (Optional) Allows viewing sidecars and repositories for this role. Defaults to `false`.
- `view_audit_logs` - (Optional) Allows viewing audit logs for this role. Defaults to `false`.
- `modify_integrations` - (Optional) Allows modifying integrations for this role. Defaults to `false`.
- `modify_roles` - (Optional) Allows modifying roles for this role. Defaults to `false`.
- `view_datamaps` - (Optional) Allows viewing datamaps for this role. Defaults to `false`.

## Attribute Reference

- `id` - The ID of this resource.
