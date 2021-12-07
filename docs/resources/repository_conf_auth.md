# Repository Authentication Configuration Resource

Provides a resource that allows configuring the [Repository Authentication settings in the Advanced tab](https://cyral.com/docs/manage-repositories/repo-advanced-settings/#authentication).

## Example Usage

```hcl
resource "cyral_repository_conf_auth" "some_resource_name" {
    repository_id = ""
    allow_native_auth = true|false
    client_tls = "enable|disable|enabledAndVerifyCertificate"
    identity_provider = ""
    repo_tls = "enable|disable|enabledAndVerifyCertificate"
}
```

## Argument Reference

* `repository_id` - (Required) The ID of the repository to be configured.
* `allow_native_auth` - (Optional) Should the communication allow native authentication?
* `client_tls` - (Optional) Is the repo Client using TLS?
* `identity_provider` - (Optional) The name of the identity provider.
* `repo_tls` - (Optional) Is TLS enabled for the repository?


## Attribute Reference

* `id` - The ID of this resource.
