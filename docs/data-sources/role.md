---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cyral_role Data Source - cyral"
subcategory: ""
description: |-
  Retrieve and filter roles https://cyral.com/docs/account-administration/acct-manage-cyral-roles/ that exist in the Cyral Control Plane.
---

# cyral_role (Data Source)

Retrieve and filter [roles](https://cyral.com/docs/account-administration/acct-manage-cyral-roles/) that exist in the Cyral Control Plane.

## Example Usage

```terraform
data "cyral_role" "admin_roles" {
  # Optional. Filter roles with name that matches regular expression.
  name = "^.*Admin$"
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Optional

- `name` (String) Filter the results by a regular expression (regex) that matches names of existing roles.

### Read-Only

- `id` (String) The ID of this resource.
- `role_list` (List of Object) List of existing roles satisfying given filter criteria. (see [below for nested schema](#nestedatt--role_list))

<a id="nestedatt--role_list"></a>

### Nested Schema for `role_list`

Read-Only:

- `description` (String)
- `id` (String)
- `members` (List of String)
- `name` (String)
- `roles` (List of String)
- `sso_groups` (List of Object) (see [below for nested schema](#nestedobjatt--role_list--sso_groups))

<a id="nestedobjatt--role_list--sso_groups"></a>

### Nested Schema for `role_list.sso_groups`

Read-Only:

- `group_name` (String)
- `id` (String)
- `idp_id` (String)
- `idp_name` (String)