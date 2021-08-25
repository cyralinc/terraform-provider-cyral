# Logstash Integration

CRUD operations for Logstash integration.

## Usage

```hcl
resource "cyral_integration_logstash" "logstash" {
    name = ""
    endpoint = ""
    use_mutual_authentication = false|true
    use_private_certificate_chain = false|true
    use_tls = false|true
}
```

## Variables

|  Name         |  Default  |  Description                                                          | Required |
|:--------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `name`        |           | Integration name that will be used internally in Control Plane.       | Yes      |
| `endpoint`    |           | The endpoint used to connect to logstash.                             | Yes      |
| `use_mutual_authentication`     |           | Should logstash use mutual authentication?          | Yes      |
| `use_private_certificate_chain` |           | Should logstash use private certificate chain?      | Yes      |
| `use_tls`     |           | Should logstash use mutual TLS?                                       | Yes      |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
