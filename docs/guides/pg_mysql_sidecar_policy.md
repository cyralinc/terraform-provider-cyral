---
page_title: "Setup policy control over PostgreSQL and MySQL"
---

In this guide, we will set up PostgreSQL and MySQL repositories, bind them to a
sidecar, and attach data access policies to them.

In the end, it should have become clear how to set up fine-grained control over
databases using [Cyral policies](https://cyral.com/docs/policy/overview/).

## Sidecar

First, let's set up the sidecar and the data repositories. Follow the code
below, replacing variable values according to your needs. This set up is
minimalistic: we will just use the necessary Cyral features to get your sidecar
up and running.

```terraform
terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = ">= 2.7.0"
    }
  }
}

locals {
  # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
  control_plane = "[TENANT].dev.cyral.com"

  sidecar = {
    # The sidecar name prefix is a unique name that is used by the sidecar
    # module to create resources in the target AWS account. For this reason, we
    # use 'cyral-zzzzzz' where zzzzzz are the last 6 digits of the sidecar id
    # created in the control plane. This explanation is purely informational and
    # we advise you to keep this variable as is.
    sidecar_name_prefix = "cyral-${substr(lower(cyral_sidecar.main_sidecar.id), -6, -1)}"

    # If you would like to use other log integration, download a new template
    # from the UI and copy the log integration configuration or follow this
    # module documentation.
    log_integration = "cloudwatch"
    # Set the desired EC2 instance type for the auto scaling group.
    instance_type = "t3.medium"
    # Set to true if you want a sidecar deployed with an internet-facing load
    # balancer (requires a public subnet).
    public_sidecar = false

    # Set the AWS region that the sidecar will be deployed to
    region = ""
    # Set the ID of VPC that the sidecar will be deployed to
    vpc_id = ""
    # Set the IDs of the subnets that the sidecar will be deployed to
    subnets = [""]

    # Set the allowed CIDR block for SSH access to the sidecar
    ssh_inbound_cidr = ["0.0.0.0/0"]
    # Set the allowed CIDR block for database access through the sidecar
    db_inbound_cidr = ["0.0.0.0/0"]
    # Set the allowed CIDR block for health check requests to the sidecar
    healthcheck_inbound_cidr = ["0.0.0.0/0"]

    # Set the parameters to access the private Cyral container registry.  These
    # parameters can be found on the sidecar Terraform template downloaded from
    # the UI.
    container_registry = {
      name         = "" # see container_registry in the downloaded template
      username     = "" # see container_registry_username in the downloaded template
      registry_key = "" # see container_registry_key in the downloaded template
    }
  }
}

provider "aws" {
    region = local.sidecar.region
}

provider "cyral" {
  # Follow the instructions in the Cyral Terraform Provider page to set up the
  # credentials:
  #
  # * https://registry.terraform.io/providers/cyralinc/cyral/latest/docs
  client_id     = ""
  client_secret = ""
  control_plane = "${local.control_plane}:8000"
}

resource "cyral_sidecar" "main_sidecar" {
  name              = "main_sidecar"
  deployment_method = "terraform"
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
  sidecar_id = cyral_sidecar.main_sidecar.id
}

resource "cyral_repository" "pg_repo" {
  name = "pg_repo"
  type = "postgresql"
  host = "postgresql.mycompany.com"
  port = 5432
}

resource "cyral_repository_binding" "pg_repo_binding" {
  repository_id = cyral_repository.pg_repo.id
  sidecar_id    = cyral_sidecar.main_sidecar.id
  listener_port = 5432
}

resource "cyral_repository" "mysql_repo" {
  name = "mysql_repo"
  type = "mysql"
  host = "mysql.mycompany.com"
  port = 3306
}

resource "cyral_repository_binding" "mysql_repo_binding" {
  repository_id = cyral_repository.mysql_repo.id
  sidecar_id    = cyral_sidecar.main_sidecar.id
  listener_port = 3306
}

module "cyral_sidecar" {
  # Set the desired sidecar version. This information can be extracted from
  # the template downloaded from the UI.
  sidecar_version = "v2.32.0"

  source = "cyralinc/sidecar-aws/cyral"
  # Use the module version that is compatible with your sidecar. This
  # information can be extracted from the template downloaded from the UI.
  version = "2.8.1"

  sidecar_id = cyral_sidecar.main_sidecar.id

  name_prefix = local.sidecar.sidecar_name_prefix

  control_plane = local.control_plane

  repositories_supported = ["postgresql", "mysql"]

  sidecar_ports = [cyral_repository.pg_repo.port, cyral_repository.mysql_repo.port]

  mongodb_port_alloc_range_low  = 0
  mongodb_port_alloc_range_high = 0

  instance_type   = local.sidecar.instance_type
  log_integration = local.sidecar.log_integration
  vpc_id          = local.sidecar.vpc_id
  subnets         = local.sidecar.subnets

  ssh_inbound_cidr         = local.sidecar.ssh_inbound_cidr
  db_inbound_cidr          = local.sidecar.db_inbound_cidr
  healthcheck_inbound_cidr = local.sidecar.healthcheck_inbound_cidr

  load_balancer_scheme        = local.sidecar.public_sidecar ? "internet-facing" : "internal"
  associate_public_ip_address = local.sidecar.public_sidecar

  deploy_secrets   = true
  secrets_location = "/cyral/sidecars/${cyral_sidecar.main_sidecar.id}/secrets"

  container_registry          = local.sidecar.container_registry.name
  container_registry_username = local.sidecar.container_registry.username
  container_registry_key      = local.sidecar.container_registry.registry_key
  client_id                   = cyral_sidecar_credentials.sidecar_credentials.client_id
  client_secret               = cyral_sidecar_credentials.sidecar_credentials.client_secret
}

output "sidecar_dns" {
  value = module.cyral_sidecar.sidecar_dns
}

output "sidecar_load_balancer_dns" {
  value = module.cyral_sidecar.sidecar_load_balancer_dns
}
```

## Policy

Set up your policies following the example code below:

```terraform
locals {
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

## Accessing the data repositories

To test your access to the repositories through the sidecar, check out [Connect
to a repository](https://cyral.com/docs/connect/repo-connect).
