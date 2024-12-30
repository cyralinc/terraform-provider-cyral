# cyral_rego_policy_instance (Resource)

Manages a Rego Policy instance.

-> **Note** This resource can be used to create repo-level policies by specifying the repo IDs associated to the policy `scope`. For more information, see the [scope](#nestedblock--scope) field.

-> Import ID syntax is `{category}/{policy_id}`.

## Example Usage

```terraform
### Global rego policy instance
resource "cyral_rego_policy_instance" "policy" {
  name = "some-rate-limit-policy"
  category = "SECURITY"
  description = "Some policy description."
  template_id = "rate-limit"
  parameters = "{\"rateLimit\":7,\"labels\":[\"EMAIL\"],\"alertSeverity\":\"high\",\"block\":false,\"identities\":{\"included\":{\"groups\":[\"analysts\"]}},\"dbAccounts\":{\"included\":[\"admin\"]}}"
  enabled = true
  tags = ["tag1", "tag2"]
}

output "policy_last_updated" {
  value = cyral_rego_policy_instance.policy.last_updated
}

output "policy_created" {
  value = cyral_rego_policy_instance.policy.created
}

### Repo-level policy
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
      host = "mysql.cyral.com"
      port = 3306
  }
}

resource "cyral_rego_policy_instance" "policy" {
  name = "some-data-masking-policy"
  category = "SECURITY"
  description = "Some policy description."
  template_id = "data-masking"
  parameters = "{\"labels\":[\"ADDRESS\"],\"maskType\":\"NULL_MASK\"}"
  enabled = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}

### Rego policy instance with duration
resource "cyral_rego_policy_instance" "policy" {
  name = "some-data-masking-policy"
  category = "SECURITY"
  description = "Some policy description."
  template_id = "data-masking"
  parameters = "{\"labels\":[\"ADDRESS\"],\"maskType\":\"NULL_MASK\"}"
  enabled = true
  tags = ["tag1", "tag2"]
  duration = "10s"
}
```

## Template Parameters

All templates use parameters defined as JSON, below is a list of all the corresponding parameters for the predefined templates.

-> You can also use the Cyral API `GET` `/v1/regopolicies/templates` to retrieve all existing templates and their corresponding parameters schema.

### Object Protection (object-protection)

-   `objectType` (String) The type of object to monitor or protect. Supported types include tables, views, roles/users, and schemas. Specific actions depend on the object type.
-   `block` (Boolean) Indicates whether unauthorized operations should be blocked. If true, operations violating the policy are prevented..
-   `monitorCreates` (Boolean) Specifies whether to monitor 'CREATE' operations for the defined object type. Applies only to relevant object types.
-   `monitorDrops` (Boolean) Specifies whether to monitor 'DROP' operations for the defined object type. Applies only to relevant object types.
-   `monitorAlters` (Boolean) Specifies whether to monitor 'ALTER' operations for the defined object type. Applies only to relevant object types.
-   `objects` (Array) A list of specific objects (e.g., tables or views) to monitor or protect. Required for 'table' or 'view' object types. Not applicable to 'role/user' or 'schema'.
-   `identities` (Object) Defines users, groups, or emails that are included or excluded from the policy. If included identities are defined, only those users are exempt from policy enforcement. Excluded identities are always subject to the policy.
-   `dbAccounts` (Object) Defines database accounts to include or exclude from the policy. Excluded accounts are not subject to the policy, while included accounts must adhere to it.

### Service Account Abuse (service-account-abuse)

-   `block` (Boolean) Policy action to enforce.
-   `serviceAccounts` (Array) Service accounts for which end user attribution is always required.
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.
