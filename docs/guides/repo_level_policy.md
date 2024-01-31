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

## Example: data firewall

Limit which rows users can read from a table:

```terraform
# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy instance from template to filter table
# 'finance.cards', returning only data where
# finance.cards.country = 'US' for users not in 'Admin' group
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-firewall-policy"
  category    = "SECURITY"
  description = "Returns only data where finance.cards.country = 'US' in table 'finance.cards' for users not in 'Admin' group"
  template_id = "data-firewall"
  parameters  = "{ \"dataSet\": \"finance.cards\", \"dataFilter\": \" finance.cards.country = 'US' \", \"labels\": [\"CCN\"], \"excludedIdentities\": { \"groups\": [\"Admin\"] } }"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
```

## Example: data masking

Mask fields for specific users:

```terraform
# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy instance from template to apply null masking to
# any data labeled as CCN for users in group 'Marketing'
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-masking-policy"
  category    = "SECURITY"
  description = "Apply null masking to any data labeled as CCN for users in group 'Marketing'"
  template_id = "data-masking"
  parameters  = "{ \"maskType\": \"NULL_MASK\", \"labels\": [\"CCN\"], \"identities\": { \"included\": { \"groups\": [\"Marketing\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
```

## Example: data protection

Protect against unauthorized updates:

```terraform
# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy instance from template to raise a 'high' alert
# and block updates and deletes on label CCN
resource "cyral_rego_policy_instance" "policy" {
  name        = "data-protection-policy"
  category    = "SECURITY"
  description = "Raise a 'high' alert and block updates and deletes on label CCN"
  template_id = "data-protection"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorDeletes\": true, \"labels\": [\"CCN\"]}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
```

## Example: rate limit

Set up a threshold on sensitive data reads over time:

```terraform
# Creates pg data repository
resource "cyral_repository" "pg1" {
  type = "postgresql"
  name = "pg-1"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# Creates a policy instance from template to raise a 'high' alert
# and set a rate limit of 500 rows per hour for group 'Marketing'
# and any data labeled as CCN
resource "cyral_rego_policy_instance" "policy" {
  name        = "rate-limit-policy"
  category    = "SECURITY"
  description = "Raise a 'high' alert and set a rate limit of 500 rows per hour for group 'Marketing' and any data labeled as CCN"
  template_id = "rate-limit"
  parameters  = "{ \"rateLimit\": 500, \"block\": true, \"alertSeverity\": \"high\", \"labels\": [\"CCN\"], \"identities\": { \"included\": { \"groups\": [\"Marketing\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
  tags = ["tag1", "tag2"]
}
```

## Example: read limit

Prevent certain records from being read beyond a specified limit:

```terraform
# Creates pg data repository
resource "cyral_repository" "pg1" {
  type = "postgresql"
  name = "pg-1"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# Creates a policy instance from template to limits to 100 the
# amount of rows that can be read per query on the entire
# repository for group 'Devs'
resource "cyral_rego_policy_instance" "policy" {
  name        = "read-limit-policy"
  category    = "SECURITY"
  description = "Limits to 100 the amount of rows that can be read per query on the entire repository for group 'Devs'"
  template_id = "read-limit"
  parameters  = "{ \"rowLimit\": 100, \"block\": true, \"alertSeverity\": \"high\", \"appliesToAllData\": true, \"identities\": { \"included\": { \"groups\": [\"Devs\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
}
```

## Example: repository protection

Alert when more than a specified number of records are updated or deleted:

```terraform
# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy instance from template to limits to 100 the
# amount of rows that can be updated or deleted per query on
# all repository data for anyone except group 'Admin'
resource "cyral_rego_policy_instance" "policy" {
  name        = "repository-protection-policy"
  category    = "SECURITY"
  description = "Limits to 100 the amount of rows that can be updated or deleted per query on all repository data for anyone except group 'Admin'"
  template_id = "repository-protection"
  parameters  = "{ \"rowLimit\": 100, \"block\": true, \"alertSeverity\": \"high\", \"monitorUpdates\": true, \"monitorDeletes\": true, \"identities\": { \"excluded\": { \"groups\": [\"Admin\"] } }}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
}
```

## Example: service account abuse

Ensure service accounts can only be used by intended applications:

```terraform
# Creates pg data repository
resource "cyral_repository" "pg1" {
  type = "postgresql"
  name = "pg-1"

  repo_node {
    host = "pg.cyral.com"
    port = 5432
  }
}

# Creates a policy instance from template to alert and block
# whenever the following service accounts john try to read,
# update, or delete data from the repository without end
# user attribution.
resource "cyral_rego_policy_instance" "policy" {
  name        = "service account abuse policy"
  category    = "SECURITY"
  description = "Alert and block whenever the following service accounts john try to read, update, or delete data from the repository without end user attribution"
  template_id = "service-account-abuse"
  parameters  = "{ \"block\": true, \"alertSeverity\": \"high\", \"serviceAccounts\": [\"john\"]}"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.pg1.id]
  }
}
```

## Example: user segmentation

Limit which rows a set of users can read from your database:

```terraform
# Creates a MySQL data repository named "mysql-1"
resource "cyral_repository" "mysql1" {
  type = "mysql"
  name = "mysql-1"

  repo_node {
    host = "mysql.cyral.com"
    port = 3306
  }
}

# Creates a policy instance from template to filter table
# 'finance.cards' when users in group 'Marketing' read label
# CCN, returning only data where finance.cards.country = 'US'
resource "cyral_rego_policy_instance" "policy" {
  name        = "user-segmentation-policy"
  category    = "SECURITY"
  description = "Filter table 'finance.cards' when users in group 'Marketing' read label CCN, returning only data where finance.cards.country = 'US'"
  template_id = "user-segmentation"
  parameters  = "{ \"dataSet\": \"finance.cards\", \"dataFilter\": \" finance.cards.country = 'US' \", \"labels\": [\"CCN\"], \"includedIdentities\": { \"groups\": [\"Marketing\"] } }"
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
```
