terraform {
  required_providers {
    cyral = {
      source  = "local/terraform/cyral"
      version = "~> 4.3"
    }
    okta = {
      source = "okta/okta"
    }
  }
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

# Refer to okta provider documentation:
#
# * https://registry.terraform.io/providers/okta/okta/latest/docs
#
provider "okta" {
  org_name  = "dev-123456" # your organization name
  base_url  = "okta.com"   # your organization url
  api_token = "xxxx"
}

provider "aws" {
  region = "us-east-1"
}

locals {
  # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
  control_plane_host = "[TENANT].app.cyral.com"
  # Set the name to be displayed for this integration in Okta's UI.
  okta_app_name = "Cyral"
  # Set the name to be displayed for this integration in Cyral's UI.
  okta_integration_name = "Okta"
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

    # Set the allowed CIDR block for SSH access to the sidecar
    ssh_inbound_cidr = ["0.0.0.0/0"]
    # Set the allowed CIDR block for database access through the
    # sidecar
    db_inbound_cidr = ["0.0.0.0/0"]
    # Set the allowed CIDR block for monitoring requests to the
    # sidecar
    monitoring_inbound_cidr = ["0.0.0.0/0"]

    # Set the parameters to access the private Cyral container
    # registry.  These parameters can be found on the sidecar
    # Terraform template downloaded from the UI.
    container_registry = {
      name         = "" # see container_registry in the downloaded template
      username     = "" # see container_registry_username in the downloaded template
      registry_key = "" # see container_registry_key in the downloaded template
    }

    # Specify the maximum number of nodes you expect this cluster to
    # have, taking into consideration future growth. This number must be
    # at least equal to the number of nodes currently in your
    # cluster. This number is used for port reservation in the
    # sidecar. This is the value that will be used for the `max_nodes`
    # argument of the `properties` block in the repository resource (see
    # resource `mongodb_repo` below).
    mongodb_max_nodes = 3

    mongodb_ports_low  = 27017
    mongodb_ports_high = local.mongodb_ports_low + local.mongodb_max_nodes

    # All ports that will be used by MongoDB. This range must span at
    # least the `local.mongodb_max_nodes` number of ports. Note that the
    # port number you pass as the second argument to this function is
    # not included in the range. For example, to set port 27021 as your
    # uppermost port number, the second argument must be 27022.
    mongodb_ports = range(local.mongodb_ports_low, local.mongodb_ports_high)
  }
}

# Configures the Okta app and the integration in Cyral control plane.
module "cyral_idp_okta" {
  source  = "cyralinc/idp/okta"
  version = "~> 4.0"

  tenant = "default"

  control_plane = local.control_plane_host

  okta_app_name        = local.okta_app_name
  idp_integration_name = local.okta_integration_name
}
resource "cyral_repository" "mongodb_repo" {
  name = "mongodb_repo"
  type = "mongodb"

  # Specify the address or hostname of the endpoint of at least one node
  # in the MongoDB replica set. You can explictly specify the host and
  # port of additional nodes, or you can mark nodes as 'dynamic' and Cyral
  # will identify the remaining nodes of the replication cluster.
  repo_node {
    name = "node_1"
    host = "mongodb-node1.cyral.com"
    port = 27017
  }

  # You can explictly specify the host and port of additional nodes,
  # or you can mark nodes as 'dynamic' and Cyral will identify the
  # remaining nodes of the replication cluster. However, you will
  # still have to explictly define listeners for each node's port.
  repo_node {
    name    = "node_2"
    dynamic = true
  }

  repo_node {
    name    = "node_3"
    dynamic = true
  }

  mongodb_settings {
    # Specify the replica set identifier, a string value that
    # identifies the MongoDB replica set cluster. To find your
    # replica set ID, see our article:
    #
    # * https://cyral.freshdesk.com/a/solutions/articles/44002241594
    replica_set_name = "some-replica-set"
    server_type      = "replicaset"
  }
}

resource "cyral_repository_conf_auth" "mongodb_repo_auth_config" {
  repository_id     = cyral_repository.mongodb_repo.id
  identity_provider = module.cyral_idp_okta.integration_idp_okta_id
  # Repo TLS is required to allow the sidecar to communicate with
  # MongoDB Atlas.
  repo_tls = "enable"
}

resource "cyral_repository_access_rules" "access_rules" {
  repository_id   = cyral_repository.mongodb_repo.id
  user_account_id = cyral_repository_user_account.mongodb_user_account.user_account_id
  rule {
    identity {
      type = "username"
      name = "me@mycompany.com"
    }
  }
}
# Create listeners for each MongoDB repo node.
resource "cyral_sidecar_listener" "mongodb_listener_node_1" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    host = "mongodb-node1.cyral.com"
    port = 27017
  }
}

resource "cyral_sidecar_listener" "mongodb_listener_node_2" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    host = "mongodb-node1.cyral.com"
    port = 27018
  }
}

resource "cyral_sidecar_listener" "mongodb_listener_node_3" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["mongodb"]
  network_address {
    host = "mongodb-node1.cyral.com"
    port = 27019
  }
}

# Bind the sidecar listeners to the repository.
resource "cyral_repository_binding" "mongodb_repo_binding" {
  repository_id = cyral_repository.mongodb_repo.id
  sidecar_id    = cyral_sidecar.sidecar.id
  enabled       = true
  listener_binding {
    listener_id = cyral_sidecar_listener.mongodb_listener_node_1.listener_id
    node_index  = 0
  }
  listener_binding {
    listener_id = cyral_sidecar_listener.mongodb_listener_node_2.listener_id
    node_index  = 1
  }
  listener_binding {
    listener_id = cyral_sidecar_listener.mongodb_listener_node_3.listener_id
    node_index  = 2
  }
}

# Set the access gateway for the repository.
resource "cyral_repository_access_gateway" "mongodb_access_gateway" {
  repository_id = cyral_repository.mongodb_repo.id
  sidecar_id    = cyral_sidecar.sidecar.id
  binding_id    = cyral_repository_binding.mongodb_repo_binding.binding_id
}


resource "cyral_sidecar" "sidecar" {
  name              = "MongoDBSidecar"
  deployment_method = "terraform"
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
  sidecar_id = cyral_sidecar.sidecar.id
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

output "sidecar_dns" {
  value = module.cyral_sidecar.sidecar_dns
}

output "sidecar_load_balancer_dns" {
  value = module.cyral_sidecar.sidecar_load_balancer_dns
}
