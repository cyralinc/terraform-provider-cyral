---
page_title: "Setup repo-level policy"
---

Cyral offers several pre-built [repo-level policy types](https://cyral.com/docs/policy/repo-level/).
In this guide, we provide different examples on how to use them.

Recommended further reading:

- Refer to the [Cyral policies](https://cyral.com/docs/policy/overview/) page in our public
  docs for a complete documentation about the Cyral policy framework.
- Refer to the [`cyral_rego_policy_instance`](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/resources/rego_policy_instance)
  resource for more details about the [template parameters](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/resources/rego_policy_instance#template-parameters)
  and how to use the pre-built repo-level policies in Terraform.

## Data Firewall policy

Limit which rows users can read from a table:

```terraform
# Creates MySQL data repository
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-firewall-policy"
  category    = "SECURITY"
  description = "Filter 'finance.cards' when someone (except 'Admin' group) reads it"
  template_id = "data-firewall"
  parameters  = "{ \"dataSet\": \"finance.cards\", \"dataFilter\": \" finance.cards.country = 'US' \", \"labels\": [\"CCN\"], \"excludedIdentities\": { \"groups\": [\"Admin\"] } }"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
```

## Data Masking policy

Mask fields for specific users:

```terraform
# Creates MySQL data repository
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-masking-policy"
  category    = "SECURITY"
  description = "Masks label CCN for identities in Marketing group"
  template_id = "data-masking"
  parameters  = "{ \"maskType\": \"NULL_MASK\", \"labels\": [\"CCN\"], \"identities\": { \"included\": { \"groups\": [\"Marketing\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
```

## Data Protection policy

Protect against unauthorized updates:

```terraform
# Creates MySQL data repository
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-protection-policy"
  category    = "SECURITY"
  description = "Protect label CCN for update and delete queries"
  template_id = "data-protection"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorDeletes\": true, \"labels\": [\"CCN\"]}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
```

## Dataset Protection policy

-> **Note** The Dataset Protection policy template is only enabled by default in control planes
`v4.13` and later. If you have a previous version, please reach out to our customer success
team to enable it.

Restrict access to specific tables or schemas in the data repositories:

```terraform
# Creates pg data repository
resource "cyral_repository" "repo" {
  type = "postgresql"
  name = "my_pg"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "dataset-protection"
  category    = "SECURITY"
  description = "Blocks reads and updates over schema 'finance' and dataset 'cyral.customers'."
  template_id = "dataset-protection"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorReads\": true, \"datasets\": {\"disallowed\": [\"finance.*\", \"cyral.customers\"]}}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
}
```

## Rate Limit policy

Set up a threshold on sensitive data reads over time:

```terraform
# Creates pg data repository
resource "cyral_repository" "repo" {
  type = "postgresql"
  name = "my_pg"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "rate-limit-policy"
  category    = "SECURITY"
  description = "Implement a threshold on label CCN for group Marketing of 500 rows per hour"
  template_id = "rate-limit"
  parameters  = "{ \"rateLimit\": 500, \"block\": true, \"alertSeverity\": \"high\", \"labels\": [\"CCN\"], \"identities\": { \"included\": { \"groups\": [\"Marketing\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
```

## Read Limit policy

Prevent certain records from being read beyond a specified limit:

```terraform
# Creates pg data repository
resource "cyral_repository" "repo" {
  type = "postgresql"
  name = "my_pg"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "read-limit-policy"
  category    = "SECURITY"
  description = "Limits to 100 the amount of rows that can be read per query on all repository data for group 'Devs'"
  template_id = "read-limit"
  parameters  = "{ \"rowLimit\": 100, \"block\": true, \"alertSeverity\": \"high\", \"appliesToAllData\": true, \"identities\": { \"included\": { \"groups\": [\"Devs\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
}
```

## Repository Protection policy

Alert when more than a specified number of records are updated or deleted:

```terraform
# Creates MySQL data repository
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
    host = "mysql.cyral.com"
    port = 5432
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "repository-protection-policy"
  category    = "SECURITY"
  description = "Limits to 100 the amount of rows that can be updated or deleted per query on all repository data for anyone except group 'Admin'"
  template_id = "repository-protection"
  parameters  = "{ \"rowLimit\": 100, \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorDeletes\": true, \"identities\": { \"excluded\": { \"groups\": [\"Admin\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
}
```

## Service Account Abuse policy

Ensure service accounts can only be used by intended applications:

```terraform
# Creates pg data repository
resource "cyral_repository" "repo" {
  type = "postgresql"
  name = "my_pg"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "service account abuse policy"
  category    = "SECURITY"
  description = "Always require user attribution for service acount 'john'"
  template_id = "service-account-abuse"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"serviceAccounts\": [\"john\"]}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
}
```

## User Segmentation policy

Limit which rows a set of users can read from your database:

```terraform
# Creates MySQL data repository
resource "cyral_repository" "repo" {
  type = "mysql"
  name = "my_mysql"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# create policy instance from template
resource "cyral_rego_policy_instance" "policy" {
  name        = "user-segmentation-policy"
  category    = "SECURITY"
  description = "Applies a data filter in 'finance.cards' when someone from group 'Marketing' reads data labeled as 'CCN'"
  template_id = "user-segmentation"
  parameters  = "{ \"dataSet\": \"finance.cards\", \"dataFilter\": \" finance.cards.country = 'US' \", \"labels\": [\"CCN\"], \"includedIdentities\": { \"groups\": [\"Marketing\"] } }"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.repo.id]
  }
  tags = ["tag1", "tag2"]
}
```
