# Provider

The provider is the base element and it must be used to inform application-wide parameters.

## Usage

- Terraform v12

```hcl
provider "cyral" {
    client_id = ""     # optional
    client_secret = "" # optional
    control_plane = "some-cp.cyral.com:8000"
}
```

- Terraform v13 and v14

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

## Variables

|  Name             |  Default     |  Description                                                                           | Required |
|:------------------|:------------:|:---------------------------------------------------------------------------------------|:--------:|
| `auth_provider`   | `"keycloak"` | Authorization provider in use by the Control Plane (valid values: `auth0`, `keycloak`) | No       |
| `auth0_audience`  |              | Auth0 audience (ex: `cyral-api.com`)                                                   | No       |
| `auth0_domain`    |              | Auth0 domain name (ex: `dev-cyral.auth0.com`)                                          | No       |
| `client_id`       |              | Client id (ex: `abcdef1234567`)                                                        | No       |
| `client_secret`   |              | Client secret (ex: `0123456789QWERTYUIOPASDFGHJKLZXCVBNM`)                             | No       |
| `control_plane`   |              | Control plane host and API port (ex: `some-cp.cyral.com:8000`)                         | Yes      |
