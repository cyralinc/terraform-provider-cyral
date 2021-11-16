# Microsoft Teams Integration Resource

Provides integration with Microsoft Teams.

## Example Usage

```hcl
resource "cyral_integration_microsoft_teams" "some_resource_name" {
    name = ""
    url = ""
}
```

## Argument Reference

* `name` - (Required) Integration name that will be used internally in Control Plane.
* `url` - (Required) Microsoft Teams webhook URL.

## Attribute Reference

* `id` - The ID of this resource.
