# Repository

CRUD operations for Splunk integration.

## Usage

```hcl
resource "cyral_integration_splunk" "SOME_RESOURCE_NAME" {
    name = ""
    access_token = ""
    port = 0
    host = ""
    index = ""
    use_tls = false|true
}
```

## Variables

|  Name          |  Default  |  Description                                                          | Required |
|:---------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`         |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `access_token` |           | Splunk Access Token.                                                  | Yes      |
| `port`         |           | Splunk Host Port.                                                     | Yes      |
| `host`         |           | Splunk Host.                                                          | Yes      |
| `index`        |           | Splunk data index name.                                               | Yes      |
| `use_tls`      |           | Should the comunication with Splunk use TLS encryption?               | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |