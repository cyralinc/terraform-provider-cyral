# SAML Certificate Data Source

Retrieves a X.509 certificate used for signing SAML requests.

## Example Usage

```hcl
data "cyral_saml_certificate" "some_data_source_name" {
}
```

## Attribute Reference

- `id` - The ID of this data source.
- `certificate` - The X.509 certificate used for signing SAML requests.
