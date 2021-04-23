# Identity Map

CRUD operations for identity maps.

## Usage

```hcl
resource "cyral_identity_map" "SOME_RESOURCE_NAME" {
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

## Variables

|  Name                         |  Default  |  Description                                                                         | Required |
|:------------------------------|:---------:|:-------------------------------------------------------------------------------------|:--------:|
| `repository_id`               |           | ID of the repository that will this identity will be associated to.                  | Yes      |
| `repository_local_account_id` |           | ID of the local account that will this identity will be associated to.               | Yes      |
| `identity_type`               |           | Identity type: `user` or `group`.                                                    | Yes      |
| `identity_name`               |           | Identity name. Ex: `myusername`, `me@myemail.com`.                                   | Yes      |
| `access_duration`             |           | Access duration defined as a sum of days, hours, minutes and seconds. If omitted, the access duration will be infinity. | No       |


## Computed Variables

|  Name        |  Description                                                                     |
|:-------------|:---------------------------------------------------------------------------------|
| `id`         | Unique ID defined by joining ``repository_id` and `repository_local_account_id`. |
