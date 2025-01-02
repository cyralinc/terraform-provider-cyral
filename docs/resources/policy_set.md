# cyral_policy_set (Resource)

This resource allows management of policy sets in the Cyral platform.

-> Import ID syntax is `{policy_set_id}`.

## Example Usage

```terraform
resource "cyral_repository" "myrepo" {
    type = "mongodb"
    name = "myrepo"

    repo_node {
        name = "node-1"
        host = "mongodb.cyral.com"
        port = 27017
    }

    mongodb_settings {
      server_type = "standalone"
    }
}

resource "cyral_policy_set" "repo_lockdown_example" {
  wizard_id    = "repo-lockdown"
  name        = "default block with failopen"
  description = "This default policy will block by default all queries for myrepo except the ones not parsed by Cyral"
  enabled = true
  tags = ["default block", "fail open"]
    scope {
    repo_ids = [cyral_repository.myrepo.id]
  }
  wizard_parameters    = jsonencode({
    denyByDefault = true
    failClosed = false
    })
}
```

## Available Policy Wizards

The following policy wizards are available for creating policy sets. The wizard parameters,
specified as a JSON object, are described below for each wizard as well.

-> You can also use the Cyral API `GET` `/v1/regopolicies/templates` to retrieve all existing templates and their corresponding parameters schema.

### Data Firewall (data-firewall) - Ensure that sensitive data can only be read by specified individuals.

