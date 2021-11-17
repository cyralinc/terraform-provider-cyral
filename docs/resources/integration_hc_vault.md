# Hashicorp Vault Integration Resource

Provides integration with Hashicorp Vault to store secrets.

## Example Usage

```hcl
resource "cyral_integration_hc_vault" "some_resource_name" {
  name = ""
  auth_type = ""
  server = ""
  auth_method = ""
}
```

## Argument Reference

* `name` - (Required) Integration name that will be used internally in Control Plane.
* `auth_type` - (Required) Authentication type for the integration.
* `server` - (Required) The server on which the vault service is running.
* `auth_method` - (Required) The authentication method for the integration.

## Attribute Reference

* `id` - The ID of this resource.
