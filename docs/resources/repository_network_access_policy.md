---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cyral_repository_network_access_policy Resource - cyral"
subcategory: ""
description: |-
  Manages the network access policy of a repository. Network access policies are also known as the Network Shield https://cyral.com/docs/manage-repositories/network-shield/. This feature is supported for the following repository types:
    - sqlserver
    - oracle
---

# cyral_repository_network_access_policy (Resource)

Manages the network access policy of a repository. Network access policies are also known as the [Network Shield](https://cyral.com/docs/manage-repositories/network-shield/). This feature is supported for the following repository types:

- `sqlserver`
- `oracle`

## Example Usage

```terraform
# Repository the policy refers to
resource "cyral_repository" "my_sqlserver_repo" {
    name = "my-sqlserver-repo"
    type = "sqlserver"
    host = "sqlserver.cyral.com"
    port = 1433
}

# Allow access from IPs 1.2.3.4 and 4.3.2.1 for Admin database
# account, and from any IP address for accounts Engineer and
# Analyst.
resource "cyral_repository_network_access_policy" "my_sqlserver_repo_policy" {
    repository_id = cyral_repository.my_sqlserver_repo.id
    network_access_rule {
        name = "rule1"
        db_accounts = ["Admin"]
        source_ips = ["1.2.3.4", "4.3.2.1"]
    }
    network_access_rule {
        name = "rule2"
        db_accounts = ["Engineer", "Analyst"]
    }
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

- `repository_id` (String) ID of the repository for which to configure a network access policy.

### Optional

- `enabled` (Boolean) Is the network access policy enabled? Default is true.
- `network_access_rule` (Block Set) Network access policy that decides whether access should be granted based on a set of rules. (see [below for nested schema](#nestedblock--network_access_rule))
- `network_access_rules_block_access` (Boolean) Determines what happens if an incoming connection matches one of the rules in `network_access_rule`. If set to true, the connection is blocked if it matches some rule (and allowed otherwise). Otherwise set to false, the connection is allowed only if it matches some rule. Default is false.

### Read-Only

- `id` (String) ID of this resource in the Cyral environment.

<a id="nestedblock--network_access_rule"></a>

### Nested Schema for `network_access_rule`

Required:

- `name` (String) Name of the rule.

Optional:

- `db_accounts` (List of String) Specify which accounts this rule applies to. The account name must match an existing account in your database.
- `description` (String) Description of the network access policy.
- `source_ips` (List of String) Specify IPs to restrict the range of allowed IP addresses for this rule.