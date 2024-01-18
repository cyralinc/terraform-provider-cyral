---
page_title: "Setup repo-level policy"
---

In this guide, we will attach repo-level data access policies to PostgreSQL and MySQL data
repositories. After reading this guide, you will understand how to setup repo-level policies
and use Cyral's built-in policy templates to control data access.

We recommend that you also read the [Cyral policies](https://cyral.com/docs/policy/overview/)
documentation for more information.

## Repo-level policy types

Cyral offers nine pre-built repo-level policy types. Learn more about them [here](https://cyral.com/docs/policy/repo-level/).
This guide demonstrates creating instances of each type.
Additionally, review the [cyral_rego_policy_instance](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/resources/rego_policy_instance) resource for clarity in parameters.

## Dataset Protection policy

Add a Dataset Protection policy to restrict access to
specific tables or schemas in the data repositories:

-> **Note** The Dataset Protection policy template is only enabled by default in control planes
`v4.13` and later. If you have a previous version, please reach out to our customer success
team to enable it.

### Example Usage

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

## Data Masking policy

Implement a repo-level policy to mask fields for specific users:

### Example Usage

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

Add a repo-level policy to guard against unauthorized updates:

### Example Usage

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

## Data Firewall policy

Set up a repo-level policy to limit which rows users can read from a table:

### Example Usage

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

## User Segmentation policy

Implement a repo-level policy to limit which rows a set of users can read from your database:

### Example Usage

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

## Rate Limit policy

Add a repo-level policy to implement a threshold on sensitive data reads over time:

### Example Usage

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

Implement a repo-level policy to prevent certain records from being read beyond a specified limit:

### Example Usage

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

Set up a repo-level policy to alert when more than a specified number of records are updated or deleted:

### Example Usage

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

Implement a repo-level policy to ensure service accounts can only be used by intended applications:

### Example Usage

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

## Next steps

This guide presents a very simple example of Cyral repo-level policy for each one of the pre-built templates.
Cyral policies have many more capabilities. Check out all parameters that each repo-level policy type supports and use them however you see fit:
[template_parameters in cyral_rego_policy_instance](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/resources/rego_policy_instance#template-parameters).
