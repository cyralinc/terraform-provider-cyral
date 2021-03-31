# Policy Rule

CRUD operations for policy rules.


## Usage

```hcl
resource "cyral_policy_rule" "SOME_RESOURCE_NAME" {
    policy_id = ""
    hosts = [""]
    identities {
        db_roles = [""]
        groups = [""]
        services = [""]
        users = [""]
    }
    deletes {
        additional_checks = ""
        data = [""]
        dataset_rewrites {
            dataset = ""
            repo = ""
            substitution = ""
            parameters = [""]
        }
        rows = 1
        severity = "low"
    }
    reads {
        additional_checks = ""
        data = [""]
        dataset_rewrites {
            dataset = ""
            repo = ""
            substitution = ""
            parameters = [""]
        }
        rows = 1
        severity = "low"
    }
    updates {
        additional_checks = ""
        data = [""]
        dataset_rewrites {
            dataset = ""
            repo = ""
            substitution = ""
            parameters = [""]
        }
        rows = 1
        severity = "low"
    }
}
```

## See also

- [Resource Policy](./resource_policy.md)

## Variables

|  Name           |  Default    |  Description                                                                         | Required |
|:----------------|:-----------:|:-------------------------------------------------------------------------------------|:--------:|
|`policy_id`      |             | The ID of the policy you are adding this rule to.                                    | Yes      |
|`identities`     |             | Identities specification that specifies the people, applications, or groups this rule applies to. Every rule except your default rule has one. It can have 4 fields: `db_roles`, `groups`, `users` and `services`. | No |
|`reads`          |             | A contexted rule for accesses of the type `read`.                                    | No       |
|`updates`        |             | A contexted rule for accesses of the type `update`.                                  | No       |
|`deletes`        |             | A contexted rule for accesses of the type `delete`.                                  | No       |
|`hosts`          |             | Hosts specification that limits access to only those users connecting from a certain network location. | No |


> Notes: 
> 1. Unless you create a default rule, users and groups only have the rights you explicitly grant them.  
> 2. Each contexted rule comprises these fields: `data`, `rows`, `severity` `additional_checks`, `dataset_rewrites`. The only required fields are `data` and `rows`.
> 3. The rules block does not need to include all three operation types (reads, updates and deletes); actions you omit are disallowed.
> 4. If you do not include a hosts block, Cyral does not enforce limits based on the connecting client's host address.

For more information, please see the [Policy Guide](https://cyral.com/docs/policy#the-rules-block-of-a-policy).

## Outputs

|  Name        |  Description                                                        |
|:-------------|:--------------------------------------------------------------------|
| `id`         | Unique ID of the resource in the control plane.                     |
