locals {
  # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
  control_plane = "[TENANT].app.cyral.com"

  sidecar = {
    # The sidecar name prefix is a unique name that is used by the sidecar
    # module to create resources in the target AWS account. For this reason, we
    # use 'cyral-zzzzzz' where zzzzzz are the last 6 digits of the sidecar id
    # created in the control plane. This explanation is purely informational and
    # we advise you to keep this variable as is.
    sidecar_name_prefix = "cyral-${substr(lower(cyral_sidecar.mongodb_sidecar.id), -6, -1)}"

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

  # Specify the maximum number of nodes you expect this cluster to have, taking
  # into consideration future growth. This number must be at least equal to the
  # number of nodes currently in your cluster. This number is used for port
  # reservation in the sidecar. This is the value that will be used for the
  # `max_nodes` argument of the `properties` block in the repository resource
  # (see resource `mongodb_repo` below).
  mongodb_max_nodes = 5

  # See `mongodb_port_alloc_range_low` and `mongodb_port_alloc_range_high` in
  # the cyral_sidecar module configuration.
  mongodb_ports_low  = 27017
  mongodb_ports_high = local.mongodb_ports_low + local.mongodb_max_nodes

  # All ports that will be used by MongoDB. This list must contain at least
  # `local.mongodb_max_nodes` ports
  mongodb_ports = range(local.mongodb_ports_low, local.mongodb_ports_high)
}

resource "cyral_repository" "mongodb_repo" {
  name = "mongodb_repo"
  type = "mongodb"

  # Specify the address or hostname of the endpoint of one node in the MongoDB
  # replica set. Cyral will automatically/dynamically identify the remaining
  # nodes of the replication cluster.
  host = "mycluster-shard-00-01.example.mongodb.net"

  port = local.mongodb_ports_low
  properties {
    mongodb_replica_set {
      max_nodes = local.mongodb_max_nodes

      # Specify the replica set identifier, a string value that identifies the
      # MongoDB replica set cluster. To find your replica set ID, see our
      # article:
      #
      # * https://cyral.freshdesk.com/a/solutions/articles/44002241594
      replica_set_id = "my-replica-set-id"
    }
  }
}

resource "cyral_repository_conf_auth" "mongodb_repo_auth_config" {
  repository_id     = cyral_repository.mongodb_repo.id
  identity_provider = module.cyral_idp_okta.integration_idp_okta_id
  # Repo TLS is required to allow the sidecar to communicate with MongoDB Atlas.
  repo_tls          = "enable"
}

resource "cyral_repository_binding" "mongodb_repo_binding" {
  repository_id                 = cyral_repository.mongodb_repo.id
  sidecar_id                    = cyral_sidecar.mongodb_sidecar.id
  listener_port                 = local.mongodb_ports_low
  sidecar_as_idp_access_gateway = true
}

resource "cyral_sidecar" "mongodb_sidecar" {
  name              = "MongoDBSidecar"
  deployment_method = "terraform"
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
  sidecar_id = cyral_sidecar.mongodb_sidecar.id
}

module "cyral_sidecar" {
  # Set the desired sidecar version. This information can be extracted from the
  # template downloaded from the UI.
  sidecar_version = "v2.32.0"

  source = "cyralinc/sidecar-aws/cyral"
  # Use the module version that is compatible with your sidecar. This
  # information can be extracted from the template downloaded from the UI.
  version = "2.8.1"

  sidecar_id = cyral_sidecar.mongodb_sidecar.id

  name_prefix = local.sidecar.sidecar_name_prefix

  control_plane = local.control_plane

  repositories_supported = ["mongodb"]

  # Specify all the ports that can be used in the sidecar. Below, we allocate
  # ports for MongoDB only. If you wish to bind this sidecar to other types of
  # repositories, make sure to allocate additional ports for them.
  sidecar_ports = local.mongodb_ports

  # Lower and upper limit values for the port allocation range reserved for
  # MongoDB. This range must correspond to the range of ports declared in
  # sidecar_ports that will be used for MongoDB. If you assign to sidecar_ports
  # the consecutive ports 27017, 27018 and 27019 for MongoDB utilization, it
  # means that the corresponding mongodb_port_alloc_range_low is 27017 and
  # mongodb_port_alloc_range_high is 27019. If you want to use a range of 10
  # ports for MongoDB, then you need to add all consecutive ports to
  # sidecar_ports (ex: 27017, 27018, 27019, 27020, 27021, 27022, 27023, 27024,
  # 27025, 27026) and define mongodb_port_alloc_range_low = 27017 and
  # mongodb_port_alloc_range_high = 27026.
  mongodb_port_alloc_range_low  = local.mongodb_ports_low
  mongodb_port_alloc_range_high = local.mongodb_ports_high

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
  secrets_location = "/cyral/sidecars/${cyral_sidecar.mongodb_sidecar.id}/secrets"

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
