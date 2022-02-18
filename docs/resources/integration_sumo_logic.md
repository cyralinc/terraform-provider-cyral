# Sumo Logic Integration Resource

Provides integration with [Sumo Logic to push sidecar logs](https://cyral.com/docs/integrations/siem/sumo-logic/).

## Example Usage

```hcl
resource "cyral_integration_sumo_logic" "some_resource_name" {
    name = ""
    address = ""
}
```

## Argument Reference

- `name` - (Required) Integration name that will be used internally in Control Plane.
- `address` - (Required) Sumo Logic Address.

## Attribute Reference

- `id` - The ID of this resource.
