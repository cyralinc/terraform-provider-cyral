# cyral_policy_rule (Resource)

~> **DEPRECATED** For control planes `>= v4.15`, use resource `cyral_policy_v2` instead.

-> Import ID syntax is `{policy_id}/{policy_rule_id}`, where `{policy_rule_id}` is the ID of the policy rule in the Cyral Control Plane.

## Example Usage

```terraform
# An example of a policy and a policy rule with a rego policy
# in `additional_checks`.
resource "cyral_policy" "this" {
  name = "My first policy"
  description = "This is my first policy"
  enabled = true
  data = ["EMAIL"]
  metadata_tags = ["Risk Level 1"]
}

resource "cyral_policy_rule" "this" {
  policy_id = cyral_policy.this.id
  deletes {
    additional_checks = <<EOT
is_valid_request {
  filter := request.filters[_]
  filter.field == "entity.user.is_real"
  filter.op == "="
  filter.value == false
}
EOT
    data = ["EMAIL"]
    rows = -1
    severity = "low"
  }
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

-   `policy_id` (String) The ID of the policy you are adding this rule to.

### Optional

-   `deletes` (Block List) A contexted rule for accesses of the type `delete`. (see [below for nested schema](#nestedblock--deletes))
-   `hosts` (List of String) Hosts specification that limits access to only those users connecting from a certain network location.
-   `identities` (Block List, Max: 1) Identities specifies the people, applications, or groups this rule applies to. Every rule except your default rule has one. It can have 4 fields: `db_roles`, `groups`, `users` and `services`. (see [below for nested schema](#nestedblock--identities))
-   `reads` (Block List) A contexted rule for accesses of the type `read`. (see [below for nested schema](#nestedblock--reads))
-   `updates` (Block List) A contexted rule for accesses of the type `update`. (see [below for nested schema](#nestedblock--updates))

### Read-Only

-   `id` (String) The ID of this resource.
-   `policy_rule_id` (String) The ID of the policy rule.

<a id="nestedblock--deletes"></a>

### Nested Schema for `deletes`

Required:

-   `data` (List of String) The data locations protected by this rule. Use `*` if you want to define `any` data location. For more information, see the [policy rules](https://cyral.com/docs/policy/rules#contexted-rules) documentation.
-   `rows` (Number) The number of records (for example, rows or documents) that can be accessed/affected in a single statement. Use positive integer numbers to define how many records. If you want to define `any` number of records, set to `-1`.

Optional:

-   `additional_checks` (String) Constraints on the data access specified in [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/). See [Additional checks](https://cyral.com/docs/policy/rules/#additional-checks).
-   `dataset_rewrites` (Block List) Defines how requests should be rewritten in the case of policy violations. See [Request rewriting](https://cyral.com/docs/policy/rules/#request-rewriting). (see [below for nested schema](#nestedblock--deletes--dataset_rewrites))
-   `rate_limit` (Number) Rate Limit specifies the limit of calls that a user can make within a given time period.
-   `severity` (String) severity level that's recorded when someone violate this rule. This is an informational value. Settings: (`low` | `medium` | `high`). If not specified, the severity is considered to be low.

<a id="nestedblock--deletes--dataset_rewrites"></a>

### Nested Schema for `deletes.dataset_rewrites`

Required:

-   `dataset` (String) The dataset that should be rewritten.In the case of Snowflake, this denotes a fully qualified table name in the form: `<database>.<schema>.<table>`
-   `parameters` (List of String) The set of parameters used in the substitution request, these are references to fields in the activity log as described in the [Additional Checks section](https://cyral.com/docs/policy/rules/#additional-checks).
-   `repo` (String) The name of the repository that the rewrite applies to.
-   `substitution` (String) The request used to substitute references to the dataset.

<a id="nestedblock--identities"></a>

### Nested Schema for `identities`

Optional:

-   `db_roles` (List of String) Database roles that this rule will apply to.
-   `groups` (List of String) Groups that this rule will apply to.
-   `services` (List of String) Services that this rule will apply to.
-   `users` (List of String) Users that this rule will apply to.

<a id="nestedblock--reads"></a>

### Nested Schema for `reads`

Required:

-   `data` (List of String) The data locations protected by this rule. Use `*` if you want to define `any` data location. For more information, see the [policy rules](https://cyral.com/docs/policy/rules#contexted-rules) documentation.
-   `rows` (Number) The number of records (for example, rows or documents) that can be accessed/affected in a single statement. Use positive integer numbers to define how many records. If you want to define `any` number of records, set to `-1`.

Optional:

-   `additional_checks` (String) Constraints on the data access specified in [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/). See [Additional checks](https://cyral.com/docs/policy/rules/#additional-checks).
-   `dataset_rewrites` (Block List) Defines how requests should be rewritten in the case of policy violations. See [Request rewriting](https://cyral.com/docs/policy/rules/#request-rewriting). (see [below for nested schema](#nestedblock--reads--dataset_rewrites))
-   `rate_limit` (Number) Rate Limit specifies the limit of calls that a user can make within a given time period.
-   `severity` (String) severity level that's recorded when someone violate this rule. This is an informational value. Settings: (`low` | `medium` | `high`). If not specified, the severity is considered to be low.

<a id="nestedblock--reads--dataset_rewrites"></a>

### Nested Schema for `reads.dataset_rewrites`

Required:

-   `dataset` (String) The dataset that should be rewritten.In the case of Snowflake, this denotes a fully qualified table name in the form: `<database>.<schema>.<table>`
-   `parameters` (List of String) The set of parameters used in the substitution request, these are references to fields in the activity log as described in the [Additional Checks section](https://cyral.com/docs/policy/rules/#additional-checks).
-   `repo` (String) The name of the repository that the rewrite applies to.
-   `substitution` (String) The request used to substitute references to the dataset.

<a id="nestedblock--updates"></a>

### Nested Schema for `updates`

Required:

-   `data` (List of String) The data locations protected by this rule. Use `*` if you want to define `any` data location. For more information, see the [policy rules](https://cyral.com/docs/policy/rules#contexted-rules) documentation.
-   `rows` (Number) The number of records (for example, rows or documents) that can be accessed/affected in a single statement. Use positive integer numbers to define how many records. If you want to define `any` number of records, set to `-1`.

Optional:

-   `additional_checks` (String) Constraints on the data access specified in [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/). See [Additional checks](https://cyral.com/docs/policy/rules/#additional-checks).
-   `dataset_rewrites` (Block List) Defines how requests should be rewritten in the case of policy violations. See [Request rewriting](https://cyral.com/docs/policy/rules/#request-rewriting). (see [below for nested schema](#nestedblock--updates--dataset_rewrites))
-   `rate_limit` (Number) Rate Limit specifies the limit of calls that a user can make within a given time period.
-   `severity` (String) severity level that's recorded when someone violate this rule. This is an informational value. Settings: (`low` | `medium` | `high`). If not specified, the severity is considered to be low.

<a id="nestedblock--updates--dataset_rewrites"></a>

### Nested Schema for `updates.dataset_rewrites`

Required:

-   `dataset` (String) The dataset that should be rewritten.In the case of Snowflake, this denotes a fully qualified table name in the form: `<database>.<schema>.<table>`
-   `parameters` (List of String) The set of parameters used in the substitution request, these are references to fields in the activity log as described in the [Additional Checks section](https://cyral.com/docs/policy/rules/#additional-checks).
-   `repo` (String) The name of the repository that the rewrite applies to.
-   `substitution` (String) The request used to substitute references to the dataset.
