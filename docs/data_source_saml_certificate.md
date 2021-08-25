# SAML Certificate

X.509 certificate for SAML integration validation.

## Usage

```hcl
data "cyral_saml_certificate" "SOME_DATA_SOURCE_NAME" {
}
```

## Inputs

| Name | Default | Description | Required |
|:-----|:--------|:------------|:---------|


## Outputs

| Name | Description |
|:-----|:------------|
| `certificate` | the certificate used for signing saml requests.| 