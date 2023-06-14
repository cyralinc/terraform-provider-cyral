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

  listeners = {
    mysql = {
      # This is the port the SIDECAR will expose to
      # clients connecting to all `mysql` databases.
      sidecar_port = 3306
    }
  }
  repos = {
    mysql1 = {
      # Name that will be shown in the Cyral UI
      name = "mysql-1"
      host = "your-mysql-1-db-host"
      # This is the port the DATABASE accepts connections.
      db_port = 3309
      type    = "mysql"
      database_credentials = {
        # Set the credentials that should be used to connect to the database
        username = ""
        password = ""
      }
    }
    mysql2 = {
      # Name that will be shown in the Cyral UI
      name = "mysql-2"
      host = "your-mysql-2-db-host"
      # This is the port the DATABASE accepts connections.
      db_port = 3310
      type    = "mysql"
      database_credentials = {
        # Set the credentials that should be used to connect to the database
        username = ""
        password = ""
      }
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
    # registry. These parameters can be found in the sidecar
    # Terraform template downloaded from the UI. Use the
    # commented values to locate the variables and copy the
    # values from the downloaded template.
    container_registry = {
      name         = "" # container_registry
      username     = "" # container_registry_username
      registry_key = "" # container_registry_key
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
  log_integration_id = cyral_integration_logging.cloudwatch.id
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
  sidecar_id = cyral_sidecar.sidecar.id
}

resource "cyral_repository" "all_repositories" {
  for_each = local.repos
  name     = each.value.name
  type     = each.value.type

  repo_node {
    host = each.value.host
    port = each.value.db_port
  }
}

resource "cyral_sidecar_listener" "all_listeners" {
    for_each = local.listeners
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = [each.key]
  // Clients will connect to both MySQL repos through
  // the same port
  network_address {
    port = each.value.sidecar_port
  }
  # MySQL version that will be shown to clients
  # connecting to both MySQL instances
  mysql_settings {
    db_version = "8.0.4"
  }
}

resource "cyral_repository_binding" "bindings" {
  for_each      = local.repos
  repository_id = cyral_repository.all_repositories[each.key].id
  sidecar_id    = cyral_sidecar.sidecar.id
  # Smart ports will be automatically be activated as both
  # repos are bound to the same listener
  listener_binding {
    listener_id = cyral_sidecar_listener.all_listeners[each.value.type].listener_id
  }
}

#####################################################################
# Deploys the credentials that the sidecar will use to access the
# databases and associate them to the repositories as user accounts
resource "aws_secretsmanager_secret" "all_secrets" {
  for_each = local.repos
  # The sidecar deployed using our AWS sidecar module has access to
  # all secrets with the prefix '/cyral/' in the region it is
  # deployed.
  name = join("", [
    "/cyral/dbsecrets/",
    cyral_repository.all_repositories[each.key].id
  ])
}

resource "aws_secretsmanager_secret_version" "all_secret_versions" {
  for_each      = local.repos
  secret_id     = aws_secretsmanager_secret.all_secrets[each.key].id
  secret_string = jsonencode(each.value.database_credentials)
}

resource "cyral_repository_user_account" "all_user_accounts" {
  for_each      = local.repos
  repository_id = cyral_repository.all_repositories[each.key].id
  name          = each.value.database_credentials.username
  auth_scheme {
    aws_secrets_manager {
      secret_arn = aws_secretsmanager_secret.all_secrets[each.key].arn
    }
  }
}
#####################################################################

data "cyral_integration_idp_saml" "saml" {
  display_name = "<IDP_NAME_AS_SHOWN_IN_THE_UI>"
}

# Allow users from SSO group `Everyone` access the database
resource "cyral_repository_access_rules" "all_access_rules" {
  for_each        = local.repos
  repository_id   = cyral_repository.all_repositories[each.key].id
  user_account_id = cyral_repository_user_account.all_user_accounts[each.key].user_account_id
  rule {
    identity {
      type = "group"
      name = "Everyone"
    }
  }
}

# Let users from the provided `identity_provider` use SSO
# to access the database
resource "cyral_repository_conf_auth" "all_conf_auth" {
  for_each          = local.repos
  repository_id     = cyral_repository.all_repositories[each.key].id
  identity_provider = "cyral-engineering"
  #identity_provider = data.cyral_integration_idp_saml.saml.idp_list[0].id
  allow_native_auth = true
}

# Enables the access portal for this repository in the
# especified sidecar
resource "cyral_repository_access_gateway" "all_access_gateways" {
  for_each      = local.repos
  repository_id = cyral_repository.all_repositories[each.key].id
  sidecar_id    = cyral_sidecar.sidecar.id
  binding_id    = cyral_repository_binding.bindings[each.key].binding_id
}

module "cyral_sidecar" {
  source = "cyralinc/sidecar-ec2/aws"

  # Use the module version that is compatible with your sidecar.
  version = "~> 4.0"

  sidecar_version = local.sidecar.sidecar_version

  sidecar_id = cyral_sidecar.sidecar.id

  control_plane = local.control_plane_host

  sidecar_ports = [for l in values(local.listeners) : l.sidecar_port]

  vpc_id  = local.sidecar.vpc_id
  subnets = local.sidecar.subnets

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

output "sidecar_load_balancer_dns" {
  value = module.cyral_sidecar.sidecar_load_balancer_dns
}
