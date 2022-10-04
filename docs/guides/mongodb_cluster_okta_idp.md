---
page_title: "Setup SSO access to MongoDB cluster using Okta IdP"
---

In this guide we deploy a sidecar, a repository associated with a MongoDB
replica set, and an Okta integration with the Cyral control plane. This setup
enables you to allow your users to access the database using SSO authentication
with Okta.

The guide is self-contained, so there are no prerequisites, except that you must
have the right credentials for Cyral, Okta and AWS providers. In each step
below, simply copy the code and paste locally, adjusting the argument values to
your needs. In some cases, we suggest the names of the files, but these names
don't need to be followed strictly.

## Configure required providers

Set required provider versions:

```terraform
terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = ">= 3.0.0"
    }
    okta = {
      source = "okta/okta"
    }
  }
}
```

Configure the providers:

```terraform
locals {
  # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
  control_plane_host = ""
  # Set the control plane API port
  control_plane_port = 443
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
```

## Configure Sidecar and MongoDB repository

Put the following Terraform configuration in `sidecar.tf`. Follow the comments
and replace argument values according to your needs. The `template` mentioned in
the comments is the sidecar Terraform deployment template for AWS you may download from
the Cyral control plane or see [in the public repository](https://github.com/cyralinc/terraform-cyral-sidecar-aws).

```terraform
locals {
  sidecar = {
    # If you would like to use other log integration, download a new
    # template from the UI and copy the log integration configuration
    # or follow this module documentation.
    log_integration = "cloudwatch"
    # Set the desired EC2 instance type for the auto scaling group.
    instance_type = "t3.medium"
    # Set to true if you want a sidecar deployed with an
    # internet-facing load balancer (requires a public subnet).
    public_sidecar = true

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
    # registry.  These parameters can be found on the sidecar
    # Terraform template downloaded from the UI.
    container_registry = {
      name         = "" # see container_registry in the downloaded template
      username     = "" # see container_registry_username in the downloaded template
      registry_key = "" # see container_registry_key in the downloaded template
    }
  }

  # Specify the maximum number of nodes you expect this cluster to
  # have, taking into consideration future growth. This number must be
  # at least equal to the number of nodes currently in your
  # cluster. This number is used for port reservation in the
  # sidecar. This is the value that will be used for the `max_nodes`
  # argument of the `properties` block in the repository resource (see
  # resource `mongodb_repo` below).
  mongodb_max_nodes = 5

  # See `mongodb_port_alloc_range_low` and
  # `mongodb_port_alloc_range_high` in the cyral_sidecar module
  # configuration.
  mongodb_ports_low  = 27017
  mongodb_ports_high = local.mongodb_ports_low + local.mongodb_max_nodes

  # All ports that will be used by MongoDB. This range must span at
  # least the `local.mongodb_max_nodes` number of ports. Note that the
  # port number you pass as the second argument to this function is
  # not included in the range. For example, to set port 27021 as your
  # uppermost port number, the second argument must be 27022.
  mongodb_ports = range(local.mongodb_ports_low, local.mongodb_ports_high)
}

resource "cyral_repository" "mongodb_repo" {
  name = "mongodb_repo"
  type = "mongodb"

  # Specify the address or hostname of the endpoint of one node in the
  # MongoDB replica set. Cyral will automatically/dynamically identify
  # the remaining nodes of the replication cluster.
  host = ""

  port = local.mongodb_ports_low
  properties {
    mongodb_replica_set {
      max_nodes = local.mongodb_max_nodes

      # Specify the replica set identifier, a string value that
      # identifies the MongoDB replica set cluster. To find your
      # replica set ID, see our article:
      #
      # * https://cyral.freshdesk.com/a/solutions/articles/44002241594
      replica_set_id = ""
    }
  }
}

resource "cyral_repository_conf_auth" "mongodb_repo_auth_config" {
  repository_id     = cyral_repository.mongodb_repo.id
  identity_provider = module.cyral_idp_okta.integration_idp_okta_id
  # Repo TLS is required to allow the sidecar to communicate with
  # MongoDB Atlas.
  repo_tls = "enable"
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
  # Set the desired sidecar version. This information can be extracted
  # from the template downloaded from the UI.
  sidecar_version = "v3.0.0"

  source = "cyralinc/sidecar-ec2/aws"
  # Use the module version that is compatible with your sidecar. This
  # information can be extracted from the template downloaded from the
  # UI.
  version = "~> 3.0.0"

  sidecar_id = cyral_sidecar.mongodb_sidecar.id

  control_plane = local.control_plane_host

  repositories_supported = ["mongodb"]

  # Specify all the ports that can be used in the sidecar. Below, we
  # allocate ports for MongoDB only. If you wish to bind this sidecar
  # to other types of repositories, make sure to allocate additional
  # ports for them.
  sidecar_ports = local.mongodb_ports

  # Lower and upper limit values for the port allocation range
  # reserved for MongoDB. This range must correspond to the range of
  # ports declared in sidecar_ports that will be used for MongoDB. If
  # you assign to sidecar_ports the consecutive ports 27017, 27018 and
  # 27019 for MongoDB utilization, it means that the corresponding
  # mongodb_port_alloc_range_low is 27017 and
  # mongodb_port_alloc_range_high is 27019. If you want to use a range
  # of 10 ports for MongoDB, then you need to add all consecutive
  # ports to sidecar_ports (ex: 27017, 27018, 27019, 27020, 27021,
  # 27022, 27023, 27024, 27025, 27026) and define
  # mongodb_port_alloc_range_low = 27017 and
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
```

## Configure a user account with the database credentials

Put the following in `user_account.tf`.

```terraform
locals {
  database_credentials = {
    # Native database credentials.
    username = ""
    password = ""
  }
}

resource "aws_secretsmanager_secret" "mongodb_creds" {
  # The sidecar deployed using our AWS sidecar module has access to
  # all secrets with the prefix '/cyral/' in the region it is
  # deployed.
  name = join("", [
    "/cyral/dbsecrets/",
    cyral_repository.mongodb_repo.id
  ])
}

resource "aws_secretsmanager_secret_version" "mongodb_creds_version" {
  secret_id     = aws_secretsmanager_secret.mongodb_creds.id
  secret_string = jsonencode(local.database_credentials)
}

resource "cyral_repository_user_account" "mongodb_user_account" {
  repository_id = cyral_repository.mongodb_repo.id
  # Set the name of the user account. This can be chosen freely.
  name = ""
  # Set the name of the target MongoDB database.
  auth_database_name = ""
  auth_scheme {
    aws_secrets_manager {
      secret_arn = aws_secretsmanager_secret.mongodb_creds.arn
    }
  }
}
```

## Configure Okta IdP

Finally, configure the Okta integration with the Cyral control plane. Put the
code in the file `integration.tf`.

```terraform
locals {
  # Set the name to be displayed for this integration in Okta's UI.
  okta_app_name = "Cyral"
  # Set the name to be displayed for this integration in Cyral's UI.
  okta_integration_name = "my-okta-integration"
}

module "cyral_idp_okta" {
  source  = "cyralinc/idp/okta"
  version = ">= 3.0.2"

  tenant = "default"

  control_plane = "${local.control_plane_host}:${local.control_plane_port}"

  okta_app_name        = local.okta_app_name
  idp_integration_name = local.okta_integration_name
}

resource "cyral_repository_access_rules" "access_rules" {
  repository_id   = cyral_repository.mongodb_repo.id
  user_account_id = cyral_repository_user_account.mongodb_user_account.user_account_id
  rule {
    identity {
      type = "username"
      name = ""
    }
  }
}
```

## Testing

To learn how to access a repository through the sidecar, see [Connect to a
repository](https://cyral.com/docs/connect/repo-connect/#connect-to-a-data-repository-with-sso-credentials).

## Next steps

In this guide, we configured a _user_ identity from Okta. You may also choose to
use group identities. For more information on Okta SSO integration, visit [SSO with
Okta](https://cyral.com/docs/sso/okta/sso) or our
[Terraform IdP integration module for Okta](https://registry.terraform.io/modules/cyralinc/idp-okta/cyral/latest).
