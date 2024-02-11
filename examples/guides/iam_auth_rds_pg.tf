terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = "~> 4.7"
    }
  }
}

locals {
  # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
  control_plane_host = "[TENANT].app.cyral.com"

  # Use the name of the IdP that will be used to access the RDS instance
  idp = {
    name = "<IDP_NAME_AS_SHOWN_IN_THE_UI>"
  }

  repos = {
    pg = {
      host = "<RDS_INSTANCE_ADDRESS>"
      port = 5432
    }
  }

  sidecar = {
    # Set to true if you want a sidecar deployed with an
    # internet-facing load balancer (requires a public subnet).
    public_sidecar = false

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

    # Optionally set the hosted zone ID that will be used to create the
    # DNS name in parameter `dns_name`
    dns_hosted_zone_id = ""
    # Optionally set the DNS name that will be used by your sidecar. Ex:
    # sidecar.mycompany.com
    dns_name = ""
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

# The log group is created in AWS by module.cyral_sidecar
# when the sidecar is deployed.
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
  activity_log_integration_id = cyral_integration_logging.cloudwatch.id
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
  sidecar_id = cyral_sidecar.sidecar.id
}

resource "cyral_repository" "pg" {
  name = "pgRepo"
  type = "postgresql"

  repo_node {
    host = local.repos.pg.host
    port = local.repos.pg.port
  }
}

resource "cyral_sidecar_listener" "pg" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["postgresql"]
  network_address {
    port = local.repos.pg.port
  }
}

resource "cyral_repository_binding" "pg" {
  sidecar_id    = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.pg.id
  listener_binding {
    listener_id = cyral_sidecar_listener.pg.listener_id
  }
}

data "cyral_integration_idp_saml" "saml" {
  display_name = local.idp.name
}

# Let users from the provided `identity_provider` use SSO
# to access the database
resource "cyral_repository_conf_auth" "pg" {
  repository_id = cyral_repository.pg.id

  client_tls = "enable"
  repo_tls   = "enable"

  identity_provider = data.cyral_integration_idp_saml.saml.idp_list[0].id
}

# Enables the access portal for this repository in the
# especified sidecar
resource "cyral_repository_access_gateway" "pg" {
  repository_id = cyral_repository.pg.id
  sidecar_id    = cyral_sidecar.sidecar.id
  binding_id    = cyral_repository_binding.pg.binding_id
}

###########################################################################
# Creates an IAM policy that the sidecar will assume in order to access
# the RDS instance. In this example, the policy attached to the role will
# let the sidecar connect to all databases in all available accounts and
# regions.
#
# This should NOT be used in production. Refer to the AWS documentation
# for guidance on how to restrict to the database you plan to protect.
#
data "aws_iam_policy_document" "rds_access_policy" {
  statement {
    actions   = ["rds-db:connect"]
    resources = [
      "*"
    ]
  }
}

resource "aws_iam_policy" "rds_access_policy" {
  name        = "my-sidecar_access_policy"
  path        = "/"
  description = "Allow sidecar to connect to all RDS instances"
  policy      = data.aws_iam_policy_document.rds_access_policy.json
}

data "aws_iam_policy_document" "sidecar_trust_policy" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "AWS"
      identifiers = [module.cyral_sidecar.aws_iam_role_arn]
    }
  }
}

resource "aws_iam_role" "rds_role" {
  name               = "my-sidecar_rds_access_role"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.sidecar_trust_policy.json
}

resource "aws_iam_role_policy_attachment" "rds_role_policy_attachment" {
  role       = aws_iam_role.rds_role.name
  policy_arn = aws_iam_policy.rds_access_policy.arn
}
###########################################################################

resource "cyral_repository_user_account" "pg_repo_user_account" {
  # You may opt for a better name here as this is the name that will
  # be shown in the UI
  name = "my-sidecar_rds_access_role"
  repository_id = cyral_repository.pg.id
  auth_scheme {
    aws_iam {
      role_arn = aws_iam_role.rds_role.arn
    }
  }
}

# Set the proper identity for the username, email or group that will
# be allowed to access the PG database using SSO
resource "cyral_repository_access_rules" "access_rule" {
  repository_id = cyral_repository.pg.id
  user_account_id = cyral_repository_user_account.pg_repo_user_account.user_account_id
  rule {
    identity {
      type = "email"
      name = "myuser@mycompany.com"
    }
  }
}

module "cyral_sidecar" {
  source = "cyralinc/sidecar-ec2/aws"

  # Use the module version that is compatible with your sidecar.
  version = "~> 4.3"

  sidecar_id = cyral_sidecar.sidecar.id
  control_plane = local.control_plane_host
  client_id     = cyral_sidecar_credentials.sidecar_credentials.client_id
  client_secret = cyral_sidecar_credentials.sidecar_credentials.client_secret

  sidecar_ports = [local.repos.pg.port]

  vpc_id  = local.sidecar.vpc_id
  subnets = local.sidecar.subnets

  ssh_inbound_cidr        = local.sidecar.ssh_inbound_cidr
  db_inbound_cidr         = local.sidecar.db_inbound_cidr
  monitoring_inbound_cidr = local.sidecar.monitoring_inbound_cidr

  load_balancer_scheme        = local.sidecar.public_sidecar ? "internet-facing" : "internal"
  associate_public_ip_address = local.sidecar.public_sidecar

  sidecar_dns_hosted_zone_id = local.sidecar.dns_hosted_zone_id
  sidecar_dns_name           = local.sidecar.dns_name
}

output "sidecar_load_balancer_dns" {
  value = module.cyral_sidecar.sidecar_load_balancer_dns
}
