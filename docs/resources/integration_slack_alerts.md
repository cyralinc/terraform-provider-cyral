# Slack Alerts Integration Resource

Provides integration with Slack to push alerts.

## Example Usage

```hcl
resource "cyral_integration_slack_alerts" "some_resource_name" {
    name = ""
    url = ""
}
```

## Argument Reference

* `name` - (Required) Integration name that will be used internally in Control Plane.
* `url` - (Required) Slack Alert Webhook url.

## Attribute Reference

* `id` - The ID of this resource.
