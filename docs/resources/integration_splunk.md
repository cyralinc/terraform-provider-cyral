# Splunk Integration Resource

Provides [integration with Splunk](https://cyral.com/docs/integrations/siem/splunk/#procedure).

## Example Usage

```hcl
resource "cyral_integration_splunk" "some_resource_name" {
    name = ""
    access_token = ""
    port = 0
    host = ""
    index = ""
    use_tls = false|true
}
```

## Argument Reference

* `name` - (Required) Integration name that will be used internally in Control Plane.
* `access_token` - (Required) Splunk Access Token.
* `port` - (Required) Splunk Host Port.
* `host` - (Required) Splunk Host.
* `index` - (Required) Splunk data index name.
* `use_tls` - (Required) Should the communication with Splunk use TLS encryption?

## Attribute Reference

* `id` - The ID of this resource.
