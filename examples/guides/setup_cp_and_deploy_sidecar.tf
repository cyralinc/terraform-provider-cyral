terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = "~> 3.0"
    }
  }
}

locals {
    # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
    control_plane_host = "[TENANT].app.cyral.com"
    # Set the control plane API port
    control_plane_port = 443

    sidecar = {
        # If you would like to use other log integration, download a
        # new template from the UI and copy the log integration
        # configuration or follow this module documentation.
        log_integration = "cloudwatch"
        # Set the desired EC2 instance type for the auto scaling
        # group.
        instance_type = "t3.medium"
        # Set to true if you want a sidecar deployed with an
        # internet-facing load balancer (requires a public subnet).
        public_sidecar = false

        # Set the AWS region that the sidecar will be deployed to
        region = ""
        # Set the ID of VPC that the sidecar will be deployed to
        vpc_id = ""
        # Set the IDs of the subnets that the sidecar will be deployed to
        subnets = [""]

        # Set the allowed CIDR block for SSH access to the sidecar
        ssh_inbound_cidr = ["0.0.0.0/0"]
        # Set the allowed CIDR block for database access through the
        # sidecar
        db_inbound_cidr = ["0.0.0.0/0"]
        # Set the allowed CIDR block for health check requests to the
        # sidecar
        healthcheck_inbound_cidr = ["0.0.0.0/0"]

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

    control_plane = "${local.control_plane_host}:${local.control_plane_port}"
}

resource "cyral_sidecar" "main_sidecar" {
  name              = "MainSidecar"
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
  # Set the desired sidecar version. This information can be extracted
  # from the template downloaded from the UI.
  sidecar_version = "v3.0.0"

  source = "cyralinc/sidecar-ec2/aws"
  # Use the module version that is compatible with your sidecar. This
  # information can be extracted from the template downloaded from the
  # UI.
  version = "~> 3.0"

  sidecar_id = cyral_sidecar.main_sidecar.id

  control_plane = local.control_plane_host

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

  client_id     = cyral_sidecar_credentials.sidecar_credentials.client_id
  client_secret = cyral_sidecar_credentials.sidecar_credentials.client_secret
}

output "sidecar_dns" {
  value = module.cyral_sidecar.sidecar_dns
}

output "sidecar_load_balancer_dns" {
  value = module.cyral_sidecar.sidecar_load_balancer_dns
}
