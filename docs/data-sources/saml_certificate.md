# SAML Certificate Data Source

Retrieves a service provider X.509 certificate used for signing SAML requests.

## Example Usage

```hcl
data "cyral_saml_certificate" "some_data_source_name" {
}
```

## Attribute Reference

* `id` - The ID of this data source.
* `certificate` - The service provider X.509 certificate used for signing SAML requests.