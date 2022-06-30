---
page_title: "Setup policy control over PostgreSQL and MySQL"
---

In this guide, we will set up PostgreSQL and MySQL repositories, bind them to a
sidecar, and attach data access policies to them.

After reading the guide, you will understand how to enforce fine-grained control
over databases using [Cyral policies](https://cyral.com/docs/policy/overview/).

## Sidecar

To set up your sidecar, please follow [this
guide](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/guides/setup_cp_and_deploy_sidecar).
Deploy two `cyral_repository` and `cyral_repository_binding` resources: one for
PostgreSQL and one for MySQL. Also, make sure the Terraform Provider version is
set to at least `2.7.0`, like below:

```terraform
terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = ">= 2.7.0"
    }
  }
}

resource "cyral_repository" "pg_repo" {
...
}

resource "cyral_repository_binding" "pg_repo_binding" {
...
}

resource "cyral_repository" "mysql_repo" {
...
}

resource "cyral_repository_binding" "mysql_repo_binding" {
...
}
```

## Policy

Set up your policies following the example below:

```terraform
locals {
  # PHONE is a predefined label. It exists by default in your control plane.
  phone_label = "PHONE"
}

resource "cyral_datalabel" "custom_label" {
  name        = "CUSTOM_LABEL"
  description = "This is a custom label."
}

resource "cyral_repository_datamap" "pg_datamap" {
  repo_id = cyral_repository.pg_repo.id
  mapping {
    label      = cyral_datalabel.custom_label.name
    attributes = ["customer_schema.table1.col1", "customer_schema.table1.col2"]
  }
}

resource "cyral_repository_datamap" "mysql_datamap" {
  repo_id = cyral_repository.mysql_repo.id
  mapping {
    label      = local.phone_label
    attributes = ["customer_schema.phone.number"]
  }
}

resource "cyral_policy" "customer_data" {
  name        = "customerData"
  data        = [local.phone_label, cyral_datalabel.custom_label.name]
  description = "Control how customer data is handled."
  enabled     = true
  tags        = ["customer"]
}

# To learn more about Cyral policies, see:
#
# * https://cyral.com/docs/policy/overview
#
resource "cyral_policy_rule" "customer_data_rule" {
  policy_id = cyral_policy.customer_data.id

  identities {
    groups = ["client_support", "client_onboarding"]
  }

  # Expect max one entry to be deleted per operation.
  deletes {
    data     = [local.phone_label, cyral_datalabel.custom_label.name]
    rows     = 1
    severity = "high"
  }
  # Expect max one entry updated per operation.
  updates {
    data     = [local.phone_label, cyral_datalabel.custom_label.name]
    rows     = 1
    severity = "high"
  }
  # A query to read more than 100 entries is not normal.
  reads {
    data     = [local.phone_label, cyral_datalabel.custom_label.name]
    rows     = 100
    severity = "medium"
  }
}
```

That's it! If the policies reference the right data locations in your databases,
they will be used to monitor compliance.

## Accessing the data repositories

To learn how to access a repository through the sidecar, see [Connect to a
repository](https://cyral.com/docs/connect/repo-connect).

## Next steps

This guide presents a very simple Cyral policy. Cyral policies have many more
capabilities. Check out all attributes that the policy rule resource supports:
[`cyral_policy_rule`](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/resources/policy_rule).
