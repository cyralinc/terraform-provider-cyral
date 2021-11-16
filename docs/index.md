# Provider

The provider is the base element and must be used to inform application-wide
parameters, like the Cyral control plane reference and authentication secrets.

## Example Usage

- Terraform v12

```hcl
provider "cyral" {
    client_id = ""     # optional
    client_secret = "" # optional
    control_plane = "some-cp.cyral.com:8000"
}
```

- Terraform v13+

```hcl
terraform {
  required_providers {
    cyral = {
      source = "cyral.com/terraform/cyral"
    }
  }
}

provider "cyral" {
    client_id = ""     # optional
    client_secret = "" # optional
    control_plane = "some-cp.cyral.com:8000"
}
```

## Argument Reference

* `auth_provider` - (Optional) Authorization provider in use by the Control Plane (valid values: `auth0`, `keycloak`). Default: `keycloak`.
* `auth0_audience` - (Optional) Auth0 audience.
* `auth0_domain` - (Optional) Auth0 domain name.
* `client_id` - (Optional) Client id used to authenticate against the Control Plane.
* `client_secret` - (Optional) Client secret used to authenticate against the Control Plane.
* `control_plane` - (Required) Control plane host and API port (ex: `some-cp.cyral.com:8000`)

----

Authentication parameters `client_id` and `client_secret` are defined as optional in the provider body once they can be set through environment variables in order to avoid storing secrets in source code repositories. The environment variables corresponds to `CYRAL_TF_CLIENT_ID` and `CYRAL_TF_CLIENT_SECRET` respectivelly and can be defined as follows:

- Linux/Mac

```bash
export CYRAL_TF_CLIENT_ID=""
export CYRAL_TF_CLIENT_SECRET=""
```

- Windows

```
set CYRAL_TF_CLIENT_ID=""
set CYRAL_TF_CLIENT_SECRET=""
```