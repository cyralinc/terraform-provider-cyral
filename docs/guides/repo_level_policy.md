---
page_title: "Setup repo-level policies"
---

Cyral offers several [policy wizards](https://cyral.com/docs/policy/repo-level/).
These wizards generate policies for common use cases based on the parameters you provide. The created policies are part of a _policy set_.
This guide shows how to define policy sets that use these wizards to create policies in Terraform.

Recommended further reading:

-   Refer to the [Cyral policies](https://cyral.com/docs/policy/overview/) page in our public
    docs for a complete documentation about the Cyral policy framework.
-   Refer to the [`cyral_policy_set`](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/resources/policy_set)
    resource for more details about how to create policy sets in Terraform.

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

# Creates a policy set using the data firewall wizard to filter table
# 'finance.cards', returning only data where
# finance.cards.country = 'US' for users not in 'Admin' group
resource "cyral_policy_set" "data_firewall_policy" {
  name        = "data firewall policy"
  description = "Returns only data where finance.cards.country = 'US' in table 'finance.cards' for users not in 'Admin' group"
  wizard_id   = "data-firewall"
  parameters  = jsonencode(
    {
      "dataset" = "finance.cards"
      "dataFilter" = " finance.cards.country = 'US' "
      "labels" = ["CCN"]
      "excludedIdentities" = { "groups" = ["Admin"] }
    }
  )
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

# Creates a policy set using the data masking wizard to apply null masking to
# any data labeled as CCN for users in group 'Marketing'
resource "cyral_policy_set" "data_masking_policy" {
  name        = "data masking policy"
  description = "Apply null masking to any data labeled as CCN for users in group 'Marketing'"
  wizard_id   = "data-masking"
  parameters  = jsonencode(
    {
      "maskType" = "null"
      "labels" = ["CCN"]
      "identities" = { "included": { "groups" = ["Marketing"] } }
    }
  )
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

# Creates a policy set using the data protection wizard to raise
# an alert and block updates and deletes on label CCN
resource "cyral_policy_set" "data_protection_policy" {
  name        = "data protection policy"
  description = "Raise an alert and block updates and deletes on label CCN"
  wizard_id   = "data-protection"
  parameters  = jsonencode(
    {
      "block" = true
      "alertSeverity" = "high"
      "governedOperations" = ["update", "delete"]
      "labels" = ["CCN"]
    }
  )
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

# Creates a policy set using the rate limit wizard to raise an alert
# and set a rate limit of 500 rows per hour for group 'Marketing'
# and any data labeled as CCN
resource "cyral_policy_set" "rate_limit_policy" {
  name        = "rate limit policy"
  description = "Raise an alert and set a rate limit of 500 rows per hour for group 'Marketing' and any data labeled as CCN"
  wizard_id   = "rate-limit"
  parameters  = jsonencode(
    {
      "rateLimit" = 500
      "enforce" = true
      "labels" = ["CCN"]
      "identities" = { "included": { "groups" = ["Marketing"] } }
    }
  )
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

# Creates a policy set using the read limit wizard to limits to 100 the
# amount of rows that can be read per query on the entire
# repository for group 'Devs'
resource "cyral_policy_set" "read_limit_policy" {
  name        = "read limit policy"
  description = "Limits to 100 the amount of rows that can be read per query on the entire repository for group 'Devs'"
  wizard_id   = "read-limit"
  parameters  = jsonencode(
    {
      "rowLimit" = 100
      "enforce" = true
      "datasets" = "*"
      "identities" = { "included": { "groups" = ["Devs"] } }
    }
  )
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

# Creates a policy set using the repository protection wizard to alert if more than
# 100 rows are updated or deleted per query on all repository data by anyone except group 'Admin'
resource "cyral_policy_set" "repository_protection_policy" {
  name        = "repository protection policy"
  description = "Alert if more than 100 rows are updated or deleted per query on all repository data by anyone except group 'Admin'"
  wizard_id   = "repository-protection"
  parameters  = jsonencode(
    {
      "rowLimit" = 100
      "datasets" = "*"
      "governedOperations" = ["update", "delete"]
      "identities" = { "excluded": { "groups" = ["Admin"] } }
    }
  )
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

# Creates a policy set using the service account abuse wizard to alert and block
# whenever the service accounts john is used without end user attribution.
resource "cyral_policy_set" "service_account_abuse_policy" {
  name        = "service account abuse policy"
  description = "Alert and block whenever the service accounts john is used without end user attribution"
  wizard_id   = "service-account-abuse"
  parameters  = jsonencode(
    {
      "block" = true
      "alertSeverity" = "high"
      "serviceAccounts" = ["john"]
    }
  )
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

# Creates a policy set using the user segmentation wizard to filter table
# 'finance.cards' when users in group 'Marketing' read label
# CCN, returning only data where finance.cards.country = 'US'
resource "cyral_policy_set" "user_segmentation_policy" {
  name        = "user segmentation policy"
  description = "Filter table 'finance.cards' when users in group 'Marketing' read label CCN, returning only data where finance.cards.country = 'US'"
  wizard_id   = "user-segmentation"
  parameters  = jsonencode(
    {
      "dataset" = "finance.cards"
      "dataFilter" = " finance.cards.country = 'US' "
      "labels" = ["CCN"]
      "includedIdentities" = { "groups" = ["Marketing"] }
    }
  )
  enabled     = true
  scope {
    repo_ids = [cyral_repository.mysql1.id]
  }
  tags = ["tag1", "tag2"]
}
```
