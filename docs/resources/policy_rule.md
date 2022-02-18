# Policy Rule Resource

Provides a resource to handle [policy rules](https://cyral.com/docs/reference/policy/#rules). See also: [Policy](./policy.md)

## Example Usage

```hcl
resource "cyral_policy_rule" "some_resource_name" {
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

## Argument Reference

- `policy_id` - (Required) The ID of the policy you are adding this rule to.
- `identities` - (Optional) Identities specification that specifies the people, applications, or groups this rule applies to. Every rule except your default rule has one. It can have 4 fields: `db_roles`, `groups`, `users` and `services`.
- `reads` - (Optional) A contexted rule for accesses of the type `read`.
- `updates` - (Optional) A contexted rule for accesses of the type `update`.
- `deletes` - (Optional) A contexted rule for accesses of the type `delete`.
- `hosts` - (Optional) Hosts specification that limits access to only those users connecting from a certain network location.

> Notes:
>
> 1. Unless you create a default rule, users and groups only have the rights you explicitly grant them.
> 2. Each contexted rule comprises these fields: `data`, `rows`, `severity` `additional_checks`, `dataset_rewrites`. The only required fields are `data` and `rows`.
> 3. The rules block does not need to include all three operation types (reads, updates and deletes); actions you omit are disallowed.
> 4. If you do not include a hosts block, Cyral does not enforce limits based on the connecting client's host address.

For more information, see the [Policy Guide](https://cyral.com/docs/policy#the-rules-block-of-a-policy).

## Attribute Reference

- `id` - The ID of this resource.
