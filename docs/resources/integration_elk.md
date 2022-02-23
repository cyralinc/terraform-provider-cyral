# ELK Integration Resource

Provides [integration with ELK](https://cyral.com/docs/integrations/siem/elk/) to push sidecar metrics.

## Example Usage

```hcl
resource "cyral_integration_elk" "some_resource_name" {
    name = ""
    kibana_url = ""
    es_url = ""
}
```

## Argument Reference

- `name` - (Required) Integration name that will be used internally in Control Plane.
- `kibana_url` - (Required) Kibana URL.
- `es_url` - (Required) Elastic Search URL.

## Attribute Reference

- `id` - The ID of this resource.
