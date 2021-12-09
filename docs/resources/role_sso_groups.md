# Role SSO Groups Resource

Provides a resource to [map SSO groups to a specific role](https://cyral.com/docs/account-administration/acct-manage-cyral-roles/#map-an-sso-group-to-a-cyral-administrator-role) on Cyral control plane. See also: [Role](./role.md).

## Example Usage

### Role with Single SSO Group

```hcl
resource "cyral_integration_idp_okta" "some_idp_okta" {
  samlp {
    config {
      single_sign_on_service_url = "https://some-sso-url.com"
    }
  }
}

resource "cyral_role" "some_role" {
  name="some-role-name"
}

resource "cyral_role_sso_groups" "some_role_sso_groups" {
  role_id=cyral_role.some_role.id
  sso_group {
    group_name="Everyone"
    idp_id=cyral_integration_idp_okta.some_idp_okta.id
  }
}
```

### Role with Multiple SSO Groups

```hcl
resource "cyral_integration_idp_okta" "some_idp_okta" {
  samlp {
    config {
      single_sign_on_service_url = "https://some-sso-url.com"
    }
  }
}

resource "cyral_role" "some_role" {
  name="some-role-name"
}

resource "cyral_role_sso_groups" "some_role_sso_groups" {
  role_id=cyral_role.some_role.id
  sso_group {
    group_name="Admin"
    idp_id=cyral_integration_idp_okta.some_idp_okta.id
  }
  sso_group {
    group_name="Dev"
    idp_id=cyral_integration_idp_okta.some_idp_okta.id
  }
}
```

## Argument Reference

* `role_id` - (Required) The ID of the role resource that will be configured.
* `sso_group` - (Required) A block responsible for mapping an SSO group to a role. See [sso_group](#sso_group) below for more details.

### sso_group

The `sso_group` object supports the following:

* `group_name` - (Required) The name of the SSO group to be mapped.
* `idp_id` - (Required) The ID of the identity provider integration to be mapped.


## Attribute Reference

* `id` - The ID of this resource.
* `sso_group.*.id` - The ID of an SSO group mapping.
* `sso_group.*.idp_name` - The name of the identity provider integration of an SSO group mapping.
