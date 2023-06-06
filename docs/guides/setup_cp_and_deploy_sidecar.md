---
page_title: "Create an AWS EC2 sidecar to protect PostgreSQL and MySQL databases"
---

Use the following guide to create the minimum required configuration in both Cyral
Control Plane and your AWS account to deploy a Cyral Sidecar to AWS EC2 in front
of two database instances: MySQL and PostgreSQL.

By running this example you will have a fully functional sidecar on your AWS
account. Read the comments and update the necessary parameters as instructed.

See also the [Cyral Sidecar module for AWS EC2](https://registry.terraform.io/modules/cyralinc/sidecar-ec2/aws/latest)
for more details on how the sidecar is deployed to AWS and more advanced configurations.

```terraform
terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = "~> 4.0"
    }
  }
}

locals {
    # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
    control_plane_host = "[TENANT].app.cyral.com"

    repos = {
        postgresql = {
            host  = "your-postgres-db-host"
            # This is the port the DATABASE accepts connections.
            db_port = 5432
            # This is the port the SIDECAR will expose to
            # clients connecting to this DB. In this case,
            # it is different only for education purposes.
            sidecar_port = 5433
            type  = "postgresql"
        }
        mysql = {
            host  = "your-mysql-db-host"
            db_port = 3306
            sidecar_port = 3307
            type  = "mysql"
        }
    }

    sidecar = {
        # Set to true if you want a sidecar deployed with an
        # internet-facing load balancer (requires a public subnet).
        public_sidecar = false

        # Set the desired sidecar version.
        sidecar_version = "v4.7.0"

        # Set the AWS region that the sidecar will be deployed to
        region = ""
        # Set the ID of VPC that the sidecar will be deployed to
        vpc_id = ""
        # Set the IDs of the subnets that the sidecar will be deployed to
        subnets = [""]
        # Name of the CloudWatch log group used to push logs
        cloudwatch_log_group_name = "cyral-example-loggroup"

        # Set the allowed CIDR block for SSH access to the sidecar
        ssh_inbound_cidr = ["0.0.0.0/0"]
        # Set the allowed CIDR block for database access through the
        # sidecar
        db_inbound_cidr = ["0.0.0.0/0"]
        # Set the allowed CIDR block for monitoring requests to the
        # sidecar
        monitoring_inbound_cidr = ["0.0.0.0/0"]

        # Set the parameters to access the private Cyral container
        # registry. These parameters can be found on the sidecar
        # Terraform template downloaded from the UI.
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

# Follow the instructions in the Cyral Terraform Provider page to set
# up the credentials:
#
# * https://registry.terraform.io/providers/cyralinc/cyral/latest/docs
provider "cyral" {
    client_id     = ""
    client_secret = ""

    control_plane = local.control_plane_host
}

resource "cyral_integration_logging" "cloudwatch" {
    name = "my-cloudwatch"
    cloudwatch {
        region = local.sidecar.region
        group  = local.sidecar.cloudwatch_log_group_name
        stream = "cyral-sidecar"
    }
}

resource "cyral_sidecar" "sidecar" {
    name               = "my-sidecar"
    deployment_method  = "terraform"
    log_integration_id = cyral_integration_logging.cloudwatch.id
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
    sidecar_id = cyral_sidecar.sidecar.id
}

resource "cyral_repository" "all_repositories" {
    for_each = local.repos
    name  = each.key
    type  = each.value.type

    connection_draining {
        auto      = false
        wait_time = 0
    }

    repo_node {
        host = each.value.host
        port = each.value.db_port
    }
}

resource "cyral_sidecar_listener" "all_listeners" {
    for_each = local.repos
    sidecar_id = cyral_sidecar.sidecar.id
    repo_types = [each.value.type]
    network_address {
        port = each.value.sidecar_port
    }
}

resource "cyral_repository_binding" "all_repo_binding" {
    for_each = local.repos
    repository_id = cyral_repository.all_repositories[each.key].id
    sidecar_id = cyral_sidecar.sidecar.id

    listener_binding {
        listener_id = cyral_sidecar_listener.all_listeners["${each.value.type}"].listener_id
    }
}

module "cyral_sidecar" {
    source = "cyralinc/sidecar-ec2/aws"

    # Use the module version that is compatible with your sidecar.
    version = "~> 4.0"

    sidecar_version = local.sidecar.sidecar_version

    sidecar_id = cyral_sidecar.sidecar.id

    control_plane = local.control_plane_host

    cloudwatch_log_group_name = local.sidecar.cloudwatch_log_group_name

    sidecar_ports = [for repo in values(local.repos) : repo.sidecar_port]

    vpc_id          = local.sidecar.vpc_id
    subnets         = local.sidecar.subnets

    ssh_inbound_cidr        = local.sidecar.ssh_inbound_cidr
    db_inbound_cidr         = local.sidecar.db_inbound_cidr
    monitoring_inbound_cidr = local.sidecar.monitoring_inbound_cidr

    load_balancer_scheme        = local.sidecar.public_sidecar ? "internet-facing" : "internal"
    associate_public_ip_address = local.sidecar.public_sidecar

    deploy_secrets   = true
    secrets_location = "/cyral/sidecars/${cyral_sidecar.sidecar.id}/secrets"

    container_registry          = local.sidecar.container_registry.name
    container_registry_username = local.sidecar.container_registry.username
    container_registry_key      = local.sidecar.container_registry.registry_key

    client_id     = cyral_sidecar_credentials.sidecar_credentials.client_id
    client_secret = cyral_sidecar_credentials.sidecar_credentials.client_secret
}

output "sidecar_dns" {
    value = module.cyral_sidecar.sidecar_dns
}

output "sidecar_load_balancer_dns" {
    value = module.cyral_sidecar.sidecar_load_balancer_dns
}
```

## Accessing the data repositories

To learn how to access a repository through the sidecar, see [Connect to a
repository](https://cyral.com/docs/connect/repo-connect).

## Enforcing access policies

To attach access policies to the created data repositories, please follow the
guide [Setup policy control over PostgreSQL and MySQL](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/guides/pg_mysql_sidecar_policy).
