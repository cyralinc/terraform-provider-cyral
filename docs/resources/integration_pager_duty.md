# Pager Duty Integration Resource

Provides [integration with PagerDuty](https://cyral.com/docs/integrations/incident-response/pagerduty/#in-cyral).

## Example Usage

```hcl
resource "cyral_integration_pager_duty" "some_resource_name" {
    name = ""
    api_token = ""
}
```

## Argument Reference

- `name` - (Required) Integration name that will be used internally in Control Plane.
- `api_token` - (Required) API token for the PagerDuty integration.

## Attribute Reference

- `id` - The ID of this resource.
