# Repository

CRUD operations for Slack Alerts integration.

## Usage

```hcl
resource "cyral_integration_slack_alerts" "SOME_RESOURCE_NAME" {
    name = ""
    url = ""
}
```

## Variables

|  Name         |  Default  |  Description                                                          | Required |
|:--------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`        |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `url`     |           | Slack Alert Webhook url.                                                   | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |