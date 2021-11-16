# Repository Identity Map Resource

Provides Repository Identity Maps configuration.

## Example Usage

```hcl
resource "cyral_repository_identity_map" "some_resource_name" {
    repository_id = ""
    repository_local_account_id = ""
    identity_type = "user|group"
    identity_name = ""
    access_duration {
        days    = 0
        hours   = 0
        minutes = 0
        seconds = 0
    }
}
```

## Argument Reference

* `repository_id` - (Required) ID of the repository that this identity will be associated to.
* `repository_local_account_id` - (Required) ID of the local account that this identity will be associated to.
* `identity_type` - (Required) Identity type: `user` or `group`.
* `identity_name` - (Required) Identity name. Ex: `myusername`, `me@myemail.com`.
* `access_duration` - (Optional) Access duration defined as a sum of days, hours, minutes and seconds. If omitted or all fields are set to zero, the access duration will be infinity.

## Attribute Reference

* `id` - The ID of this resource.
