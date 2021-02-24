# Provider

The provider is the base element and it must be used to inform application-wide parameters.

## Usage

- Terraform v12

```hcl
provider "cyral" {
    auth0_domain = ""
    auth0_audience = ""
    control_plane = ""
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
    auth0_domain = ""
    auth0_audience = ""
    control_plane = ""
}
```

----

Auth0 authentication parameters (`client ID` and `client secret`) must be configured as environment variables **before** running Terraform command line. Fill the parameters with the corresponding values taken from the Auth0 application configuration and run the following commands to create the environment variables:

- Linux/Mac

```bash
export AUTH0_CLIENT_ID=""
export AUTH0_CLIENT_SECRET=""
```

- Windows

```
set AUTH0_CLIENT_ID=""
set AUTH0_CLIENT_SECRET=""
```

## Variables

|  Name                    |  Default  |  Description                                                      | Required |
|:-------------------------|:---------:|:------------------------------------------------------------------|:--------:|
| `auth0_domain`           |           | Auth0 domain name (ex: `dev-cyral.auth0.com`)                     | Yes      |
| `auth0_client_id`        |           | Auth0 client id (ex: `1nrd81340lskf`)                             | Yes      |
| `auth0_client_secret`    |           | Auth0 client secret (ex: `klfd;3rf-0e13jklehgjlkhjf31J:LkfdsjfA`) | Yes      |
| `auth0_audience`         |           | Auth0 audience (ex: `cyral-api.com`)                              | Yes      |
| `control_plane`          |           | Control plane host and API port (ex: `some-cp.cyral.com:8000`)    | Yes      |
