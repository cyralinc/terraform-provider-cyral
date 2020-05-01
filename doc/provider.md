# Provider

The provider is the base element and it must be used to inform application-wide parameters.

## Usage

```hcl
provider "cyral" {
    auth0_domain = ""
    auth0_client_id = ""
    auth0_client_secret = ""
    auth0_audience = ""
    control_plane = ""
}
```

## Variables

|  Name                 |  Default  |  Description                                                      | Required |
|:----------------------|:---------:|:------------------------------------------------------------------|:--------:|
| `auth0_domain`        |           | Auth0 domain name (ex: `dev-cyral.auth0.com`)                     | Yes      |
| `auth0_client_id`     |           | Auth0 client id (ex: `1nrd81340lskf`)                             | Yes      |
| `auth0_client_secret` |           | Auth0 client secret (ex: `klfd;3rf-0e13jklehgjlkhjf31J:LkfdsjfA`) | Yes      |
| `auth0_audience`      |           | Auth0 audience (ex: `cyral-api.com`)                              | Yes      |
| `control_plane`       |           | Control plane host (ex: `yourcp.cyral.com`)                       | Yes      |