-   `dataset` (String) Data Set (table, collection, etc.) to which the policy applies.
-   `dataFilter` (String) Data filter that will be applied when anyone tries to read the specified data labels from the data set.
-   `substitutionQuery` (String) A query that will be used to replace all occurrences of the dataset in the original query. Only one of `dataFilter` and `substitutionQuery` can be specified.
-   `excludedIdentities` (Object) Identities that will be excluded from this policy. See [identityList](#objects--identityList).

### Data Masking (data-masking) - Mask fields for specific users and applications.

-   `maskType` (String) Mask Type (E.g.: `null`, `constant`, `format-preserving`).
-   `maskArguments` (Array) Mask Argument associated to the given Mask Type (E.g.: Replacement Value).
-   `tags` (Array) Data Tags to which the policy applies.
-   `labels` (Array) Data Labels to which the policy applies.
-   `identities` (Object) Identities to which the policy applies. If empty, the policy will apply to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts to which the policy applies. If empty, the policy will apply to any database account. See [dbAccounts](#objects--dbAccounts).

### Data Protection (data-protection) - Guard against reads and writes of specified tables or fields.

-   `block` (Boolean) Policy action to block.
-   `governedOperations` (Array) Operations governed by this policy, can be one or more of: `read`, `update`, `delete`, and `insert`.
-   `tags` (Array) Data Tags to which the policy applies.
-   `labels` (Array) Data Labels to which the policy applies.
-   `datasets` (Array) Data Sets (tables, collections, etc.) to which the policy applies.
-   `identities` (Object) Identities to which the policy applies. If empty, the policy will be applied to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts to which the policy applies. If empty, the policy will be applied to any database account. See [dbAccounts](#objects--dbAccounts).

### Object Protection (object-protection) - Guards against operations like create, drop, and alter for specified object types.

-   `objectType` (String) The type of object to monitor or protect. The only value currently supported is `role/user`.
-   `block` (Boolean) Indicates whether unauthorized operations should be blocked. If true, operations violating the policy are prevented.
-   `governedOperations` (Array) Operations governed by this policy, can be one or more of: `create`, `drop`, and `alter`.
-   `identities` (Object) Identities to which the policy applies. If empty, the policy will be applied to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts to which the policy applies. If empty, the policy will be applied to any database account. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Alert severity. Allowed values are: `low`, `medium`, `high`.

### Rate Limit (rate-limit) - Implement threshold on sensitive data reads over a period of time.

-   `rateLimit` (Integer) Maximum number of rows that can be returned per hour. Note: the value must be an integer greater than zero.
-   `enforce` (Boolean) Whether to enforce the policy, if false, only alerts will be raised on policy violations.
-   `tags` (Array) Data Tags to which the policy applies.
-   `labels` (Array) Data Labels to which the policy applies.
-   `identities` (Object) Identities to which the policy applies. If empty, the policy will be applied to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts to which the policy applies. If empty, the policy will be applied to any database account. See [dbAccounts](#objects--dbAccounts).

### Read Limit (read-limit) - Prevent certain data from being read beyond a specified limit.

-   `rowLimit` (Integer) Maximum number of rows that can be read per query. Note: the value must be an integer greater than zero.
-   `enforce` (Boolean) Whether to enforce the policy, if false, only alerts will be raised on policy violations.
-   `tags` (Array) Data Tags to which the policy applies.
-   `labels` (Array) Data Labels to which the policy applies.
-   `datasets` (Array) Data Sets (tables, collections, etc.) to which the policy applies.
-   `identities` (Object) Identities to which the policy applies. If empty, the policy will be applied to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts to which the policy applies. If empty, the policy will be applied to any database account. See [dbAccounts](#objects--dbAccounts).

### Repository Lockdown (repo-lockdown) - Deny all statements that are not allowed by some policy and/or not understood by Cyral.

-   `failClosed` (Boolean) Whether to fail closed, if true, all statements that are not understood by Cyral will be blocked.
-   `denyByDefault` (Boolean) Whether to deny all statements by default, if true, all statements that are not explicitly allowed by some policy will be blocked.

### Repository Protection (repository-protection) - Alert when more than a specified number of records are updated, deleted, or inserted in specified datasets.

-   `rowLimit` (Integer) Maximum number of rows that can be modified per query. Note: the value must be an integer greater than zero.
-   `governedOperations` (Array) Operations governed by this policy, can be one or more of: `update`, `delete` and `insert`.
-   `datasets` (Array) Data Sets (tables, collections, etc.) to which the policy applies.
-   `identities` (Object) Identities to which the policy applies. If empty, the policy will be applied to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts to which the policy applies. If empty, the policy will be applied to any database account. See [dbAccounts](#objects--dbAccounts).

### Schema Protection (schema-protection) - Protect database schema against unauthorized creation, deletion, or modification of tables and views.

-   `block` (Boolean) Whether to block unauthorized schema changes.
-   `schemas` (Array) Schemas to which the policy applies.
-   `excludedIdentities` (Object) Identities that are exempt from the policy. See [identities](#objects--identityList).

### Service Account Abuse (service-account-abuse) - Ensure service accounts can only be used by intended applications.

-   `block` (Boolean) Policy action to enforce.
-   `serviceAccounts` (Array) Service accounts for which end user attribution is always required.
-   `alertSeverity` (String) Alert severity. Allowed values are: `low`, `medium`, `high`.

### Stored Procedure Governance (stored-procedure-governance) - Restrict execution of stored procedures..

-   `enforced` (Boolean) Whether to enforce the policy, if false, only alerts will be raised on policy violations.
-   `governedProcedures` (Array) Stored procedures to which the policy applies.
-   `identities` (Object) Identities to which the policy applies. If empty, the policy will be applied to all identities. See [identities](#objects--identities).
-   `dbAccounts` (Object) Database Accounts to which the policy applies. If empty, the policy will be applied to any database account. See [dbAccounts](#objects--dbAccounts).
-   `alertSeverity` (String) Alert severity. Allowed values are: `low`, `medium`, `high`.

### User Segmentation (user-segmentation) - Restrict specific users to a subset of data.

-   `dataset` (String) Data Set (table, collection, etc.) to which the policy applies.
-   `dataFilter` (String) Data filter that will be applied when anyone tries to read the specified data labels from the data set.
-   `substitutionQuery` (String) A query that will be used to replace all occurrences of the dataset in the original query. Only one of `dataFilter` and `substitutionQuery` can be specified.
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

-   `name` (String) Name of the policy set.
-   `wizard_id` (String) The ID of the policy wizard used to create this policy set.
-   `wizard_parameters` (String) Parameters passed to the wizard while creating the policy set.

### Optional

-   `description` (String) Description of the policy set.
-   `enabled` (Boolean) Indicates if the policy set is enabled.
-   `scope` (Block List, Max: 1) Scope of the policy set. (see [below for nested schema](#nestedblock--scope))
-   `tags` (List of String) Tags associated with the policy set.

### Read-Only

-   `created` (Map of String) Information about when and by whom the policy set was created.
-   `id` (String) Identifier for the policy set.
-   `last_updated` (Map of String) Information about when and by whom the policy set was last updated.
-   `policies` (List of Object) List of policies that comprise the policy set. (see [below for nested schema](#nestedatt--policies))

<a id="nestedblock--scope"></a>

### Nested Schema for `scope`

Optional:

-   `repo_ids` (List of String) List of repository IDs that are in scope. Empty list means all repositories are in scope.

<a id="nestedatt--policies"></a>

### Nested Schema for `policies`

Read-Only:

-   `id` (String)
-   `type` (String)
