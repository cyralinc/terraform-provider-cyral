---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cyral_role Resource - terraform-provider-cyral"
subcategory: ""
description: |-
    Manages roles for Cyral control plane users https://cyral.com/docs/user-administration/manage-cyral-roles/#create-and-manage-administrator-roles-for-cyral-control-plane-users. See also: Role SSO Groups ./role_sso_groups.md.
---

# cyral_role (Resource)

Manages [roles for Cyral control plane users](https://cyral.com/docs/user-administration/manage-cyral-roles/#create-and-manage-administrator-roles-for-cyral-control-plane-users). See also: [Role SSO Groups](./role_sso_groups.md).

## Example Usage

```terraform
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
    view_audit_logs = false
    modify_integrations = false
    modify_roles = false
    view_datamaps = false
    repo_crawler = false
    approval_management = false
  }
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

-   `name` (String) The name of the role.

### Optional

-   `permissions` (Block Set, Max: 1) A block responsible for configuring the role permissions. (see [below for nested schema](#nestedblock--permissions))

### Read-Only

-   `id` (String) ID of this resource in Cyral environment

<a id="nestedblock--permissions"></a>

### Nested Schema for `permissions`

Optional:

-   `approval_management` (Boolean) Allows approving or denying approval requests on Cyral Control Plane. Defaults to `false`.
-   `modify_integrations` (Boolean) Allows modifying integrations on Cyral Control Plane. Defaults to `false`.
-   `modify_policies` (Boolean) Allows modifying policies on Cyral Control Plane. Defaults to `false`.
-   `modify_roles` (Boolean) Allows modifying roles on Cyral Control Plane. Defaults to `false`.
-   `modify_sidecars_and_repositories` (Boolean) Allows modifying sidecars and repositories on Cyral Control Plane. Defaults to `false`.
-   `modify_users` (Boolean) Allows modifying users on Cyral Control Plane. Defaults to `false`.
-   `repo_crawler` (Boolean) Allows running the Cyral repo crawler data classifier and user discovery. Defaults to `false`.
-   `view_audit_logs` (Boolean) Allows viewing audit logs on Cyral Control Plane. Defaults to `false`.
-   `view_datamaps` (Boolean) Allows viewing datamaps on Cyral Control Plane. Defaults to `false`.
-   `view_integrations` (Boolean) Allows viewing integrations on Cyral Control Plane. Defaults to `false`.
-   `view_policies` (Boolean) Allows viewing policies on Cyral Control Plane. Defaults to `false`.
-   `view_roles` (Boolean) Allows viewing roles on Cyral Control Plane. Defaults to `false`.
-   `view_users` (Boolean) Allows viewing users on Cyral Control Plane. Defaults to `false`.
