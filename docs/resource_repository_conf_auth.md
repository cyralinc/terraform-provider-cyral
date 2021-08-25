# Repository Authentication Configuration

CRUD operations for Repository Authentication Configuration.

## Usage

```hcl
resource "cyral_repository_conf_auth" "my-repository-conf-auth" {
    repository_id = ""
    allow_native_auth = true|false
    client_tls = "enable|disable|enabledAndVerifyCertificate"
    identity_provider = ""
    repo_tls = "enable|disable|enabledAndVerifyCertificate"
}
```

## Variables

|  Name               |  Default  |  Description                                                          | Required |
|:--------------------|:---------:|:----------------------------------------------------------------------|:--------:|
| `repository_id`     |           | The ID of the repository to be configured.                            | Yes      |
| `allow_native_auth` |           | Should the comunication allow native authentication?                  | No       |
| `client_tls`        |           | Is the repo Client using TLS?                                         | No       |
| `identity_provider` |           | The name of the okta identity provider.                               | No       |
| `repo_tls`          |           | Is TLS enabled for the repository?                                    | No       |


## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the Control Plane.                     |
