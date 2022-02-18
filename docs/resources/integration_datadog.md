# Datadog Integration Resource

Provides [integration with DataDog](https://cyral.com/docs/integrations/apm/datadog/) to push sidecar logs and/or metrics.

## Example Usage

```hcl
resource "cyral_integration_datadog" "some_resource_name" {
    name = ""
    api_key = ""
}
```

## Argument Reference

- `name` - (Required) Integration name that will be used internally in Control Plane.
- `api_key` - (Required)Datadog API key.

## Attribute Reference

- `id` - The ID of this resource.
