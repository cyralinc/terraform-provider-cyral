# cyral_rego_policy_instance (Resource)

Manages a Rego Policy instance.

-> **Note** This resource can be used to create repo-level policies by specifying the repo IDs associated to the policy `scope`. For more information, see the [scope](#nestedblock--scope) field.

-> Import ID syntax is `{category}/{policy_id}`.

## Example Usage

```terraform
### Global rego policy instance
resource "cyral_rego_policy_instance" "policy" {
  name = "User Management"
  category = "SECURITY"
  description = "Policy to govern user management operations"
  template_id = "object-protection"
  parameters = jsonencode(
    {
      "objectType" = "role/user"
      "block" = true
      "monitorCreates" = true
      "monitorAlters" = true
      "monitorDrops" = true
      "identities" = {
        "excluded" = {
          "groups" = ["dba"]
        }
      }
    }
  )
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
  name = "User Management"
  category = "SECURITY"
  description = "Policy to govern user management operations"
  template_id = "object-protection"
  parameters = jsonencode(
    {
      "objectType" = "role/user"
      "block" = true
      "monitorCreates" = true
      "monitorAlters" = true
      "monitorDrops" = true
      "identities" = {
        "excluded" = {
          "groups" = ["dba"]
        }
      }
    }
  )
  enabled = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}

### Rego policy instance with duration
resource "cyral_rego_policy_instance" "policy" {
  name = "User Management"
  category = "SECURITY"
  description = "Policy to govern user management operations"
  template_id = "object-protection"
  parameters = jsonencode(
    {
      "objectType" = "role/user"
      "block" = true
      "monitorCreates" = true
      "monitorAlters" = true
      "monitorDrops" = true
      "identities" = {
        "excluded" = {
          "groups" = ["dba"]
        }
      }
    }
  )
  enabled = true
  tags = ["tag1", "tag2"]
  duration = "10s"
}
```

## Template Parameters

All templates use parameters defined as JSON, below is a list of all the corresponding parameters for the predefined templates.

-> You can also use the Cyral API `GET` `/v1/regopolicies/templates` to retrieve all existing templates and their corresponding parameters schema.

### Fail Closed (fail-closed) - Protect against statements that are not understood by Cyral.

-   `block` (Boolean) Indicates whether unauthorized operations should be blocked. If true, operations violating the policy are prevented.
-   `identities` (Object) Defines users, groups, or emails that are included or excluded from the policy. If included identities are defined, only those users are exempt from policy enforcement. Excluded identities are always subject to the policy. See [identities](#objects--identities).
-   `dbAccounts` (Object) Defines database accounts to include or exclude from the policy. Excluded accounts are not subject to the policy, while included accounts must adhere to it. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.

### Object Protection (object-protection) - Guards against operations like create, drop, and alter for specified object types.

-   `objectType` (String) The type of object to monitor or protect. Supported types include tables, views, roles/users, and schemas. Specific actions depend on the object type.
-   `block` (Boolean) Indicates whether unauthorized operations should be blocked. If true, operations violating the policy are prevented.
-   `monitorCreates` (Boolean) Specifies whether to monitor 'CREATE' operations for the defined object type. Applies only to relevant object types.
-   `monitorDrops` (Boolean) Specifies whether to monitor 'DROP' operations for the defined object type. Applies only to relevant object types.
-   `monitorAlters` (Boolean) Specifies whether to monitor 'ALTER' operations for the defined object type. Applies only to relevant object types.
-   `objects` (Array) A list of specific objects (e.g., tables or views) to monitor or protect. Required for 'table' or 'view' object types. Not applicable to 'role/user' or 'schema'.
-   `identities` (Object) Defines users, groups, or emails that are included or excluded from the policy. If included identities are defined, only those users are exempt from policy enforcement. Excluded identities are always subject to the policy. See [identities](#objects--identities).
-   `dbAccounts` (Object) Defines database accounts to include or exclude from the policy. Excluded accounts are not subject to the policy, while included accounts must adhere to it. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.

### Service Account Abuse (service-account-abuse) - Ensure service accounts can only be used by intended applications.

-   `block` (Boolean) Policy action to enforce.
-   `serviceAccounts` (Array) Service accounts for which end user attribution is always required.
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.

### Stored Procedure Governance (stored-procedure-governance) - Restrict execution of stored procedures.

-   `governedProcedures` (Array) List of stored procedures to be governed.
-   `enforce` (Boolean) Whether to enforce the policy, if false, only alerts will be raised on policy violations.
-   `identities` (Object) Defines users, groups, or emails that are included or excluded from the policy. If included identities are defined, only those users are exempt from policy enforcement. Excluded identities are always subject to the policy. See [identities](#objects--identities).
-   `dbAccounts` (Object) Defines database accounts to include or exclude from the policy. Excluded accounts are not subject to the policy, while included accounts must adhere to it. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.

### Ungoverned Statements (ungoverned-statements) - Control execution of statements not governed by other policies.

-   `block` (Boolean) Indicates whether unauthorized operations should be blocked. If true, operations violating the policy are prevented.
-   `identities` (Object) Defines users, groups, or emails that are included or excluded from the policy. If included identities are defined, only those users are exempt from policy enforcement. Excluded identities are always subject to the policy. See [identities](#objects--identities).
-   `dbAccounts` (Object) Defines database accounts to include or exclude from the policy. Excluded accounts are not subject to the policy, while included accounts must adhere to it. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.

### Deprecated policy templates

The remaining list of policy templates have been deprecated in v4.18.X of the Cyral Control Plane
and can not be used for creating new policies. Managing existing policy instances is still supported.
Please visit [`cyral_policy_set`](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/resources/policy_set)
resource to find replacements for the deprecated policy templates.

#### Data Firewall (data-firewall)

-   `dataSet` (String) Data Set.
-   `dataFilter` (String) Data filter that will be applied when anyone tries to read the specified data labels from the data set.
-   `tags` (Array) Tags.
-   `labels` (Array) Data Labels.
-   `excludedIdentities` (Object) Identities that will be excluded from this policy. See [identityList](#objects--identityList).

#### Data Masking (data-masking)

-   `maskType` (String) Mask Type (E.g.: `NULL_MASK`, `CONSTANT_MASK`, `MASK`).
-   `maskArguments` (Array) Mask Argument associated to the given Mask Type (E.g.: Replacement Value).
-   `tags` (Array) Tags.
-   `labels` (Array) Data Labels.
-   `identities` (Object) Identities associated to the policy. If empty, the policy will be associated to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts associated to the policy. If empty, the policy will be associated to any database account. See [dbAccounts](#objects--dbAccounts).

#### Data Protection (data-protection)

-   `block` (Boolean) Policy action to block.
-   `monitorReads` (Boolean) Monitor read operations.
-   `monitorUpdates` (Boolean) Monitor update operations.
-   `monitorDeletes` (Boolean) Monitor delete operations.
-   `tags` (Array) Tags.
-   `labels` (Array) Data Labels.
-   `identities` (Object) Identities associated to the policy. If empty, the policy will be associated to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts associated to the policy. If empty, the policy will be associated to any database account. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.

#### Ephemeral Grant (EphemeralGrantPolicy)

-   `repoAccount` (String) Repository Account Name.
-   `repo` (String) Repository Name.
-   `allowedSensitiveAttributes` (Array) Allowed Sensitive Attributes.

#### Rate Limit (rate-limit)

-   `rateLimit` (Integer) Maximum number of rows that can be returned per hour. Note: the value must be an integer greater than zero.
-   `block` (Boolean) Policy action to enforce.
-   `tags` (Array) Tags.
-   `labels` (Array) Data Labels.
-   `identities` (Object) Identities associated to the policy. If empty, the policy will be associated to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts associated to the policy. If empty, the policy will be associated to any database account. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.

#### Read Limit (read-limit)

-   `rowLimit` (Integer) Maximum number of rows that can be read per query. Note: the value must be an integer greater than zero.
-   `block` (Boolean) Policy action to enforce.
-   `appliesToAllData` (Boolean) Whether the policy should apply to the entire repository data.
-   `tags` (Array) Tags.
-   `labels` (Array) Data Labels.
-   `identities` (Object) Identities associated to the policy. If empty, the policy will be associated to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts associated to the policy. If empty, the policy will be associated to any database account. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.

#### Repository Protection (repository-protection)

-   `rowLimit` (Integer) Maximum number of rows that can be modified per query. Note: the value must be an integer greater than zero.
-   `monitorUpdates` (Boolean) Monitor update operations.
-   `monitorDeletes` (Boolean) Monitor delete operations.
-   `identities` (Object) Identities associated to the policy. If empty, the policy will be associated to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts associated to the policy. If empty, the policy will be associated to any database account. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Policy action to alert, using the respective severity. Allowed values are: `low`, `medium`, `high`.

#### User Segmentation (user-segmentation)

-   `dataSet` (String) Data Set.
-   `dataFilter` (String) Data filter that will be applied when anyone tries to read the specified data labels from the data set.
-   `tags` (Array) Tags.
-   `labels` (Array) Data Labels.
-   `includedIdentities` (Object) Identities that cannot see restricted records. See [identityList](#objects--identityList).
-   `includedDbAccounts` (Array) Database accounts cannot see restricted records.

<a id="parameter-objects"></a>

### Objects

<a id="objects--identities"></a>

-   `identities` (Object) Identities. See properties below:
    -   `included` (Object) Included Identities. See [identityList](#objects--identityList).
    -   `excluded` (Object) Excluded Identities. See [identityList](#objects--identityList).
        <a id="objects--dbAccounts"></a>
-   `dbAccounts` (Object) Database Accounts. See properties below:
    -   `included` (Array) Included Database Accounts.
    -   `excluded` (Array) Excluded Database Accounts.
        <a id="objects--identityList"></a>
-   `identityList` (Object) Identity List. See properties below:
    -   `userNames` (Array) Identity Emails.
    -   `emails` (Array) Identity Usernames.
    -   `groups` (Array) Identity Groups.

<!-- schema generated by tfplugindocs -->

## Schema

### Required

-   `category` (String) Policy category. List of supported categories:
    -   `SECURITY`
    -   `GRANT`
    -   `USER_DEFINED`
-   `name` (String) Policy name.
-   `template_id` (String) Policy template identifier. Predefined templates are:
    -   `data-firewall`
    -   `data-masking`
    -   `data-protection`
    -   `EphemeralGrantPolicy`
    -   `rate-limit`
    -   `read-limit`
    -   `repository-protection`
    -   `service-account-abuse`
    -   `user-segmentation`

### Optional

-   `description` (String) Policy description.
-   `duration` (String) Policy duration. The policy expires after the duration specified. Should follow the protobuf duration string format, which corresponds to a sequence of decimal numbers suffixed by a 's' at the end, representing the duration in seconds. For example: `300s`, `60s`, `10.50s`, etc.
-   `enabled` (Boolean) Enable/disable the policy. Defaults to `false` (Disabled).
-   `parameters` (String) Policy parameters. The parameters vary based on the policy template schema.
-   `scope` (Block Set, Max: 1) Determines the scope that the policy applies to. It can be used to create a repo-level policy by specifying the corresponding `repo_ids` that this policy should be applied. (see [below for nested schema](#nestedblock--scope))
-   `tags` (List of String) Tags that can be used to categorize the policy.

### Read-Only

-   `created` (Set of Object) Information regarding the policy creation. (see [below for nested schema](#nestedatt--created))
-   `id` (String) The resource identifier. It is a composed ID that follows the format `{category}/{policy_id}`.
-   `last_updated` (Set of Object) Information regarding the policy last update. (see [below for nested schema](#nestedatt--last_updated))
-   `policy_id` (String) ID of this rego policy instance in Cyral environment.

<a id="nestedblock--scope"></a>

### Nested Schema for `scope`

Required:

-   `repo_ids` (List of String) A list of repository identifiers that belongs to the policy scope. The policy will be applied at repo-level for every repository ID included in this list. This is equivalent of creating a repo-level policy in the UI for a given repository.

<a id="nestedatt--created"></a>

### Nested Schema for `created`

Read-Only:

-   `actor` (String)
-   `actor_type` (String)
-   `timestamp` (String)

<a id="nestedatt--last_updated"></a>

### Nested Schema for `last_updated`

Read-Only:

-   `actor` (String)
-   `actor_type` (String)
-   `timestamp` (String)
