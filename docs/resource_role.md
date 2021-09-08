# Role

CRUD operations for Cyral roles.

## Usage

```hcl
resource "cyral_role" "SOME_ROLE_NAME" {
    name="some-role-name"
    permissions {
      view_sidecars_and_repositories = false|true
      modify_sidecars_and_repositories = false|true
      modify_policies = false|true
      modify_integrations = false|true
      modify_users = false|true
      modify_roles = false|true
      view_audit_logs = false|true
    }
}
```

## Observation

This resource only supports a single block of `permissions`.

## Variables

| Name          |          Default          | Description                         | Required |
| :------------ | :-----------------------: | :---------------------------------- | :------: |
| `name`        |                           | Role name that will be created.     |   Yes    |
| `permissions` | `false` for all permissions | Block of permissions for this role. Grants a specific permission for this role if set as `true`, otherwise all permissions are set as `false` by default. |    No    |

## Outputs

| Name | Description                                     |
| :--- | :---------------------------------------------- |
| `id` | Unique ID of the resource in the Control Plane. |
