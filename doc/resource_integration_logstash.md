# Repository

CRUD operations for Datadog integration.

## Usage

```hcl
resource "cyral_integration_datadog" "SOME_RESOURCE_NAME" {
    name = ""
    api_key = ""
}
```

## Variables

|  Name         |  Default  |  Description                                                          | Required |
|:--------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`        |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `endpoint`        |           | The endpoint used to connect to logstash.       | Yes      |
| `use_mutual_authentication`        |           | Should logstash use mutual authentication?       | Yes      |
| `use_private_certificate_chain`        |           | Should logstash use private certificate chain?       | Yes      |
| `use_tls`        |           | Should logstash use mutual TLS?       | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
