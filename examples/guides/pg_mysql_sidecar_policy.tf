locals {
  # PHONE is a predefined label. It exists by default in your control plane.
  phone_label = "PHONE"
}

resource "cyral_datalabel" "custom_label" {
  name        = "CUSTOM_LABEL"
  description = "This is a custom label."
}

resource "cyral_repository_datamap" "pg_datamap" {
  repository_id = cyral_repository.pg_repo.id
  mapping {
    label      = cyral_datalabel.custom_label.name
    attributes = ["customer_schema.table1.col1", "customer_schema.table1.col2"]
  }
}

resource "cyral_repository_datamap" "mysql_datamap" {
  repository_id = cyral_repository.mysql_repo.id
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
  # A query to read more than 100 entries is not considered normal.
  reads {
    data     = [local.phone_label, cyral_datalabel.custom_label.name]
    rows     = 100
    severity = "medium"
  }
}
