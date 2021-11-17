# Logstash Integration Resource

Provides integration with Logstash.

## Example Usage

```hcl
resource "cyral_integration_logstash" "some_resource_name" {
    name = ""
    endpoint = ""
    use_mutual_authentication = false|true
    use_private_certificate_chain = false|true
    use_tls = false|true
}
```

## Argument Reference

* `name` - (Required) Integration name that will be used internally in Control Plane.
* `endpoint` - (Required) The endpoint used to connect to logstash.
* `use_mutual_authentication` - (Required) Should logstash use mutual authentication?
* `use_private_certificate_chain` - (Required) Should logstash use private certificate chain?
* `use_tls` - (Required) Should logstash use mutual TLS?


## Attribute Reference

* `id` - The ID of this resource.
