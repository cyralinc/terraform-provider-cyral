# Create the repository
resource "cyral_repository" "example-pg" {
  name = "example-pg"
  type = "postgresql"

  repo_node {
    host = "pg.example.com"
    port = 5432
  }
}

# Create custom labels
resource "cyral_datalabel" "NAME" {
  name        = "NAME"
  description = "Customer name"
  tags        = ["PII"]
}

resource "cyral_datalabel" "DOB" {
  name        = "DOB"
  description = "Customer date of birth"
  tags        = ["PII"]
}

# Create data map for the repository, using the custom labels
resource "cyral_repository_datamap" "example-pg_datamap" {
  repository_id = cyral_repository.example-pg.id

  mapping {
    label = cyral_datalabel.NAME.name
    attributes = [
      "FINANCE.CUSTOMERS.FIRST_NAME",
      "FINANCE.CUSTOMERS.MIDDLE_NAME",
      "FINANCE.CUSTOMERS.LAST_NAME"
    ]
  }

  mapping {
    label = cyral_datalabel.DOB.name
    attributes = [
      "FINANCE.CUSTOMERS.DOB",
    ]
  }
}
