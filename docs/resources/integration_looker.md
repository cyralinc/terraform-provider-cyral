# Looker Integration Resource

Provides integration with Looker.

## Example Usage

```hcl
resource "cyral_integration_looker" "some_resource_name" {
    client_id = ""
    client_secret = ""
    url = ""
}
```

## Argument Reference

* `client_id` - (Required) Looker client id.
* `client_secret` - (Required) Looker client secret.
* `url` - (Required) Looker integration url.

## Attribute Reference

* `id` - The ID of this resource.
