# Policy

This resource provides CRUD operations in Cyral policies, allowing users to Create, Read, Update and Delete policies.

## Usage

```hcl
resource "cyral_policy" "SOME_RESOURCE_NAME" {
  data = [""]
  description = ""
  enabled = true | false
  name = ""
  tags = [""]
}
```

## Variables

|  Name           |  Default  |  Description                                                                         | Required |
|:----------------|:---------:|:-------------------------------------------------------------------------------------|:--------:|
| `data`         |           | List that specify which data fields a policy manages. Each field is represented by the LABEL you established for it in your data map. The actual location of that data (the names of fields, columns, or databases that hold it) is listed in the data map.                   | No     |
| `description`  |           | String that describes the policy (ex: `your_policy_description`).  | No      |
| `enabled`      | `true`      | Boolean that causes a policy to be enabled or disabled.  | Yes      |
| `name`         |     | Policy name that will be used internally in Control Plane (ex: `your_policy_name`).   | No      |
| `properties`   |           |   | No      |
| `tags`         |           | Tags that can be used to organize and/or classify your policies (ex: `[your_tag1, your_tag2]`).  | No      |
| `type`         |           |   | No      |


## Outputs

|  Name          |  Description                                                        |
|:---------------|:--------------------------------------------------------------------|
| `id`           | Unique ID of the resource in the Control Plane.                     |
| `created`      | Policy creation timestamp.                                          |
| `last_updated` | Last update timestamp.                                              |
| `version`      | Incremental counter for every update on the policy.                 |

